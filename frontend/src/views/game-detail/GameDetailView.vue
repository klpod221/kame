<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  ArrowLeft, 
  Download, 
  Star, 
  Calendar, 
  HardDrive, 
  Terminal, 
  ExternalLink, 
  Loader2, 
  Monitor
} from 'lucide-vue-next'
import type { Game } from '../../types'
import { useGameStore } from '../../stores/game.store'
import { formatBytes } from '../../utils/format'
import GameSpecsTable from './components/GameSpecsTable.vue'
import Button from '../../components/ui/Button.vue'
import Badge from '../../components/ui/Badge.vue'

const route = useRoute()
const router = useRouter()

const gameStore = useGameStore()

const game = ref<Game | null>(null)
const isLoadingGame = ref(false)

const loadGame = async () => {
  const id = route.params.id as string
  if (!id) return
  isLoadingGame.value = true
  try {
    game.value = await gameStore.fetchGameById(id)
  } finally {
    isLoadingGame.value = false
  }
}

onMounted(() => {
  loadGame()
})

watch(() => route.params.id, () => {
  loadGame()
})

const goBack = () => {
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-[#070a0f] text-slate-100 flex flex-col relative animate-fadeIn pb-24">
    
    <!-- 1. Hero Backdrop Banner (Cinematic Movie Style matching Jellyfin reference) -->
    <div class="relative w-full h-80 sm:h-96 md:h-115 lg:h-130 overflow-hidden select-none">
      <img
        v-if="game?.background_url || game?.cover_url"
        :src="game?.background_url || game?.cover_url"
        :alt="game?.title"
        class="w-full h-full object-cover object-center filter brightness-50 contrast-110"
      />
      <!-- Bottom Fade to Cinema Dark Background -->
      <div class="absolute inset-0 bg-linear-to-t from-[#070a0f] via-[#070a0f]/60 to-transparent"></div>
      <!-- Left Vignette Gradient -->
      <div class="absolute inset-0 bg-linear-to-r from-[#070a0f]/90 via-[#070a0f]/40 to-transparent"></div>
      <!-- Radial Ambient Highlight -->
      <div class="absolute top-0 right-1/4 w-125 h-75 bg-indigo-600/10 rounded-full blur-3xl pointer-events-none"></div>
    </div>

    <!-- 2. Top Floating Navigation Bar -->
    <header class="absolute top-0 left-0 right-0 z-20 px-4 sm:px-8 py-5 flex items-center justify-between">
      <button
        @click="goBack"
        class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-slate-900/80 hover:bg-slate-800 text-slate-300 hover:text-white backdrop-blur-md border border-white/10 text-xs sm:text-sm font-semibold transition-all shadow-xl active:scale-95 cursor-pointer"
      >
        <ArrowLeft class="w-4 h-4 text-indigo-400" />
        <span>Back to Library</span>
      </button>

      <div class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-slate-900/60 backdrop-blur-md border border-white/10 shadow-lg">
        <img src="/kame-48.png" alt="Kame" class="w-5 h-5 object-contain" />
        <span class="text-xs font-bold text-white tracking-wider">Kame</span>
      </div>
    </header>

    <!-- 3. Loading State -->
    <div v-if="isLoadingGame" class="flex-1 flex flex-col items-center justify-center py-24 relative z-10">
      <Loader2 class="w-10 h-10 text-indigo-400 animate-spin mb-3" />
      <p class="text-xs text-slate-400 font-medium">Loading game details...</p>
    </div>

    <!-- 4. Not Found State -->
    <div v-else-if="!game" class="flex-1 flex flex-col items-center justify-center py-24 relative z-10 space-y-4">
      <p class="text-slate-400 text-sm">Game not found or removed from vault.</p>
      <Button variant="primary" size="md" @click="goBack">
        Return to Library
      </Button>
    </div>

    <!-- 5. Main Content Container -->
    <div v-else class="max-w-7xl mx-auto w-full px-4 sm:px-8 lg:px-12 -mt-44 sm:-mt-56 md:-mt-64 relative z-10">
      
      <div class="flex flex-col md:flex-row gap-8 lg:gap-12 items-center md:items-start">
        
        <!-- Left: Floating Cinema Poster (Aspect 2:3 with deep shadow) -->
        <div class="w-60 sm:w-68 md:w-76 lg:w-80 shrink-0">
          <div class="aspect-2/3 rounded-2xl overflow-hidden shadow-2xl shadow-black/90 border border-white/15 bg-slate-900 relative group">
            <img
              v-if="game.cover_url"
              :src="game.cover_url"
              :alt="game.title"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-103"
            />
            <div v-else class="w-full h-full flex flex-col items-center justify-center bg-slate-950 text-slate-600 gap-2">
              <HardDrive class="w-12 h-12 text-slate-700" />
              <span class="text-xs font-bold text-slate-500">No Cover Art</span>
            </div>
          </div>
        </div>

        <!-- Right: Information & Actions Column -->
        <div class="flex-1 space-y-6 text-center md:text-left min-w-0">
          
          <!-- Title & Badges -->
          <div class="space-y-3">
            <h1 class="text-3xl sm:text-4xl lg:text-5xl font-black text-white tracking-tight leading-tight drop-shadow-md">
              {{ game.title }}
            </h1>

            <!-- Metadata Ribbon -->
            <div class="flex flex-wrap items-center justify-center md:justify-start gap-2.5 pt-1 text-xs font-semibold">
              <Badge v-if="game.release_year" variant="slate" class="bg-white/5 border-white/10 text-slate-300">
                <Calendar class="w-3.5 h-3.5 text-indigo-400" />
                <span>{{ game.release_year }}</span>
              </Badge>

              <Badge v-if="game.rating > 0" variant="amber" class="bg-amber-500/15 border-amber-500/30 text-amber-300 font-bold">
                <Star class="w-3.5 h-3.5 fill-amber-400 text-amber-400" />
                <span>{{ game.rating.toFixed(1) }}</span>
              </Badge>

              <Badge v-if="game.file_size && game.file_size > 0" variant="slate" class="bg-white/5 border-white/10 text-slate-300 font-mono">
                <HardDrive class="w-3.5 h-3.5 text-emerald-400" />
                <span>{{ formatBytes(game.file_size) }}</span>
              </Badge>

              <Badge variant="slate" class="bg-white/5 border-white/10 text-slate-300">
                <Monitor class="w-3.5 h-3.5 text-sky-400" />
                <span>PC (Windows)</span>
              </Badge>
            </div>
          </div>

          <!-- Action Buttons Bar (Sleek Compact Icon Buttons Style) -->
          <div class="flex items-center justify-center md:justify-start gap-2.5 pt-1">
            
            <!-- Download Icon Button -->
            <button
              @click="gameStore.downloadGame(game)"
              title="Download Game (Direct Stream)"
              class="flex items-center justify-center p-3 rounded-full bg-indigo-600 hover:bg-indigo-500 text-white shadow-xl shadow-indigo-600/30 hover:scale-105 active:scale-95 transition-all cursor-pointer group"
              aria-label="Download Game"
            >
              <Download class="w-4.5 h-4.5 group-hover:-translate-y-0.5 transition-transform" />
            </button>

            <!-- cURL Icon Button -->
            <button
              @click="gameStore.copyCurl(game)"
              title="Copy cURL Command"
              class="flex items-center justify-center p-3 rounded-full bg-white/10 hover:bg-white/15 text-slate-200 hover:text-white border border-white/15 shadow-lg backdrop-blur-md hover:scale-105 active:scale-95 transition-all cursor-pointer group"
              aria-label="Copy cURL command"
            >
              <Terminal class="w-4.5 h-4.5 text-indigo-400 group-hover:scale-110 transition-transform" />
            </button>

            <!-- Google Drive Icon Button -->
            <button
              @click="gameStore.openInDrive(game)"
              title="Open in Google Drive Web"
              class="flex items-center justify-center p-3 rounded-full bg-white/10 hover:bg-white/15 text-slate-200 hover:text-white border border-white/15 shadow-lg backdrop-blur-md hover:scale-105 active:scale-95 transition-all cursor-pointer group"
              aria-label="Open in Google Drive"
            >
              <ExternalLink class="w-4.5 h-4.5 text-slate-400 group-hover:text-white transition-colors" />
            </button>
          </div>

          <!-- Overview Synopsis -->
          <div class="pt-6 space-y-2 border-t border-white/10 text-left">
            <h2 class="text-xs font-bold uppercase tracking-wider text-slate-400">
              Overview
            </h2>
            <p class="text-sm sm:text-base text-slate-300/90 leading-relaxed whitespace-pre-line max-w-4xl font-normal">
              {{ game.description || 'No overview synopsis available for this game.' }}
            </p>
          </div>

          <!-- Specifications Table -->
          <GameSpecsTable :game="game" />

        </div>

      </div>

    </div>

  </div>
</template>
