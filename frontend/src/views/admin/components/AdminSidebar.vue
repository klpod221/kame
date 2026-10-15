<script setup lang="ts">
import { useRouter } from 'vue-router'
import { 
  Gamepad2, 
  Users, 
  Server, 
  LogOut, 
  ShieldCheck, 
  X, 
  ArrowLeft,
  ChevronRight
} from 'lucide-vue-next'
import { useAuthStore } from '../../../stores/auth.store'
import { useAdminStore } from '../../../stores/admin.store'
import { useGameStore } from '../../../stores/game.store'

defineProps<{
  isOpen: boolean
  activeTab: 'games' | 'users' | 'system'
}>()

const emit = defineEmits<{
  (e: 'update:activeTab', tab: 'games' | 'users' | 'system'): void
  (e: 'closeMobileSidebar'): void
}>()

const router = useRouter()
const authStore = useAuthStore()
const adminStore = useAdminStore()
const gameStore = useGameStore()

const selectTab = (tab: 'games' | 'users' | 'system') => {
  emit('update:activeTab', tab)
  emit('closeMobileSidebar')
}

const goToLibrary = () => {
  emit('closeMobileSidebar')
  router.push('/')
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
    class="fixed top-0 bottom-0 left-0 z-50 w-64 h-screen max-h-screen bg-[#070a0f]/95 border-r border-white/5 backdrop-blur-xl flex flex-col justify-between transition-transform duration-300 ease-in-out lg:translate-x-0 overflow-hidden"
    :class="isOpen ? 'translate-x-0' : '-translate-x-full'"
  >
    <!-- 1. Top Branding & Navigation -->
    <div class="shrink-0 px-4 pt-5 pb-3">
      <!-- Brand Logo Header -->
      <div class="flex items-center justify-between pb-4 px-1 border-b border-white/5">
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
            <span class="text-[9px] font-extrabold px-1.5 py-0.5 rounded bg-purple-500/20 text-purple-300 border border-purple-500/30 uppercase tracking-widest">
              admin
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

      <!-- Navigation Menu (Tabs transferred from top-right) -->
      <nav class="space-y-1.5 pt-4">
        <!-- 1. Games Manager -->
        <button
          @click="selectTab('games')"
          class="w-full flex items-center justify-between px-3.5 py-2.5 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left"
          :class="activeTab === 'games' 
            ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30 font-bold' 
            : 'text-slate-400 hover:text-white hover:bg-white/5'"
        >
          <div class="flex items-center gap-3">
            <Gamepad2 class="w-4 h-4 shrink-0" />
            <span>Games Manager</span>
          </div>
          <span 
            class="text-[11px] px-2 py-0.5 rounded-full font-bold"
            :class="activeTab === 'games' ? 'bg-indigo-700/60 text-white' : 'bg-white/5 text-slate-400'"
          >
            {{ adminStore.stats?.total_games ?? gameStore.games.length }}
          </span>
        </button>

        <!-- 2. Users -->
        <button
          @click="selectTab('users')"
          class="w-full flex items-center justify-between px-3.5 py-2.5 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left"
          :class="activeTab === 'users' 
            ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30 font-bold' 
            : 'text-slate-400 hover:text-white hover:bg-white/5'"
        >
          <div class="flex items-center gap-3">
            <Users class="w-4 h-4 shrink-0" />
            <span>Users</span>
          </div>
          <span 
            class="text-[11px] px-2 py-0.5 rounded-full font-bold"
            :class="activeTab === 'users' ? 'bg-indigo-700/60 text-white' : 'bg-white/5 text-slate-400'"
          >
            {{ adminStore.stats?.total_users ?? adminStore.users.length }}
          </span>
        </button>

        <!-- 3. System Info -->
        <button
          @click="selectTab('system')"
          class="w-full flex items-center justify-between px-3.5 py-2.5 rounded-xl text-xs font-semibold transition-all cursor-pointer text-left"
          :class="activeTab === 'system' 
            ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/30 font-bold' 
            : 'text-slate-400 hover:text-white hover:bg-white/5'"
        >
          <div class="flex items-center gap-3">
            <Server class="w-4 h-4 shrink-0" />
            <span>System Info</span>
          </div>
        </button>
      </nav>
    </div>

    <!-- 2. Middle Spacer -->
    <div class="flex-1 min-h-0"></div>

    <!-- 3. Bottom User Section & Back to Library (Fixed at bottom) -->
    <div class="shrink-0 p-4 border-t border-white/5 bg-[#05070b]/95">
      
      <!-- Back to Library Button -->
      <button
        @click="goToLibrary"
        class="w-full flex items-center justify-between px-3.5 py-2.5 rounded-xl bg-indigo-950/40 hover:bg-indigo-900/60 border border-indigo-500/30 text-indigo-200 text-xs font-bold transition-all shadow-sm cursor-pointer group mb-3"
      >
        <div class="flex items-center gap-2.5">
          <ArrowLeft class="w-4 h-4 text-indigo-400 group-hover:-translate-x-0.5 transition-transform" />
          <span>Back to Library</span>
        </div>
        <ChevronRight class="w-3.5 h-3.5 text-indigo-400/70 opacity-0 group-hover:opacity-100 transition-opacity" />
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
        <span>Kame • Admin Console</span>
        <div class="flex items-center gap-1 text-emerald-400">
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
          <span>Online</span>
        </div>
      </div>

    </div>
  </aside>
</template>
