<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Gamepad2, 
  LogOut, 
  ShieldCheck, 
  X, 
  Layers, 
  ChevronDown,
  ChevronRight
} from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth.store'
import { useGameStore } from '../../stores/game.store'

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'closeMobileSidebar'): void
}>()

const router = useRouter()
const authStore = useAuthStore()
const gameStore = useGameStore()

const isGenresExpanded = ref(true)

const selectGenre = (genre = '') => {
  gameStore.selectedGenre = genre
  emit('closeMobileSidebar')
  if (router.currentRoute.value.path !== '/') {
    router.push('/')
  }
}

const goToAdmin = () => {
  emit('closeMobileSidebar')
  router.push('/admin')
}
</script>

<template>
  <!-- Mobile Backdrop Overlay -->
  <div 
    v-if="isOpen" 
    class="fixed inset-0 z-40 bg-black/80 backdrop-blur-md lg:hidden transition-opacity"
    @click="emit('closeMobileSidebar')"
  ></div>

  <!-- Sidebar Container -->
  <aside 
    class="fixed top-0 bottom-0 left-0 z-50 w-64 bg-[#070a0f]/95 border-r border-white/5 backdrop-blur-xl flex flex-col justify-between transition-transform duration-300 ease-in-out lg:translate-x-0"
    :class="isOpen ? 'translate-x-0' : '-translate-x-full'"
  >
    <!-- Top Branding & Navigation -->
    <div class="flex flex-col flex-1 overflow-y-auto px-4 py-5 scrollbar-none">
      
      <!-- Brand Logo Header with transparent logo -->
      <div class="flex items-center justify-between pb-6 px-1 border-b border-white/5">
        <router-link to="/" class="flex items-center gap-3 group">
          <img 
            src="/kame-48.png" 
            alt="Kame" 
            class="w-9 h-9 object-contain shrink-0 group-hover:scale-105 transition-transform drop-shadow-[0_0_12px_rgba(99,102,241,0.3)]" 
          />
          <div class="flex items-center gap-2">
            <span class="text-xl font-black tracking-tight text-white group-hover:text-indigo-400 transition-colors">
              Kame
            </span>
            <span class="text-[9px] font-extrabold px-1.5 py-0.5 rounded bg-indigo-500/20 text-indigo-300 border border-indigo-500/30 uppercase tracking-widest">
              vault
            </span>
          </div>
        </router-link>

        <!-- Mobile Close Button -->
        <button 
          @click="emit('closeMobileSidebar')"
          class="p-1.5 rounded-lg text-slate-400 hover:text-white lg:hidden cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Main Navigation Menu -->
      <nav class="space-y-1.5 pt-5">
        
        <!-- All Games (Library) -->
        <button
          @click="selectGenre('')"
          class="w-full flex items-center justify-between px-3.5 py-2.5 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left"
          :class="gameStore.selectedGenre === '' 
            ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30' 
            : 'text-slate-400 hover:text-white hover:bg-white/5'"
        >
          <div class="flex items-center gap-3">
            <Gamepad2 class="w-4 h-4 shrink-0" />
            <span>Game Vault</span>
          </div>
          <span 
            class="text-[11px] px-2 py-0.5 rounded-full font-bold"
            :class="gameStore.selectedGenre === '' ? 'bg-indigo-700/60 text-white' : 'bg-white/5 text-slate-400'"
          >
            {{ gameStore.total }}
          </span>
        </button>

      </nav>

      <!-- GENRES SECTION (Collapsible) -->
      <div class="mt-6 pt-5 border-t border-white/5">
        <button
          @click="isGenresExpanded = !isGenresExpanded"
          class="w-full flex items-center justify-between px-3 mb-2.5 text-[10px] font-bold uppercase tracking-wider text-slate-400 hover:text-slate-200 transition-colors cursor-pointer select-none"
        >
          <div class="flex items-center gap-2">
            <Layers class="w-3.5 h-3.5 text-indigo-400" />
            <span>Genres</span>
          </div>
          <component 
            :is="isGenresExpanded ? ChevronDown : ChevronRight" 
            class="w-3.5 h-3.5 text-slate-500" 
          />
        </button>

        <div v-show="isGenresExpanded" class="space-y-1">
          <button
            v-for="genre in gameStore.genres"
            :key="genre"
            @click="selectGenre(genre)"
            class="w-full flex items-center justify-between px-3.5 py-2 rounded-xl text-xs font-medium transition-colors cursor-pointer text-left"
            :class="gameStore.selectedGenre === genre 
              ? 'bg-indigo-600 text-white font-semibold shadow-md shadow-indigo-600/20' 
              : 'text-slate-400 hover:text-slate-200 hover:bg-white/5'"
          >
            <span class="truncate">{{ genre }}</span>
          </button>
        </div>
      </div>

    </div>

    <!-- Bottom User Section & Admin Entry -->
    <div class="p-4 border-t border-white/5 bg-[#05070b]/90">
      
      <!-- Admin Console Button (Separated from client browsing) -->
      <button
        v-if="authStore.isAdmin"
        @click="goToAdmin"
        class="w-full flex items-center justify-between px-3.5 py-2.5 rounded-xl bg-purple-950/40 hover:bg-purple-900/60 border border-purple-500/30 text-purple-200 text-xs font-bold transition-all shadow-sm cursor-pointer group mb-3"
      >
        <div class="flex items-center gap-2.5">
          <ShieldCheck class="w-4 h-4 text-purple-400 group-hover:scale-110 transition-transform" />
          <span>Admin Console</span>
        </div>
        <ChevronRight class="w-3.5 h-3.5 text-purple-400/70 group-hover:translate-x-0.5 transition-transform" />
      </button>

      <!-- User Profile Card -->
      <div v-if="authStore.currentUser" class="flex items-center justify-between gap-3 p-2 rounded-xl bg-white/3 border border-white/5">
        <div class="flex items-center gap-2.5 min-w-0">
          <img
            v-if="authStore.currentUser.picture"
            :src="authStore.currentUser.picture"
            :alt="authStore.currentUser.name"
            class="w-8 h-8 rounded-full border border-indigo-500/40 object-cover shrink-0"
          />
          <div class="min-w-0">
            <p class="text-xs font-bold text-white truncate flex items-center gap-1">
              {{ authStore.currentUser.name }}
              <ShieldCheck v-if="authStore.isAdmin" class="w-3 h-3 text-indigo-400 shrink-0" />
            </p>
            <p class="text-[10px] text-slate-400 truncate">{{ authStore.currentUser.email }}</p>
          </div>
        </div>

        <!-- Logout Action -->
        <button
          @click="authStore.logout"
          title="Sign out"
          class="p-1.5 text-slate-400 hover:text-rose-400 hover:bg-rose-500/10 rounded-lg transition-colors cursor-pointer shrink-0"
        >
          <LogOut class="w-4 h-4" />
        </button>
      </div>

      <!-- Footer Info -->
      <div class="mt-3 px-1 flex items-center justify-between text-[10px] text-slate-400">
        <span>Kame • by klpod221</span>
        <div class="flex items-center gap-1 text-emerald-400">
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
          <span>Drive Stream</span>
        </div>
      </div>

    </div>
  </aside>
</template>
