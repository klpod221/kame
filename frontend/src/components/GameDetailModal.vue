<script setup lang="ts">
import { ref } from 'vue'
import { X, Download, Star, Calendar, HardDrive, Trash2, Zap, AlertCircle } from 'lucide-vue-next'
import type { Game } from '../types'
import { useAuth } from '../composables/useAuth'
import { useGames } from '../composables/useGames'

const props = defineProps<{
  game: Game | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'deleted'): void
}>()

const { isAuthenticated, isAdmin, loginWithGoogle } = useAuth()
const { getDownloadUrl, deleteGame } = useGames()

const isDeleting = ref(false)
const deleteConfirm = ref(false)

const handleDownload = () => {
  if (!props.game) return

  if (!isAuthenticated.value) {
    loginWithGoogle()
    return
  }

  const url = getDownloadUrl(props.game.drive_file_id)
  
  // Trigger direct download via temporary anchor
  const a = document.createElement('a')
  a.href = url
  a.download = `${props.game.title}.zip`
  a.target = '_blank'
  a.rel = 'noopener noreferrer'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}

const handleDelete = async () => {
  if (!props.game) return
  if (!deleteConfirm.value) {
    deleteConfirm.value = true
    return
  }

  isDeleting.value = true
  try {
    await deleteGame(props.game.id)
    emit('deleted')
    emit('close')
  } catch (err) {
    alert('Failed to delete game')
  } finally {
    isDeleting.value = false
    deleteConfirm.value = false
  }
}
</script>

<template>
  <div 
    v-if="game" 
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md overflow-y-auto"
    @click.self="emit('close')"
  >
    <div class="relative w-full max-w-3xl rounded-3xl bg-slate-900 border border-slate-800 shadow-2xl overflow-hidden my-8">
      
      <!-- Hero Background Image -->
      <div class="relative h-64 md:h-80 w-full overflow-hidden bg-slate-950">
        <img
          v-if="game.background_url || game.cover_url"
          :src="game.background_url || game.cover_url"
          :alt="game.title"
          class="w-full h-full object-cover object-center filter brightness-75"
        />
        <div class="absolute inset-0 bg-linear-to-t from-slate-900 via-slate-900/60 to-transparent"></div>

        <!-- Close Button -->
        <button
          @click="emit('close')"
          class="absolute top-4 right-4 p-2.5 rounded-full bg-slate-950/60 hover:bg-slate-950 text-slate-300 hover:text-white backdrop-blur-md border border-slate-700/60 transition-colors cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Content Area -->
      <div class="px-6 md:px-8 pb-8 -mt-24 relative z-10">
        <div class="flex flex-col md:flex-row gap-6 items-start">
          
          <!-- Poster Mini -->
          <div class="w-36 md:w-44 shrink-0 aspect-3/4 rounded-2xl overflow-hidden shadow-2xl border-2 border-slate-700/60 bg-slate-950">
            <img
              v-if="game.cover_url"
              :src="game.cover_url"
              :alt="game.title"
              class="w-full h-full object-cover"
            />
          </div>

          <!-- Main Info Header -->
          <div class="flex-1 space-y-2.5 pt-2">
            <div class="flex flex-wrap items-center gap-2">
              <span 
                v-if="game.release_year" 
                class="flex items-center gap-1 text-xs font-medium px-2.5 py-1 rounded-lg bg-slate-800/80 text-slate-300 border border-slate-700/60"
              >
                <Calendar class="w-3.5 h-3.5 text-slate-400" />
                {{ game.release_year }}
              </span>

              <span 
                v-if="game.rating > 0" 
                class="flex items-center gap-1 text-xs font-semibold px-2.5 py-1 rounded-lg bg-amber-500/20 text-amber-300 border border-amber-500/30"
              >
                <Star class="w-3.5 h-3.5 fill-amber-400 text-amber-400" />
                {{ game.rating.toFixed(1) }} / 5.0
              </span>

              <span class="flex items-center gap-1 text-xs font-medium px-2.5 py-1 rounded-lg bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
                <HardDrive class="w-3.5 h-3.5" />
                Google Drive
              </span>
            </div>

            <h1 class="text-2xl md:text-3xl font-extrabold text-white tracking-tight">
              {{ game.title }}
            </h1>

            <!-- Genres -->
            <div class="flex flex-wrap gap-1.5 pt-1">
              <span
                v-for="genre in (game.genres || [])"
                :key="genre"
                class="text-xs font-medium px-2.5 py-1 rounded-lg bg-slate-800 text-slate-300"
              >
                {{ genre }}
              </span>
            </div>
          </div>
        </div>

        <!-- Description -->
        <div class="mt-6 pt-6 border-t border-slate-800">
          <h2 class="text-sm font-semibold uppercase tracking-wider text-slate-400 mb-2">About Game</h2>
          <div class="text-sm text-slate-300 leading-relaxed max-h-48 overflow-y-auto pr-2 space-y-2 whitespace-pre-line">
            {{ game.description || 'No description provided for this game.' }}
          </div>
        </div>

        <!-- Action Download Footer -->
        <div class="mt-8 pt-6 border-t border-slate-800 flex flex-col sm:flex-row items-center justify-between gap-4">
          
          <!-- Notice for non-authenticated user -->
          <div v-if="!isAuthenticated" class="flex items-center gap-2 text-xs text-amber-400">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>Sign in with Google to download file with full CDN speed.</span>
          </div>
          <div v-else class="text-xs text-slate-400 flex items-center gap-1.5">
            <Zap class="w-4 h-4 text-emerald-400" />
            <span>High-Speed Direct Stream (Google Drive API v3)</span>
          </div>

          <!-- Buttons -->
          <div class="flex items-center gap-3 w-full sm:w-auto">
            <!-- Delete Button (Admin only) -->
            <button
              v-if="isAdmin"
              @click="handleDelete"
              :disabled="isDeleting"
              class="flex items-center justify-center gap-2 px-4 py-3 rounded-xl border border-rose-500/30 text-rose-400 hover:bg-rose-500/10 transition-colors text-sm font-medium cursor-pointer"
            >
              <Trash2 class="w-4 h-4" />
              <span>{{ deleteConfirm ? 'Confirm Delete?' : 'Delete' }}</span>
            </button>

            <!-- Download Button -->
            <button
              @click="handleDownload"
              class="flex-1 sm:flex-initial flex items-center justify-center gap-2.5 px-6 py-3 rounded-xl bg-linear-to-r from-indigo-500 to-purple-600 hover:from-indigo-400 hover:to-purple-500 text-white font-semibold shadow-xl shadow-indigo-500/25 transition-all transform hover:scale-[1.02] cursor-pointer"
            >
              <Download class="w-5 h-5" />
              <span>{{ isAuthenticated ? 'Download Game' : 'Sign in to Download' }}</span>
            </button>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>
