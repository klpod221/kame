<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft, Gamepad2, Users, Server } from 'lucide-vue-next'
import { useAdminStore } from '../../../stores/admin.store'
import { useGameStore } from '../../../stores/game.store'

defineProps<{
  activeTab: 'games' | 'users' | 'system'
}>()

const emit = defineEmits<{
  (e: 'update:activeTab', tab: 'games' | 'users' | 'system'): void
}>()

const router = useRouter()
const adminStore = useAdminStore()
const gameStore = useGameStore()

const goBackToLibrary = () => {
  router.push('/')
}
</script>

<template>
  <header class="sticky top-0 z-30 w-full bg-[#070a0f]/90 backdrop-blur-xl border-b border-white/5 px-4 sm:px-8 py-3 flex flex-wrap items-center justify-between gap-4 transition-colors">
    <div class="flex items-center gap-4">
      <!-- Back to Library -->
      <button
        @click="goBackToLibrary"
        class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-white/5 hover:bg-white/10 text-slate-300 hover:text-white border border-white/10 text-xs font-semibold transition-all shadow-sm active:scale-95 cursor-pointer"
      >
        <ArrowLeft class="w-4 h-4 text-indigo-400" />
        <span>Library</span>
      </button>

      <div class="h-4 w-px bg-white/10 hidden sm:block"></div>

      <!-- Brand Title with transparent icon -->
      <div class="flex items-center gap-2.5">
        <img src="/kame-48.png" alt="Kame" class="w-7 h-7 object-contain" />
        <span class="text-sm font-black tracking-tight text-white hidden sm:inline">Kame</span>
        <span class="text-[10px] font-bold px-2 py-0.5 rounded-md bg-white/5 text-slate-300 border border-white/10 tracking-wide">
          Admin Console
        </span>
      </div>
    </div>

    <!-- Navigation Tabs (Sleek Pills) -->
    <nav class="flex items-center gap-1.5 bg-white/3 p-1 rounded-xl border border-white/5 overflow-x-auto max-w-full">
      <button
        @click="emit('update:activeTab', 'games')"
        class="flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap"
        :class="activeTab === 'games' ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-white/5'"
      >
        <Gamepad2 class="w-3.5 h-3.5" />
        <span>Games Manager</span>
        <span class="text-[10px] px-1.5 py-0.2 rounded-full font-bold" :class="activeTab === 'games' ? 'bg-indigo-700 text-white' : 'bg-white/10 text-slate-400'">
          {{ adminStore.stats?.total_games ?? gameStore.games.length }}
        </span>
      </button>

      <button
        @click="emit('update:activeTab', 'users')"
        class="flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap"
        :class="activeTab === 'users' ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-white/5'"
      >
        <Users class="w-3.5 h-3.5" />
        <span>Users</span>
        <span class="text-[10px] px-1.5 py-0.2 rounded-full font-bold" :class="activeTab === 'users' ? 'bg-indigo-700 text-white' : 'bg-white/10 text-slate-400'">
          {{ adminStore.stats?.total_users ?? adminStore.users.length }}
        </span>
      </button>

      <button
        @click="emit('update:activeTab', 'system')"
        class="flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap"
        :class="activeTab === 'system' ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30' : 'text-slate-400 hover:text-white hover:bg-white/5'"
      >
        <Server class="w-3.5 h-3.5" />
        <span>System Info</span>
      </button>
    </nav>
  </header>
</template>
