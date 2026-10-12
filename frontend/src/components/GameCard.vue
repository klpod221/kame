<script setup lang="ts">
import { Star, Download, Eye, Calendar } from 'lucide-vue-next'
import type { Game } from '../types'

defineProps<{
  game: Game
}>()

const emit = defineEmits<{
  (e: 'select', game: Game): void
  (e: 'download', game: Game): void
}>()
</script>

<template>
  <div 
    @click="emit('select', game)"
    class="group relative flex flex-col rounded-2xl bg-slate-900/60 border border-slate-800/80 overflow-hidden hover:border-indigo-500/50 hover:shadow-2xl hover:shadow-indigo-500/10 transition-all duration-300 cursor-pointer transform hover:-translate-y-1.5"
  >
    <!-- Cover Poster (Aspect Ratio 3:4) -->
    <div class="relative aspect-[3/4] w-full overflow-hidden bg-slate-950">
      <img
        v-if="game.cover_url"
        :src="game.cover_url"
        :alt="game.title"
        loading="lazy"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-slate-600 bg-slate-900">
        <span class="text-sm font-medium">No Cover Art</span>
      </div>

      <!-- Gradient Overlay on Poster -->
      <div class="absolute inset-0 bg-gradient-to-t from-slate-950 via-slate-950/20 to-transparent opacity-80 group-hover:opacity-60 transition-opacity"></div>

      <!-- Top Badges -->
      <div class="absolute top-2.5 inset-x-2.5 flex items-center justify-between">
        <!-- Release Year -->
        <span 
          v-if="game.release_year" 
          class="flex items-center gap-1 text-[11px] font-medium px-2 py-0.5 rounded-md backdrop-blur-md bg-slate-900/80 text-slate-300 border border-slate-700/50"
        >
          <Calendar class="w-3 h-3 text-slate-400" />
          {{ game.release_year }}
        </span>
        <span v-else></span>

        <!-- Rating -->
        <span 
          v-if="game.rating > 0" 
          class="flex items-center gap-1 text-[11px] font-semibold px-2 py-0.5 rounded-md backdrop-blur-md bg-amber-500/20 text-amber-300 border border-amber-500/30"
        >
          <Star class="w-3 h-3 fill-amber-400 text-amber-400" />
          {{ game.rating.toFixed(1) }}
        </span>
      </div>

      <!-- Quick Action Overlay on Hover -->
      <div class="absolute inset-0 flex items-center justify-center gap-3 opacity-0 group-hover:opacity-100 backdrop-blur-[2px] bg-slate-950/40 transition-all duration-200">
        <button
          @click.stop="emit('download', game)"
          title="Direct Download"
          class="p-3 rounded-full bg-indigo-600 hover:bg-indigo-500 text-white shadow-xl shadow-indigo-600/40 transform hover:scale-110 transition-all cursor-pointer"
        >
          <Download class="w-5 h-5" />
        </button>
        <button
          @click.stop="emit('select', game)"
          title="Game Info"
          class="p-3 rounded-full bg-slate-800/90 hover:bg-slate-700 text-slate-200 border border-slate-600/50 transform hover:scale-110 transition-all cursor-pointer"
        >
          <Eye class="w-5 h-5" />
        </button>
      </div>
    </div>

    <!-- Metadata Bottom Bar -->
    <div class="p-3.5 flex flex-col flex-1 justify-between gap-2">
      <div>
        <h3 class="text-sm font-semibold text-slate-100 line-clamp-1 group-hover:text-indigo-400 transition-colors" :title="game.title">
          {{ game.title }}
        </h3>

        <!-- Genre Tags -->
        <div class="flex flex-wrap gap-1 mt-1.5">
          <span
            v-for="genre in (game.genres || []).slice(0, 2)"
            :key="genre"
            class="text-[10px] font-medium px-1.5 py-0.5 rounded bg-slate-800 text-slate-400"
          >
            {{ genre }}
          </span>
          <span 
            v-if="(game.genres || []).length > 2" 
            class="text-[10px] font-medium px-1 py-0.5 rounded bg-slate-800/50 text-slate-500"
          >
            +{{ game.genres.length - 2 }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
