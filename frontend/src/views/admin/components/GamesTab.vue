<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Plus, Search, Check, Copy, Edit3, Trash2, Gamepad2, Star, Sparkles, RefreshCw, StopCircle, AlertTriangle, Lock, Unlock } from 'lucide-vue-next'
import { useAdminStore } from '../../../stores/admin.store'
import { useGameStore } from '../../../stores/game.store'
import Button from '../../../components/ui/Button.vue'
import Input from '../../../components/ui/Input.vue'
import Modal from '../../../components/ui/Modal.vue'
import type { Game } from '../../../types'

const emit = defineEmits<{
  (e: 'add'): void
  (e: 'edit', game: Game): void
}>()

const adminStore = useAdminStore()
const gameStore = useGameStore()

const gameSearch = ref('')
const copiedId = ref<string | null>(null)
const rescrapingSingleId = ref<string | null>(null)

// Confirmation Modal States
const showRescrapeConfirmModal = ref(false)
const showCancelRescrapeModal = ref(false)
const gamePendingDelete = ref<Game | null>(null)

onMounted(async () => {
  await adminStore.fetchRescrapeStatus()
})

const filteredGames = computed(() => {
  if (!gameSearch.value.trim()) return gameStore.games
  const q = gameSearch.value.toLowerCase()
  return gameStore.games.filter(g => 
    g.title.toLowerCase().includes(q) ||
    (g.genres || []).some(genre => genre.toLowerCase().includes(q)) ||
    g.drive_file_id.toLowerCase().includes(q)
  )
})

const copyDriveId = async (id: string) => {
  await navigator.clipboard.writeText(id)
  copiedId.value = id
  setTimeout(() => {
    copiedId.value = null
  }, 2000)
}

const openDeleteModal = (game: Game) => {
  gamePendingDelete.value = game
}

const confirmDelete = async () => {
  if (gamePendingDelete.value) {
    const game = gamePendingDelete.value
    gamePendingDelete.value = null
    await adminStore.deleteGame(game.id, game.title)
  }
}

const openRescrapeModal = () => {
  if (gameStore.games.length === 0 || adminStore.isRescraping) return
  showRescrapeConfirmModal.value = true
}

const confirmRescrape = async () => {
  showRescrapeConfirmModal.value = false
  await adminStore.triggerRescrapeAll()
}

const openCancelModal = () => {
  showCancelRescrapeModal.value = true
}

const confirmCancel = async () => {
  showCancelRescrapeModal.value = false
  await adminStore.cancelRescrape()
}

const handleToggleLock = async (game: Game) => {
  await adminStore.toggleGameLock(game)
}

const handleRescrapeSingle = async (game: Game) => {
  if (rescrapingSingleId.value) return
  rescrapingSingleId.value = game.id
  try {
    await adminStore.rescrapeSingleGame(game.id)
  } finally {
    rescrapingSingleId.value = null
  }
}
</script>

<template>
  <div class="space-y-4 animate-fadeIn">
    
    <!-- Header row with Filter, Re-scrape All & Add Game buttons -->
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 pb-3 border-b border-white/5">
      <div>
        <h1 class="text-xl sm:text-2xl font-black text-white tracking-tight flex items-center gap-2.5">
          <span>Games Manager</span>
          <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-white/5 border border-white/10 text-slate-400">
            {{ gameStore.games.length }}
          </span>
        </h1>
        <p class="text-xs text-slate-400 mt-1 font-normal">Add, edit metadata, and manage vault game packages</p>
      </div>

      <div class="flex items-center gap-2.5 w-full sm:w-auto flex-wrap sm:flex-nowrap">
        <!-- Search filter -->
        <div class="relative flex-1 sm:w-60">
          <Input
            v-model="gameSearch"
            placeholder="Filter by title, genre, ID..."
            size="sm"
            clearable
          >
            <template #prefix>
              <Search class="w-3.5 h-3.5 text-slate-400" />
            </template>
          </Input>
        </div>

        <!-- RE-SCRAPE ALL BUTTON (Background Worker) -->
        <Button
          variant="secondary"
          size="md"
          :disabled="adminStore.isRescraping || gameStore.games.length === 0"
          @click="openRescrapeModal"
          class="border-indigo-500/30 text-indigo-300 hover:bg-indigo-500/10 font-semibold cursor-pointer shrink-0"
          title="Re-scrape all games metadata and Steam 600x900 capsule covers"
        >
          <template #prefix>
            <RefreshCw v-if="adminStore.isRescraping" class="w-4 h-4 animate-spin text-indigo-400" />
            <Sparkles v-else class="w-4 h-4 text-indigo-400" />
          </template>
          <span>{{ adminStore.isRescraping ? 'Scraping...' : 'Re-scrape All' }}</span>
        </Button>

        <!-- ADD GAME BUTTON -->
        <Button
          variant="primary"
          size="md"
          @click="emit('add')"
          class="shadow-lg shadow-indigo-600/30 font-semibold shrink-0 cursor-pointer"
        >
          <template #prefix>
            <Plus class="w-4 h-4" />
          </template>
          <span>Add Game</span>
        </Button>
      </div>
    </div>

    <!-- Live Background Re-scrape Progress Banner -->
    <div 
      v-if="adminStore.rescrapeStatus?.status === 'running'"
      class="p-4 rounded-2xl bg-indigo-950/40 border border-indigo-500/30 backdrop-blur-md shadow-lg shadow-indigo-950/50 flex flex-col gap-2.5 animate-fadeIn"
    >
      <div class="flex items-center justify-between gap-3">
        <div class="flex items-center gap-2.5 min-w-0">
          <RefreshCw class="w-4 h-4 text-indigo-400 animate-spin shrink-0" />
          <div class="min-w-0">
            <span class="text-xs font-bold text-white flex items-center gap-2">
              <span>Re-scraping Steam & Metadata in background</span>
              <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-indigo-500/20 text-indigo-300 border border-indigo-500/30">
                {{ adminStore.rescrapeStatus.processed }} / {{ adminStore.rescrapeStatus.total }}
                ({{ Math.round((adminStore.rescrapeStatus.processed / (adminStore.rescrapeStatus.total || 1)) * 100) }}%)
              </span>
            </span>
            <p v-if="adminStore.rescrapeStatus.current_game" class="text-[11px] text-slate-300 truncate mt-0.5">
              Current: <span class="font-semibold text-indigo-200">"{{ adminStore.rescrapeStatus.current_game }}"</span>
            </p>
          </div>
        </div>

        <button
          @click="openCancelModal"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-rose-500/10 hover:bg-rose-500/20 border border-rose-500/30 text-rose-300 text-xs font-bold transition-colors cursor-pointer shrink-0"
        >
          <StopCircle class="w-3.5 h-3.5" />
          <span>Cancel</span>
        </button>
      </div>

      <!-- Animated Progress Bar -->
      <div class="w-full h-1.5 rounded-full bg-white/10 overflow-hidden">
        <div 
          class="h-full bg-linear-to-r from-indigo-500 to-sky-400 transition-all duration-300 rounded-full"
          :style="{ width: `${Math.round((adminStore.rescrapeStatus.processed / (adminStore.rescrapeStatus.total || 1)) * 100)}%` }"
        ></div>
      </div>
    </div>

    <!-- Games Table Container (Dark Cinema Obsidian Style) -->
    <div class="rounded-2xl bg-white/2 border border-white/5 overflow-hidden backdrop-blur-xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead class="bg-white/3 border-b border-white/5 text-slate-400 font-semibold uppercase tracking-wider text-[10px]">
            <tr>
              <th class="py-3.5 px-4">Game</th>
              <th class="py-3.5 px-4">Year</th>
              <th class="py-3.5 px-4">Rating</th>
              <th class="py-3.5 px-4">Genres</th>
              <th class="py-3.5 px-4">Drive ID</th>
              <th class="py-3.5 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5 text-slate-300">
            <tr v-if="filteredGames.length === 0">
              <td colspan="6" class="py-12 text-center text-slate-500">
                No games found matching filter criteria.
              </td>
            </tr>
            <tr 
              v-for="game in filteredGames" 
              :key="game.id"
              class="hover:bg-white/3 transition-colors"
            >
              <!-- Title & Poster -->
              <td class="py-3 px-4">
                <div class="flex items-center gap-3">
                  <img
                    v-if="game.cover_url"
                    :src="game.cover_url"
                    :alt="game.title"
                    class="w-10 h-14 rounded-lg object-cover border border-white/10 shrink-0"
                  />
                  <div v-else class="w-10 h-14 rounded-lg bg-white/5 flex items-center justify-center shrink-0">
                    <Gamepad2 class="w-4 h-4 text-slate-600" />
                  </div>
                  <div class="min-w-0">
                    <div class="flex items-center gap-2">
                      <p class="font-bold text-white text-xs truncate max-w-xs">{{ game.title }}</p>
                      <span 
                        v-if="game.is_locked" 
                        class="px-1.5 py-0.5 rounded text-[9px] font-bold bg-amber-500/20 text-amber-300 border border-amber-500/30 shrink-0 flex items-center gap-1"
                        title="Untracked: Skipped in Re-scrape All"
                      >
                        <Lock class="w-2.5 h-2.5" />
                        <span>Locked</span>
                      </span>
                    </div>
                    <p class="text-[10px] text-slate-400 truncate max-w-xs mt-0.5">{{ game.description || 'No description' }}</p>
                  </div>
                </div>
              </td>

              <!-- Year -->
              <td class="py-3 px-4 font-semibold text-slate-300">
                {{ game.release_year || 'N/A' }}
              </td>

              <!-- Rating -->
              <td class="py-3 px-4">
                <span v-if="game.rating > 0" class="flex items-center gap-1 font-bold text-amber-300">
                  <Star class="w-3 h-3 fill-amber-400 text-amber-400" />
                  {{ game.rating.toFixed(1) }}
                </span>
                <span v-else class="text-slate-500">-</span>
              </td>

              <!-- Genres -->
              <td class="py-3 px-4">
                <div class="flex flex-wrap gap-1 max-w-xs">
                  <span
                    v-for="g in (game.genres || []).slice(0, 3)"
                    :key="g"
                    class="px-2 py-0.5 rounded bg-white/5 border border-white/5 text-[10px] text-slate-300 font-medium"
                  >
                    {{ g }}
                  </span>
                </div>
              </td>

              <!-- Drive ID with Copy Button -->
              <td class="py-3 px-4 font-mono text-[11px]">
                <button
                  @click="copyDriveId(game.drive_file_id)"
                  class="flex items-center gap-1.5 px-2 py-1 rounded-md bg-black/40 hover:bg-black/60 border border-white/10 text-indigo-300 hover:text-white transition-colors cursor-pointer max-w-44"
                  title="Click to copy ID"
                >
                  <span class="truncate">{{ game.drive_file_id }}</span>
                  <Check v-if="copiedId === game.drive_file_id" class="w-3 h-3 text-emerald-400 shrink-0" />
                  <Copy v-else class="w-3 h-3 text-slate-500 shrink-0" />
                </button>
              </td>

              <!-- Actions (Compact Icon Buttons) -->
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <!-- Single Game Re-scrape Button -->
                  <button
                    @click="handleRescrapeSingle(game)"
                    :disabled="rescrapingSingleId === game.id"
                    title="Re-scrape Steam & RAWG metadata for this game"
                    class="p-2 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-indigo-300 border border-white/5 transition-all cursor-pointer disabled:opacity-50"
                  >
                    <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin text-indigo-400': rescrapingSingleId === game.id }" />
                  </button>

                  <!-- Lock / Unlock Toggle Button -->
                  <button
                    @click="handleToggleLock(game)"
                    :title="game.is_locked ? 'Locked (skipped in Re-scrape All). Click to unlock.' : 'Tracked. Click to lock (skip in Re-scrape All).'"
                    class="p-2 rounded-lg border transition-all cursor-pointer"
                    :class="game.is_locked 
                      ? 'bg-amber-500/10 text-amber-400 border-amber-500/30 hover:bg-amber-500/20' 
                      : 'bg-white/5 text-slate-400 hover:text-slate-200 border-white/5 hover:bg-white/10'"
                  >
                    <Lock v-if="game.is_locked" class="w-3.5 h-3.5" />
                    <Unlock v-else class="w-3.5 h-3.5" />
                  </button>

                  <!-- Edit Game -->
                  <button
                    @click="emit('edit', game)"
                    title="Edit Game"
                    class="p-2 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-white border border-white/5 transition-all cursor-pointer"
                  >
                    <Edit3 class="w-3.5 h-3.5 text-indigo-400" />
                  </button>

                  <!-- Delete Game -->
                  <button
                    @click="openDeleteModal(game)"
                    title="Delete Game"
                    class="p-2 rounded-lg bg-white/5 hover:bg-rose-500/20 text-slate-400 hover:text-rose-400 border border-white/5 transition-all cursor-pointer"
                  >
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 1. RE-SCRAPE CONFIRMATION MODAL -->
    <Modal
      v-if="showRescrapeConfirmModal"
      v-model="showRescrapeConfirmModal"
      title="Re-scrape Vault Metadata"
      maxWidth="md"
    >
      <div class="space-y-4 text-xs">
        <div class="p-3.5 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 flex items-start gap-3">
          <Sparkles class="w-5 h-5 text-indigo-400 shrink-0 mt-0.5" />
          <div class="space-y-1">
            <p class="font-bold text-white text-sm">Steam & RAWG Hybrid Scraper</p>
            <p class="text-slate-300 leading-relaxed">
              This operation will query RAWG for accurate ratings, descriptions, and release years, and enrich each game with official <span class="font-semibold text-sky-300">Steam 600x900 vertical capsule art</span> for all <span class="font-bold text-white">{{ gameStore.games.length }} games</span>.
            </p>
          </div>
        </div>

        <p class="text-slate-400 leading-relaxed">
          The task runs asynchronously in the background with automatic pacing. Any games marked as <span class="font-semibold text-amber-300">Locked</span> will be skipped. Your Google Drive files, download links, and IDs remain completely untouched.
        </p>
      </div>

      <template #footer>
        <div class="flex items-center justify-end gap-3">
          <Button
            variant="secondary"
            size="md"
            @click="showRescrapeConfirmModal = false"
            class="cursor-pointer"
          >
            Cancel
          </Button>

          <Button
            variant="primary"
            size="md"
            @click="confirmRescrape"
            class="shadow-lg shadow-indigo-600/30 cursor-pointer"
          >
            <template #prefix>
              <Sparkles class="w-4 h-4" />
            </template>
            <span>Start Re-scrape</span>
          </Button>
        </div>
      </template>
    </Modal>

    <!-- 2. DELETE GAME CONFIRMATION MODAL -->
    <Modal
      v-if="gamePendingDelete"
      :modelValue="!!gamePendingDelete"
      @update:modelValue="gamePendingDelete = null"
      title="Delete Game"
      maxWidth="sm"
    >
      <div class="space-y-3 text-xs">
        <div class="p-3.5 rounded-2xl bg-rose-500/10 border border-rose-500/20 flex items-start gap-3">
          <AlertTriangle class="w-5 h-5 text-rose-400 shrink-0 mt-0.5" />
          <div>
            <p class="font-bold text-white text-sm">Permanent Action</p>
            <p class="text-slate-300 mt-1 leading-relaxed">
              Are you sure you want to permanently remove <span class="font-bold text-white">"{{ gamePendingDelete?.title }}"</span> from the game vault?
            </p>
          </div>
        </div>
        <p class="text-slate-400">
          This record will be deleted from the database. Google Drive source files will not be affected.
        </p>
      </div>

      <template #footer>
        <div class="flex items-center justify-end gap-3">
          <Button
            variant="secondary"
            size="md"
            @click="gamePendingDelete = null"
            class="cursor-pointer"
          >
            Cancel
          </Button>

          <Button
            variant="danger"
            size="md"
            @click="confirmDelete"
            class="shadow-lg shadow-rose-600/30 cursor-pointer"
          >
            <template #prefix>
              <Trash2 class="w-4 h-4" />
            </template>
            <span>Delete Permanently</span>
          </Button>
        </div>
      </template>
    </Modal>

    <!-- 3. CANCEL RE-SCRAPE WORKER MODAL -->
    <Modal
      v-if="showCancelRescrapeModal"
      v-model="showCancelRescrapeModal"
      title="Stop Re-scrape Worker"
      maxWidth="sm"
    >
      <div class="space-y-3 text-xs">
        <div class="p-3.5 rounded-2xl bg-amber-500/10 border border-amber-500/20 flex items-start gap-3">
          <StopCircle class="w-5 h-5 text-amber-400 shrink-0 mt-0.5" />
          <div>
            <p class="font-bold text-white text-sm">Abort Active Task?</p>
            <p class="text-slate-300 mt-1 leading-relaxed">
              Do you want to stop the background re-scrape task now? Games updated up to this point will remain safely preserved.
            </p>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex items-center justify-end gap-3">
          <Button
            variant="secondary"
            size="md"
            @click="showCancelRescrapeModal = false"
            class="cursor-pointer"
          >
            Keep Running
          </Button>

          <Button
            variant="danger"
            size="md"
            @click="confirmCancel"
            class="shadow-lg shadow-rose-600/30 cursor-pointer"
          >
            <span>Stop Worker</span>
          </Button>
        </div>
      </template>
    </Modal>

  </div>
</template>
