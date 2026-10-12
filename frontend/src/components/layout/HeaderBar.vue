<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { Search, Menu, X } from 'lucide-vue-next'
import { debounce } from '../../utils/debounce'
import { useGameStore } from '../../stores/game.store'

const emit = defineEmits<{
  (e: 'toggleSidebar'): void
}>()

const gameStore = useGameStore()
const localSearch = ref(gameStore.searchQuery)

const debouncedSearch = debounce((query: string) => {
  gameStore.searchQuery = query
  gameStore.fetchGames()
}, 300)

watch(localSearch, (newVal) => {
  debouncedSearch(newVal)
})

watch(() => gameStore.searchQuery, (newVal) => {
  if (newVal !== localSearch.value) {
    localSearch.value = newVal
  }
})

const clearSearch = () => {
  localSearch.value = ''
  gameStore.searchQuery = ''
  gameStore.fetchGames()
}

// Global keyboard shortcut listener for Ctrl+K
const handleGlobalKeyDown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    const inputEl = document.getElementById('global-search-input') as HTMLInputElement
    if (inputEl) {
      inputEl.focus()
      inputEl.select()
    }
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeyDown)
})
</script>

<template>
  <header class="sticky top-0 z-30 w-full bg-[#070a0f]/90 backdrop-blur-xl border-b border-white/10 px-4 sm:px-6 py-3 transition-colors">
    <div class="flex items-center justify-between gap-4 w-full">
      
      <!-- Mobile Menu Toggle Button -->
      <button
        @click="emit('toggleSidebar')"
        class="p-2.5 rounded-xl text-slate-300 hover:text-white hover:bg-white/10 border border-white/10 lg:hidden cursor-pointer shrink-0 transition-colors"
        aria-label="Toggle Navigation"
      >
        <Menu class="w-5 h-5" />
      </button>

      <!-- Prominent Search Bar (Always visible & crisp) -->
      <div class="relative flex-1 max-w-2xl">
        <div class="relative flex items-center">
          <!-- Search Icon Prefix -->
          <div class="absolute left-3.5 flex items-center pointer-events-none text-slate-400">
            <Search class="w-4 h-4 text-indigo-400" />
          </div>

          <!-- Main Search Input -->
          <input
            id="global-search-input"
            v-model="localSearch"
            type="text"
            placeholder="Search games by title, genre, or keyword..."
            class="w-full pl-10 pr-20 py-2 sm:py-2.5 text-xs sm:text-sm text-white placeholder-slate-400 bg-white/5 hover:bg-white/8 focus:bg-white/10 border border-white/15 focus:border-indigo-500 rounded-xl transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 shadow-inner"
          />

          <!-- Clear Button & Ctrl+K Shortcut -->
          <div class="absolute right-3 flex items-center gap-1.5">
            <button
              v-if="localSearch"
              @click="clearSearch"
              class="p-1 text-slate-400 hover:text-white rounded-md transition-colors cursor-pointer"
              title="Clear search"
            >
              <X class="w-3.5 h-3.5" />
            </button>

            <span class="hidden sm:inline-block font-mono text-[10px] text-slate-400 bg-black/40 px-1.5 py-0.5 rounded border border-white/10 select-none">
              Ctrl K
            </span>
          </div>
        </div>
      </div>

      <!-- Right Spacer to balance layout on desktop -->
      <div class="hidden lg:block w-8"></div>

    </div>
  </header>
</template>
