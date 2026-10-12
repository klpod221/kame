<script setup lang="ts">
import { ref } from 'vue'
import { Sparkles, Save, Image as ImageIcon, Link as LinkIcon } from 'lucide-vue-next'
import type { GameCandidate } from '../../../types'
import { useAdminStore } from '../../../stores/admin.store'
import Modal from '../../../components/ui/Modal.vue'
import Input from '../../../components/ui/Input.vue'
import Textarea from '../../../components/ui/Textarea.vue'
import Button from '../../../components/ui/Button.vue'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created'): void
}>()

const adminStore = useAdminStore()

const driveInput = ref('')
const searchInput = ref('')
const isSearching = ref(false)
const isSaving = ref(false)
const searchError = ref<string | null>(null)
const candidates = ref<GameCandidate[]>([])

// Form state
const title = ref('')
const description = ref('')
const genresStr = ref('')
const coverUrl = ref('')
const backgroundUrl = ref('')
const releaseYear = ref<number | undefined>(undefined)
const rating = ref<number | undefined>(undefined)

const handleSearch = async () => {
  if (!searchInput.value.trim()) return

  isSearching.value = true
  searchError.value = null
  candidates.value = []

  try {
    const list = await adminStore.scrapeMetadata(searchInput.value)
    candidates.value = list
    if (list.length === 0) {
      searchError.value = 'No games found. You can fill the information manually.'
    } else {
      selectCandidate(list[0])
    }
  } catch (err: any) {
    searchError.value = err.message || 'Scraping failed'
  } finally {
    isSearching.value = false
  }
}

const selectCandidate = (item: GameCandidate) => {
  title.value = item.title
  description.value = item.description || ''
  genresStr.value = (item.genres || []).join(', ')
  coverUrl.value = item.cover_url || ''
  backgroundUrl.value = item.background_url || item.cover_url || ''
  releaseYear.value = item.release_year || undefined
  rating.value = item.rating ? Number((Math.round(item.rating * 10) / 10).toFixed(1)) : undefined
}

const extractDriveId = (input: string): string => {
  const trimmed = input.trim()
  const match1 = trimmed.match(/\/file\/d\/([a-zA-Z0-9_-]+)/)
  if (match1 && match1[1]) return match1[1]

  const match2 = trimmed.match(/[?&]id=([a-zA-Z0-9_-]+)/)
  if (match2 && match2[1]) return match2[1]

  return trimmed
}

const handleSubmit = async () => {
  const fileId = extractDriveId(driveInput.value)
  if (!fileId) {
    alert('Please enter a Google Drive File ID or link')
    return
  }

  if (!title.value.trim()) {
    alert('Game title is required')
    return
  }

  isSaving.value = true
  try {
    const genres = genresStr.value
      .split(',')
      .map(s => s.trim())
      .filter(Boolean)

    await adminStore.createGame({
      title: title.value.trim(),
      description: description.value.trim(),
      genres,
      cover_url: coverUrl.value.trim(),
      background_url: backgroundUrl.value.trim(),
      drive_file_id: fileId,
      release_year: releaseYear.value || 0,
      rating: rating.value || 0,
    })

    emit('created')
    emit('close')
  } catch (err) {
    console.error(err)
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <Modal title="Add Game to Vault" maxWidth="2xl" @close="emit('close')">
    <div class="space-y-6">
      
      <!-- Auto-Scrape Search Bar -->
      <div class="p-4 rounded-2xl bg-slate-950/70 border border-slate-800 space-y-3">
        <label class="text-xs font-bold text-slate-300 flex items-center gap-1.5">
          <Sparkles class="w-3.5 h-3.5 text-amber-400" />
          <span>Auto-Fetch Metadata (RAWG / IGDB)</span>
        </label>
        
        <div class="flex gap-2">
          <Input
            v-model="searchInput"
            placeholder="Search game title e.g. Cyberpunk 2077..."
            @keydown.enter.prevent="handleSearch"
          />
          <Button
            variant="primary"
            size="md"
            :loading="isSearching"
            @click="handleSearch"
          >
            <span>Search</span>
          </Button>
        </div>

        <p v-if="searchError" class="text-xs text-rose-400 font-medium">
          {{ searchError }}
        </p>

        <!-- Search Results Chips -->
        <div v-if="candidates.length > 0" class="flex flex-wrap gap-2 pt-1 max-h-36 overflow-y-auto pr-1">
          <button
            v-for="c in candidates"
            :key="c.source_id"
            type="button"
            @click="selectCandidate(c)"
            class="flex items-center gap-2 p-1.5 pr-3 rounded-xl border text-left text-xs transition-all cursor-pointer"
            :class="title === c.title 
              ? 'bg-indigo-600/30 border-indigo-500 text-white' 
              : 'bg-slate-900 border-slate-800 text-slate-400 hover:text-white'"
          >
            <img v-if="c.cover_url" :src="c.cover_url" class="w-7 h-9 rounded object-cover" />
            <div class="min-w-0">
              <p class="font-bold truncate max-w-40 text-xs">{{ c.title }}</p>
              <p class="text-[10px] text-slate-400">{{ c.release_year || 'N/A' }} • ★ {{ c.rating ? c.rating.toFixed(1) : '-' }}</p>
            </div>
          </button>
        </div>
      </div>

      <!-- Main Form Details -->
      <form @submit.prevent="handleSubmit" class="space-y-4">
        
        <!-- Drive Link or ID (Required) -->
        <div>
          <label class="text-xs font-bold text-slate-300 mb-1 flex items-center gap-1.5">
            <LinkIcon class="w-3.5 h-3.5 text-indigo-400" />
            <span>Google Drive File ID or Shareable Link *</span>
          </label>
          <Input
            v-model="driveInput"
            placeholder="10nIJS0O27dBTusVkz6sitk_wZklPp39Z or https://drive.google.com/file/d/..."
            required
          />
        </div>

        <!-- Title & Year & Rating Row -->
        <div class="grid grid-cols-1 sm:grid-cols-6 gap-3">
          <div class="sm:col-span-3">
            <label class="block text-xs font-bold text-slate-300 mb-1">Title *</label>
            <Input v-model="title" placeholder="e.g. Black Myth: Wukong" required />
          </div>

          <div class="sm:col-span-1">
            <label class="block text-xs font-bold text-slate-300 mb-1">Year</label>
            <Input v-model.number="releaseYear" type="number" placeholder="2024" />
          </div>

          <div class="sm:col-span-2">
            <label class="block text-xs font-bold text-slate-300 mb-1">Rating (0 - 5)</label>
            <Input v-model.number="rating" type="number" step="0.1" placeholder="4.5" />
          </div>
        </div>

        <!-- Genres -->
        <div>
          <label class="block text-xs font-bold text-slate-300 mb-1">Genres (comma separated)</label>
          <Input v-model="genresStr" placeholder="Action, RPG, Adventure" />
        </div>

        <!-- Description -->
        <div>
          <label class="block text-xs font-bold text-slate-300 mb-1">Overview Synopsis</label>
          <Textarea v-model="description" :rows="3" placeholder="Enter game story and overview synopsis..." />
        </div>

        <!-- Artwork URLs -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="text-xs font-bold text-slate-300 mb-1 flex items-center gap-1">
              <ImageIcon class="w-3 h-3 text-slate-400" />
              <span>Cover Poster URL</span>
            </label>
            <Input v-model="coverUrl" placeholder="https://..." />
          </div>

          <div>
            <label class="text-xs font-bold text-slate-300 mb-1 flex items-center gap-1">
              <ImageIcon class="w-3 h-3 text-slate-400" />
              <span>Backdrop Background URL</span>
            </label>
            <Input v-model="backgroundUrl" placeholder="https://..." />
          </div>
        </div>

        <!-- Modal Actions -->
        <div class="pt-4 flex items-center justify-end gap-3 border-t border-slate-800">
          <Button variant="ghost" size="md" @click="emit('close')">
            Cancel
          </Button>

          <Button
            type="submit"
            variant="primary"
            size="md"
            :loading="isSaving"
          >
            <Save class="w-4 h-4" />
            <span>Save to Vault</span>
          </Button>
        </div>

      </form>

    </div>
  </Modal>
</template>
