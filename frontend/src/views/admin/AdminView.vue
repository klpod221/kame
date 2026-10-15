<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAdminStore } from '../../stores/admin.store'
import { useGameStore } from '../../stores/game.store'
import AdminSidebar from './components/AdminSidebar.vue'
import AdminHeader from './components/AdminHeader.vue'
import GamesTab from './components/GamesTab.vue'
import UsersTab from './components/UsersTab.vue'
import SystemTab from './components/SystemTab.vue'
import AddGameModal from './components/AddGameModal.vue'
import EditGameModal from './components/EditGameModal.vue'
import type { Game } from '../../types'

const activeTab = ref<'games' | 'users' | 'system'>('games')
const isMobileSidebarOpen = ref(false)
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
  <div class="h-screen bg-[#070a0f] text-slate-100 flex overflow-hidden animate-fadeIn">
    
    <!-- Left Admin Navigation Sidebar (Mirrors Library Sidebar Structure) -->
    <AdminSidebar
      :isOpen="isMobileSidebarOpen"
      :activeTab="activeTab"
      @update:activeTab="activeTab = $event; isMobileSidebarOpen = false"
      @closeMobileSidebar="isMobileSidebarOpen = false"
    />

    <!-- Main Content Layout (Offset for Sidebar on desktop) -->
    <div class="flex-1 flex flex-col min-w-0 h-full overflow-hidden lg:pl-64">
      
      <!-- Top Admin Header Bar with Mobile Toggle -->
      <AdminHeader
        :activeTab="activeTab"
        @toggleSidebar="isMobileSidebarOpen = !isMobileSidebarOpen"
      />

      <!-- Scrollable Tab Content Area -->
      <main class="flex-1 min-h-0 overflow-y-auto p-4 sm:p-6 lg:p-8 overscroll-contain">
        <div class="max-w-7xl w-full mx-auto pb-12">
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
        </div>
      </main>

    </div>

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
