<script setup lang="ts">
import { ref } from 'vue'
import { X, Sparkles, Save, Loader2, Image as ImageIcon, Link as LinkIcon } from 'lucide-vue-next'
import type { GameCandidate } from '../types'
import { useGames } from '../composables/useGames'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'created'): void
}>()

const { scrapeMetadata, createGame } = useGames()

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
    const list = await scrapeMetadata(searchInput.value)
    candidates.value = list
    if (list.length === 0) {
      searchError.value = 'No games found. You can fill the information manually.'
    } else {
      // Auto-select top 1 candidate
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
  rating.value = item.rating || undefined
}

const handleSubmit = async () => {
  if (!title.value.trim()) {
    alert('Title is required')
    return
  }
  if (!driveInput.value.trim()) {
    alert('Google Drive File ID or link is required')
    return
  }

  isSaving.value = true
  try {
    const genres = genresStr.value
      .split(',')
      .map(g => g.trim())
      .filter(g => g.length > 0)

    await createGame({
      title: title.value.trim(),
      description: description.value.trim(),
      genres,
      cover_url: coverUrl.value.trim(),
      background_url: backgroundUrl.value.trim(),
      drive_file_id: driveInput.value.trim(),
      release_year: releaseYear.value || 0,
      rating: rating.value || 0,
    })

    emit('created')
    emit('close')
  } catch (err: any) {
    alert(err.message || 'Failed to save game')
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div 
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md overflow-y-auto"
    @click.self="emit('close')"
  >
    <div class="relative w-full max-w-2xl rounded-3xl bg-slate-900 border border-slate-800 shadow-2xl p-6 md:p-8 my-8">
      
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-4 border-b border-slate-800">
        <div>
          <h2 class="text-xl font-bold text-white tracking-tight">Add Game to Catalog</h2>
          <p class="text-xs text-slate-400 mt-0.5">Scrape metadata automatically or enter manually.</p>
        </div>
        <button
          @click="emit('close')"
          class="p-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-white transition-colors cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Scraper Section -->
      <div class="mt-6 p-4 rounded-2xl bg-slate-950/60 border border-slate-800/80 space-y-3">
        <label class="block text-xs font-semibold text-slate-300">
          Search Game Metadata (RAWG & IGDB)
        </label>
        <div class="flex gap-2">
          <input
            v-model="searchInput"
            @keyup.enter="handleSearch"
            type="text"
            placeholder="e.g. Elden Ring, Baldur's Gate 3, or raw filename"
            class="flex-1 px-3.5 py-2 text-sm rounded-xl bg-slate-900 border border-slate-700/80 text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500"
          />
          <button
            @click="handleSearch"
            :disabled="isSearching || !searchInput.trim()"
            class="flex items-center gap-2 px-4 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white text-sm font-medium transition-colors cursor-pointer disabled:opacity-50"
          >
            <Loader2 v-if="isSearching" class="w-4 h-4 animate-spin" />
            <Sparkles v-else class="w-4 h-4" />
            <span>{{ isSearching ? 'Fetching...' : 'Scrape' }}</span>
          </button>
        </div>

        <!-- Candidate Thumbnails -->
        <div v-if="candidates.length > 0" class="pt-2">
          <p class="text-[11px] font-semibold text-slate-400 mb-2">Select matching result:</p>
          <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
            <div
              v-for="item in candidates"
              :key="item.source_id"
              @click="selectCandidate(item)"
              class="group relative flex items-center gap-2.5 p-2 rounded-xl bg-slate-900 border border-slate-800 hover:border-indigo-500 cursor-pointer transition-all"
              :class="{ 'border-indigo-500 bg-indigo-500/10': title === item.title }"
            >
              <img
                v-if="item.cover_url"
                :src="item.cover_url"
                class="w-9 h-12 rounded object-cover shrink-0"
              />
              <div class="overflow-hidden">
                <div class="text-xs font-medium text-slate-200 truncate group-hover:text-indigo-400">
                  {{ item.title }}
                </div>
                <div class="text-[10px] text-slate-500 capitalize">
                  {{ item.source }} • {{ item.release_year || 'N/A' }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <p v-if="searchError" class="text-xs text-amber-400 mt-1">
          {{ searchError }}
        </p>
      </div>

      <!-- Main Fields -->
      <form @submit.prevent="handleSubmit" class="mt-6 space-y-4">
        <!-- Google Drive File ID -->
        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1 flex items-center gap-1.5">
            <LinkIcon class="w-3.5 h-3.5 text-indigo-400" />
            Google Drive File ID or Shareable Link <span class="text-rose-400">*</span>
          </label>
          <input
            v-model="driveInput"
            type="text"
            required
            placeholder="1abcXYZ... or https://drive.google.com/file/d/1abcXYZ/view?usp=sharing"
            class="w-full px-3.5 py-2 text-sm rounded-xl bg-slate-950 border border-slate-800 text-white placeholder-slate-600 focus:outline-none focus:border-indigo-500 font-mono text-xs"
          />
        </div>

        <!-- Title -->
        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1">
            Game Title <span class="text-rose-400">*</span>
          </label>
          <input
            v-model="title"
            type="text"
            required
            placeholder="Official Game Title"
            class="w-full px-3.5 py-2 text-sm rounded-xl bg-slate-950 border border-slate-800 text-white placeholder-slate-600 focus:outline-none focus:border-indigo-500"
          />
        </div>

        <!-- Release Year & Rating -->
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Release Year</label>
            <input
              v-model.number="releaseYear"
              type="number"
              placeholder="e.g. 2024"
              class="w-full px-3.5 py-2 text-sm rounded-xl bg-slate-950 border border-slate-800 text-white placeholder-slate-600 focus:outline-none focus:border-indigo-500"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Rating (0 - 5.0)</label>
            <input
              v-model.number="rating"
              type="number"
              step="0.1"
              min="0"
              max="5"
              placeholder="e.g. 4.8"
              class="w-full px-3.5 py-2 text-sm rounded-xl bg-slate-950 border border-slate-800 text-white placeholder-slate-600 focus:outline-none focus:border-indigo-500"
            />
          </div>
        </div>

        <!-- Genres -->
        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1">Genres (comma separated)</label>
          <input
            v-model="genresStr"
            type="text"
            placeholder="Action, RPG, Open World"
            class="w-full px-3.5 py-2 text-sm rounded-xl bg-slate-950 border border-slate-800 text-white placeholder-slate-600 focus:outline-none focus:border-indigo-500"
          />
        </div>

        <!-- Cover Art URL -->
        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1 flex items-center gap-1.5">
            <ImageIcon class="w-3.5 h-3.5 text-indigo-400" />
            Cover Image URL
          </label>
          <input
            v-model="coverUrl"
            type="url"
            placeholder="https://..."
            class="w-full px-3.5 py-2 text-sm rounded-xl bg-slate-950 border border-slate-800 text-white placeholder-slate-600 focus:outline-none focus:border-indigo-500 font-mono text-xs"
          />
        </div>

        <!-- Background Image URL -->
        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1">Background Banner URL</label>
          <input
            v-model="backgroundUrl"
            type="url"
            placeholder="https://..."
            class="w-full px-3.5 py-2 text-sm rounded-xl bg-slate-950 border border-slate-800 text-white placeholder-slate-600 focus:outline-none focus:border-indigo-500 font-mono text-xs"
          />
        </div>

        <!-- Description -->
        <div>
          <label class="block text-xs font-semibold text-slate-300 mb-1">Description</label>
          <textarea
            v-model="description"
            rows="3"
            placeholder="Brief game synopsis..."
            class="w-full px-3.5 py-2 text-sm rounded-xl bg-slate-950 border border-slate-800 text-white placeholder-slate-600 focus:outline-none focus:border-indigo-500 leading-relaxed"
          ></textarea>
        </div>

        <!-- Submit Button -->
        <div class="pt-4 flex justify-end gap-3 border-t border-slate-800">
          <button
            type="button"
            @click="emit('close')"
            class="px-4 py-2 text-sm font-medium rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 transition-colors cursor-pointer"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="isSaving"
            class="flex items-center gap-2 px-5 py-2 text-sm font-semibold rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white shadow-lg shadow-indigo-600/25 transition-all cursor-pointer disabled:opacity-50"
          >
            <Loader2 v-if="isSaving" class="w-4 h-4 animate-spin" />
            <Save v-else class="w-4 h-4" />
            <span>{{ isSaving ? 'Saving...' : 'Add to Catalog' }}</span>
          </button>
        </div>
      </form>

    </div>
  </div>
</template>
