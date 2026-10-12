<script setup lang="ts">
import { ref, computed } from 'vue'
import { Plus, Search, Check, Copy, Edit3, Trash2, Gamepad2, Star } from 'lucide-vue-next'
import { useAdminStore } from '../../../stores/admin.store'
import { useGameStore } from '../../../stores/game.store'
import Button from '../../../components/ui/Button.vue'
import Input from '../../../components/ui/Input.vue'
import type { Game } from '../../../types'

const emit = defineEmits<{
  (e: 'add'): void
  (e: 'edit', game: Game): void
}>()

const adminStore = useAdminStore()
const gameStore = useGameStore()

const gameSearch = ref('')
const copiedId = ref<string | null>(null)

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

const handleDelete = async (game: Game) => {
  if (confirm(`Permanently remove "${game.title}" from vault?`)) {
    await adminStore.deleteGame(game.id, game.title)
  }
}
</script>

<template>
  <div class="space-y-4 animate-fadeIn">
    
    <!-- Header row with Filter & THE ONLY ADD GAME BUTTON -->
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

      <div class="flex items-center gap-3 w-full sm:w-auto">
        <!-- Search filter -->
        <div class="relative flex-1 sm:w-64">
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

        <!-- THE ONLY ADD GAME BUTTON -->
        <Button
          variant="primary"
          size="md"
          @click="emit('add')"
          class="shadow-lg shadow-indigo-600/30 font-semibold"
        >
          <template #prefix>
            <Plus class="w-4 h-4" />
          </template>
          <span>Add Game</span>
        </Button>
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
                    <p class="font-bold text-white text-xs truncate max-w-xs">{{ game.title }}</p>
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
                  <button
                    @click="emit('edit', game)"
                    title="Edit Game"
                    class="p-2 rounded-lg bg-white/5 hover:bg-white/10 text-slate-400 hover:text-white border border-white/5 transition-all cursor-pointer"
                  >
                    <Edit3 class="w-3.5 h-3.5 text-indigo-400" />
                  </button>

                  <button
                    @click="handleDelete(game)"
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

  </div>
</template>
