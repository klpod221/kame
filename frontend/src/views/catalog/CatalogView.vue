<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Gamepad2, Loader2, ArrowUpDown, Sparkles } from 'lucide-vue-next'
import Sidebar from '../../components/layout/Sidebar.vue'
import HeaderBar from '../../components/layout/HeaderBar.vue'
import GameCard from './components/GameCard.vue'
import { useGameStore } from '../../stores/game.store'
import type { Game } from '../../types'

const router = useRouter()
const route = useRoute()
const gameStore = useGameStore()

const isMobileSidebarOpen = ref(false)

// Sync route genre query with store
watch(() => route.query.genre, (qGenre) => {
  const g = (qGenre as string) || ''
  if (g !== gameStore.selectedGenre) {
    gameStore.selectedGenre = g
    gameStore.fetchGames()
  }
})

// Sync store genre with route query
watch(() => gameStore.selectedGenre, (newGenre) => {
  if (newGenre) {
    router.replace({ query: { ...route.query, genre: newGenre } })
  } else {
    const q = { ...route.query }
    delete q.genre
    router.replace({ query: q })
  }
  gameStore.fetchGames()
})

onMounted(async () => {
  if (route.query.genre) {
    gameStore.selectedGenre = route.query.genre as string
  }
  await Promise.all([
    gameStore.fetchGames(),
    gameStore.fetchGenres(),
  ])
})

const handleGameSelect = (game: Game) => {
  router.push(`/games/${game.id}`)
}

const handleGenreFilter = (genre: string) => {
  gameStore.selectedGenre = gameStore.selectedGenre === genre ? '' : genre
}
</script>

<template>
  <div class="min-h-screen bg-[#070a0f] text-slate-100 flex animate-fadeIn">
    
    <!-- Left Navigation Sidebar -->
    <Sidebar
      :isOpen="isMobileSidebarOpen"
      @closeMobileSidebar="isMobileSidebarOpen = false"
    />

    <!-- Main Content Layout (Offset for Sidebar on desktop) -->
    <div class="flex-1 flex flex-col min-w-0 lg:pl-64">
      
      <!-- Top Header Bar with Debounced Search -->
      <HeaderBar
        @toggleSidebar="isMobileSidebarOpen = !isMobileSidebarOpen"
      />

      <!-- Catalog Main Area -->
      <main class="flex-1 p-4 sm:p-6 lg:p-8 flex flex-col gap-6 max-w-7xl w-full mx-auto">
        
        <!-- Horizontal Genre Pills Bar (Quick Filter Ribbon) -->
        <div class="flex items-center gap-2 overflow-x-auto pb-1.5 scrollbar-none select-none">
          <button
            @click="handleGenreFilter('')"
            class="px-3.5 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition-all cursor-pointer border"
            :class="!gameStore.selectedGenre 
              ? 'bg-indigo-600 border-indigo-500 text-white shadow-lg shadow-indigo-600/30' 
              : 'bg-white/5 border-white/5 text-slate-400 hover:text-white hover:bg-white/10'"
          >
            All Vault Items
          </button>

          <button
            v-for="genre in gameStore.genres"
            :key="genre"
            @click="handleGenreFilter(genre)"
            class="px-3.5 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition-all cursor-pointer border flex items-center gap-1.5"
            :class="gameStore.selectedGenre === genre 
              ? 'bg-indigo-600 border-indigo-500 text-white shadow-lg shadow-indigo-600/30' 
              : 'bg-white/5 border-white/5 text-slate-400 hover:text-white hover:bg-white/10'"
          >
            <span>{{ genre }}</span>
          </button>
        </div>

        <!-- Section Header with Count and Sort Selector -->
        <div class="flex items-center justify-between pt-1 pb-2 border-b border-white/5">
          <div class="flex items-center gap-3">
            <h2 class="text-lg sm:text-xl font-black text-white tracking-tight flex items-center gap-2">
              <Sparkles v-if="!gameStore.selectedGenre" class="w-4 h-4 text-indigo-400" />
              <span>{{ gameStore.selectedGenre ? gameStore.selectedGenre : 'Game Vault' }}</span>
            </h2>
            <span class="text-xs font-bold px-2.5 py-0.5 rounded-full bg-white/5 border border-white/10 text-slate-400">
              {{ gameStore.total }} items
            </span>
          </div>

          <!-- Sort Select Component -->
          <div class="flex items-center gap-2">
            <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-white/5 border border-white/10 text-xs text-slate-400 backdrop-blur-md">
              <ArrowUpDown class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
              <select
                v-model="gameStore.sortBy"
                class="bg-transparent text-slate-200 font-semibold focus:outline-none cursor-pointer text-xs"
              >
                <option value="rating" class="bg-slate-900 text-slate-200">Rating (Highest)</option>
                <option value="newest" class="bg-slate-900 text-slate-200">Release Year (Newest)</option>
                <option value="title" class="bg-slate-900 text-slate-200">Title (A - Z)</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Loading State -->
        <div v-if="gameStore.isLoading" class="flex-1 flex flex-col items-center justify-center py-28 gap-3 text-slate-500">
          <Loader2 class="w-8 h-8 animate-spin text-indigo-500" />
          <span class="text-xs font-medium">Loading vault games...</span>
        </div>

        <!-- Empty State -->
        <div 
          v-else-if="gameStore.sortedGames.length === 0" 
          class="flex-1 flex flex-col items-center justify-center py-24 px-4 rounded-3xl bg-white/2 border border-dashed border-white/10 text-center"
        >
          <div class="w-14 h-14 rounded-2xl bg-white/5 border border-white/10 flex items-center justify-center text-slate-500 mb-4 shadow-xl">
            <Gamepad2 class="w-7 h-7" />
          </div>
          <h2 class="text-base font-bold text-slate-200">No games found</h2>
          <p class="text-xs text-slate-400 mt-1 max-w-sm">
            {{ gameStore.searchQuery ? 'Try adjusting your search terms.' : 'No games in this section.' }}
          </p>
        </div>

        <!-- Jellyfin Cinema Poster Grid (Larger Cards) -->
        <div 
          v-else 
          class="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-5 sm:gap-6 lg:gap-7"
        >
          <GameCard
            v-for="game in gameStore.sortedGames"
            :key="game.id"
            :game="game"
            @select="handleGameSelect"
          />
        </div>

      </main>

    </div>

  </div>
</template>
