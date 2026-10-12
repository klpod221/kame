<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Gamepad2, Plus, Loader2 } from 'lucide-vue-next'
import Navbar from './components/Navbar.vue'
import GameCard from './components/GameCard.vue'
import GameDetailModal from './components/GameDetailModal.vue'
import AddGameModal from './components/AddGameModal.vue'
import { useAuth } from './composables/useAuth'
import { useGames } from './composables/useGames'
import type { Game } from './types'

const { fetchConfig, isAdmin } = useAuth()
const { games, total, genres, isLoading, loadGames, loadGenres } = useGames()

const searchQuery = ref('')
const selectedGenre = ref('')
const selectedGame = ref<Game | null>(null)
const isAddModalOpen = ref(false)

let debounceTimer: any = null

const triggerSearch = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    loadGames(searchQuery.value, selectedGenre.value)
  }, 250)
}

watch(searchQuery, () => triggerSearch())
watch(selectedGenre, () => loadGames(searchQuery.value, selectedGenre.value))

onMounted(async () => {
  await fetchConfig()
  await Promise.all([loadGames(), loadGenres()])
})

const handleGameSelect = (game: Game) => {
  selectedGame.value = game
}

const handleQuickDownload = (game: Game) => {
  selectedGame.value = game
}
</script>

<template>
  <div class="min-h-screen flex flex-col bg-slate-950 text-slate-100">
    <!-- Navbar -->
    <Navbar
      v-model:searchQuery="searchQuery"
      v-model:selectedGenre="selectedGenre"
      :genres="genres"
      @openAddModal="isAddModalOpen = true"
    />

    <!-- Main Container -->
    <main class="flex-1 max-w-7xl w-full mx-auto px-4 lg:px-8 py-6 flex flex-col gap-6">
      
      <!-- Genre Filter Pills -->
      <div class="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-none">
        <button
          @click="selectedGenre = ''"
          class="shrink-0 px-3.5 py-1.5 rounded-xl text-xs font-medium transition-all cursor-pointer"
          :class="selectedGenre === '' 
            ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30' 
            : 'bg-slate-900/80 hover:bg-slate-800 text-slate-400 hover:text-slate-200 border border-slate-800'"
        >
          All Games
        </button>

        <button
          v-for="genre in genres"
          :key="genre"
          @click="selectedGenre = selectedGenre === genre ? '' : genre"
          class="shrink-0 px-3.5 py-1.5 rounded-xl text-xs font-medium transition-all cursor-pointer"
          :class="selectedGenre === genre 
            ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30' 
            : 'bg-slate-900/80 hover:bg-slate-800 text-slate-400 hover:text-slate-200 border border-slate-800'"
        >
          {{ genre }}
        </button>
      </div>

      <!-- Catalog Header Count -->
      <div class="flex items-center justify-between text-xs text-slate-400 px-1">
        <div class="flex items-center gap-1.5">
          <Gamepad2 class="w-4 h-4 text-indigo-400" />
          <span>{{ total }} {{ total === 1 ? 'game' : 'games' }} available</span>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center py-20 gap-3 text-slate-500">
        <Loader2 class="w-8 h-8 animate-spin text-indigo-500" />
        <span class="text-sm">Loading catalog...</span>
      </div>

      <!-- Empty State -->
      <div 
        v-else-if="games.length === 0" 
        class="flex-1 flex flex-col items-center justify-center py-20 px-4 rounded-3xl bg-slate-900/30 border border-dashed border-slate-800 text-center"
      >
        <div class="w-14 h-14 rounded-2xl bg-slate-800/80 flex items-center justify-center text-slate-400 mb-4">
          <Gamepad2 class="w-7 h-7" />
        </div>
        <h3 class="text-base font-semibold text-slate-200">No games found</h3>
        <p class="text-xs text-slate-400 mt-1 max-w-sm">
          {{ searchQuery ? 'Try adjusting your search keywords or clearing genre filters.' : 'The catalog is currently empty. Admins can add games by clicking "+ Add Game".' }}
        </p>
        <button
          v-if="isAdmin && !searchQuery"
          @click="isAddModalOpen = true"
          class="mt-5 flex items-center gap-2 px-4 py-2 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white text-xs font-medium transition-all shadow-lg shadow-indigo-600/20 cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>Add First Game</span>
        </button>
      </div>

      <!-- Games Grid -->
      <div 
        v-else 
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 sm:gap-5"
      >
        <GameCard
          v-for="game in games"
          :key="game.id"
          :game="game"
          @select="handleGameSelect"
          @download="handleQuickDownload"
        />
      </div>

    </main>

    <!-- Modals -->
    <GameDetailModal
      :game="selectedGame"
      @close="selectedGame = null"
      @deleted="loadGames(searchQuery, selectedGenre)"
    />

    <AddGameModal
      v-if="isAddModalOpen"
      @close="isAddModalOpen = false"
      @created="loadGames(searchQuery, selectedGenre); loadGenres()"
    />
  </div>
</template>
