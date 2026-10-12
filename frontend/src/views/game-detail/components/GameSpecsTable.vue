<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Check, Copy, HardDrive, Calendar, Star, Monitor, Tag, Clock } from 'lucide-vue-next'
import type { Game } from '../../../types'
import { formatBytes, formatDate } from '../../../utils/format'

const props = defineProps<{
  game: Game
}>()

const router = useRouter()
const copiedId = ref(false)

const copyDriveId = async () => {
  await navigator.clipboard.writeText(props.game.drive_file_id)
  copiedId.value = true
  setTimeout(() => {
    copiedId.value = false
  }, 2000)
}

const filterGenre = (genre: string) => {
  router.push({ path: '/', query: { genre } })
}
</script>

<template>
  <div class="pt-8 border-t border-white/10 text-left space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-xs font-bold uppercase tracking-wider text-slate-400 flex items-center gap-2">
        <Monitor class="w-3.5 h-3.5 text-indigo-400" />
        <span>Media & Vault Information</span>
      </h2>
    </div>

    <!-- Specs Matrix Grid (Jellyfin Media Info Style) -->
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
      
      <!-- Platform -->
      <div class="p-3.5 rounded-xl bg-white/3 border border-white/5 space-y-1">
        <div class="flex items-center gap-1.5 text-slate-400 text-[11px] font-medium">
          <Monitor class="w-3.5 h-3.5 text-sky-400" />
          <span>Platform</span>
        </div>
        <p class="text-xs sm:text-sm font-bold text-slate-200">PC (Windows)</p>
      </div>

      <!-- Release Year -->
      <div class="p-3.5 rounded-xl bg-white/3 border border-white/5 space-y-1">
        <div class="flex items-center gap-1.5 text-slate-400 text-[11px] font-medium">
          <Calendar class="w-3.5 h-3.5 text-indigo-400" />
          <span>Release Year</span>
        </div>
        <p class="text-xs sm:text-sm font-bold text-slate-200">{{ game.release_year || 'N/A' }}</p>
      </div>

      <!-- Storage Size -->
      <div class="p-3.5 rounded-xl bg-white/3 border border-white/5 space-y-1">
        <div class="flex items-center gap-1.5 text-slate-400 text-[11px] font-medium">
          <HardDrive class="w-3.5 h-3.5 text-emerald-400" />
          <span>Package Size</span>
        </div>
        <p class="text-xs sm:text-sm font-bold font-mono text-slate-200">
          {{ game.file_size && game.file_size > 0 ? formatBytes(game.file_size) : 'Vault Stream' }}
        </p>
      </div>

      <!-- Rating Score -->
      <div class="p-3.5 rounded-xl bg-white/3 border border-white/5 space-y-1">
        <div class="flex items-center gap-1.5 text-slate-400 text-[11px] font-medium">
          <Star class="w-3.5 h-3.5 text-amber-400" />
          <span>User Rating</span>
        </div>
        <p class="text-xs sm:text-sm font-bold text-slate-200 flex items-center gap-1">
          <span v-if="game.rating > 0">{{ game.rating.toFixed(1) }} / 5.0</span>
          <span v-else class="text-slate-500 font-normal">Not rated</span>
        </p>
      </div>

      <!-- Drive File ID (Full width or span 2) -->
      <div class="col-span-2 p-3.5 rounded-xl bg-white/3 border border-white/5 space-y-1">
        <div class="flex items-center gap-1.5 text-slate-400 text-[11px] font-medium">
          <HardDrive class="w-3.5 h-3.5 text-indigo-400" />
          <span>Google Drive File ID</span>
        </div>
        <button 
          @click="copyDriveId"
          title="Click to copy Drive File ID"
          class="font-mono text-xs text-indigo-300 hover:text-white flex items-center gap-2 cursor-pointer bg-black/40 px-2.5 py-1 rounded-lg border border-white/10 w-full justify-between transition-colors"
        >
          <span class="truncate">{{ game.drive_file_id }}</span>
          <span class="shrink-0 flex items-center gap-1 text-[10px] text-slate-400">
            <Check v-if="copiedId" class="w-3.5 h-3.5 text-emerald-400" />
            <Copy v-else class="w-3.5 h-3.5 text-slate-400" />
            <span>{{ copiedId ? 'Copied' : 'Copy' }}</span>
          </span>
        </button>
      </div>

      <!-- Added Date -->
      <div class="col-span-2 sm:col-span-1 lg:col-span-2 p-3.5 rounded-xl bg-white/3 border border-white/5 space-y-1">
        <div class="flex items-center gap-1.5 text-slate-400 text-[11px] font-medium">
          <Clock class="w-3.5 h-3.5 text-slate-400" />
          <span>Added to Library</span>
        </div>
        <p class="text-xs sm:text-sm font-medium text-slate-300">
          {{ formatDate(game.created_at) }}
        </p>
      </div>

    </div>

    <!-- Genres Tag Cloud -->
    <div class="pt-2 flex items-center gap-2 flex-wrap">
      <span class="text-[11px] font-medium text-slate-400 flex items-center gap-1">
        <Tag class="w-3 h-3 text-slate-500" />
        <span>Genres:</span>
      </span>
      <div class="flex flex-wrap gap-1.5">
        <button
          v-for="genre in (game.genres || [])"
          :key="genre"
          @click="filterGenre(genre)"
          class="px-2.5 py-1 rounded-lg bg-white/5 hover:bg-white/10 border border-white/10 text-slate-300 hover:text-white text-xs font-semibold cursor-pointer transition-colors"
        >
          {{ genre }}
        </button>
        <span v-if="!game.genres || game.genres.length === 0" class="text-slate-500 text-xs">Uncategorized</span>
      </div>
    </div>
  </div>
</template>
