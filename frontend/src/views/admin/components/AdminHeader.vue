<script setup lang="ts">
import { Menu, Gamepad2, Users, Server } from 'lucide-vue-next'

const props = defineProps<{
  activeTab: 'games' | 'users' | 'system'
}>()

const emit = defineEmits<{
  (e: 'toggleSidebar'): void
}>()

const getTabInfo = () => {
  switch (props.activeTab) {
    case 'games':
      return { title: 'Games Manager', subtitle: 'Catalog, metadata enrichment & Steam assets', icon: Gamepad2 }
    case 'users':
      return { title: 'Users & Permissions', subtitle: 'Manage accounts and administrative privileges', icon: Users }
    case 'system':
      return { title: 'System Diagnostics', subtitle: 'Runtime status, database metrics & environment', icon: Server }
  }
}
</script>

<template>
  <header class="sticky top-0 z-30 w-full bg-[#070a0f]/90 backdrop-blur-xl border-b border-white/10 px-4 sm:px-6 py-3 transition-colors shrink-0">
    <div class="flex items-center justify-between gap-4 w-full">
      
      <!-- Left: Mobile Menu Toggle & Page Context -->
      <div class="flex items-center gap-3">
        <!-- Mobile Menu Toggle Button -->
        <button
          @click="emit('toggleSidebar')"
          class="p-2 rounded-xl text-slate-300 hover:text-white hover:bg-white/10 border border-white/10 lg:hidden cursor-pointer shrink-0 transition-colors"
          aria-label="Toggle Admin Navigation"
        >
          <Menu class="w-5 h-5" />
        </button>

        <!-- Current Tab Title Badge -->
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400 shrink-0">
            <component :is="getTabInfo().icon" class="w-4 h-4" />
          </div>
          <div>
            <h1 class="text-sm font-black text-white tracking-tight flex items-center gap-2">
              <span>{{ getTabInfo().title }}</span>
            </h1>
            <p class="text-[10px] text-slate-400 hidden sm:block">
              {{ getTabInfo().subtitle }}
            </p>
          </div>
        </div>
      </div>

    </div>
  </header>
</template>
