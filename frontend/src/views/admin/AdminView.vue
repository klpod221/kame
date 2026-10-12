<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAdminStore } from '../../stores/admin.store'
import { useGameStore } from '../../stores/game.store'
import AdminHeader from './components/AdminHeader.vue'
import GamesTab from './components/GamesTab.vue'
import UsersTab from './components/UsersTab.vue'
import SystemTab from './components/SystemTab.vue'
import AddGameModal from './components/AddGameModal.vue'
import EditGameModal from './components/EditGameModal.vue'
import type { Game } from '../../types'

const activeTab = ref<'games' | 'users' | 'system'>('games')
const isAddModalOpen = ref(false)
const editingGame = ref<Game | null>(null)

const adminStore = useAdminStore()
const gameStore = useGameStore()

onMounted(async () => {
  await Promise.all([
    adminStore.fetchStats(),
    adminStore.fetchUsers(),
    gameStore.fetchGames(),
  ])
})
</script>

<template>
  <div class="min-h-screen bg-[#070a0f] text-slate-100 flex flex-col animate-fadeIn pb-16">
    
    <!-- Dedicated Admin Top Header with Tab Switcher -->
    <AdminHeader v-model:activeTab="activeTab" />

    <!-- Tab View Container -->
    <main class="flex-1 p-4 sm:p-6 lg:p-8 max-w-7xl w-full mx-auto space-y-6">
      <GamesTab 
        v-if="activeTab === 'games'" 
        @add="isAddModalOpen = true" 
        @edit="editingGame = $event" 
      />
      
      <UsersTab 
        v-else-if="activeTab === 'users'" 
      />
      
      <SystemTab 
        v-else-if="activeTab === 'system'" 
      />
    </main>

    <!-- Add Game Modal -->
    <AddGameModal
      v-if="isAddModalOpen"
      @close="isAddModalOpen = false"
    />

    <!-- Edit Game Modal -->
    <EditGameModal
      v-if="editingGame"
      :game="editingGame"
      @close="editingGame = null"
    />

  </div>
</template>
