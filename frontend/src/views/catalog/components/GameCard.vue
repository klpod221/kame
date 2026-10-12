<script setup lang="ts">
import { ref } from 'vue'
import { Star, Calendar, HardDrive } from 'lucide-vue-next'
import type { Game } from '../../../types'
import { formatBytes } from '../../../utils/format'

const props = defineProps<{
  game: Game
}>()

const emit = defineEmits<{
  (e: 'select', game: Game): void
}>()

const isImageLoaded = ref(false)
const hasImageError = ref(false)

const onImageLoad = () => {
  isImageLoaded.value = true
}

const onImageError = () => {
  hasImageError.value = true
}
</script>

<template>
  <div 
    @click="emit('select', game)"
    class="group relative flex flex-col rounded-2xl bg-slate-900/40 border border-white/5 overflow-hidden transition-all duration-300 ease-out cursor-pointer transform hover:-translate-y-2 hover:border-indigo-500/80 hover:shadow-[0_20px_40px_-12px_rgba(99,102,241,0.35)] backdrop-blur-md"
  >
    <!-- Jellyfin Cinema Poster Container (Aspect 2:3) -->
    <div class="relative aspect-2/3 w-full overflow-hidden bg-slate-950">
      
      <!-- Skeleton shimmer placeholder while loading -->
      <div 
        v-if="!isImageLoaded && !hasImageError && game.cover_url"
        class="absolute inset-0 animate-shimmer"
      ></div>

      <!-- Poster Image with Smooth Zoom -->
      <img
        v-if="game.cover_url && !hasImageError"
        :src="game.cover_url"
        :alt="game.title"
        loading="lazy"
        @load="onImageLoad"
        @error="onImageError"
        class="w-full h-full object-cover transition-transform duration-700 ease-out group-hover:scale-106"
        :class="isImageLoaded ? 'opacity-100' : 'opacity-0'"
      />

      <!-- Fallback when no cover art -->
      <div 
        v-else 
        class="w-full h-full flex flex-col items-center justify-center text-slate-600 bg-slate-950 gap-2 p-4 text-center"
      >
        <HardDrive class="w-12 h-12 text-slate-700 group-hover:text-indigo-400 transition-colors" />
        <span class="text-xs font-semibold text-slate-400 line-clamp-2">{{ game.title }}</span>
      </div>

      <!-- Cinema Gradient Vignette Overlay (Deeper at bottom for text readability) -->
      <div class="absolute inset-0 bg-linear-to-t from-slate-950 via-slate-950/25 to-transparent opacity-90 group-hover:opacity-75 transition-opacity duration-300"></div>

      <!-- Subtle Cinema Hover Reflection Sheen -->
      <div class="absolute inset-0 bg-linear-to-tr from-indigo-500/0 via-white/0 to-white/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"></div>

      <!-- Top Floating Badges (Rating & Year) -->
      <div class="absolute top-3 inset-x-3 flex items-center justify-between pointer-events-none z-10">
        <!-- Release Year -->
        <span 
          v-if="game.release_year" 
          class="flex items-center gap-1 text-[11px] font-bold px-2.5 py-1 rounded-lg backdrop-blur-md bg-slate-950/85 text-slate-300 border border-white/10 shadow-sm"
        >
          <Calendar class="w-3 h-3 text-indigo-400" />
          {{ game.release_year }}
        </span>
        <span v-else></span>

        <!-- Rating with Golden Star -->
        <span 
          v-if="game.rating > 0" 
          class="flex items-center gap-1 text-[11px] font-extrabold px-2.5 py-1 rounded-lg backdrop-blur-md bg-amber-500/20 text-amber-300 border border-amber-500/40 shadow-sm shadow-amber-500/10"
        >
          <Star class="w-3.5 h-3.5 fill-amber-400 text-amber-400" />
          {{ game.rating.toFixed(1) }}
        </span>
      </div>
    </div>

    <!-- Bottom Content Metadata -->
    <div class="p-4 flex flex-col justify-between flex-1 gap-2 bg-slate-950/60 transition-colors group-hover:bg-slate-950/80">
      <div>
        <h3 
          class="text-sm sm:text-base font-bold text-white line-clamp-1 group-hover:text-indigo-400 transition-colors"
          :title="game.title"
        >
          {{ game.title }}
        </h3>
        
        <p class="text-xs text-slate-400 line-clamp-1 mt-1 font-normal">
          {{ game.description || 'Cinema Vault Item' }}
        </p>
      </div>

      <!-- Genres Tags & File Size -->
      <div class="flex items-center justify-between text-xs text-slate-400 pt-2 border-t border-white/5">
        <div class="flex items-center gap-1.5 flex-wrap">
          <span
            v-for="genre in (game.genres || []).slice(0, 1)"
            :key="genre"
            class="font-medium px-2 py-0.5 rounded-md bg-white/5 text-slate-300 border border-white/5 text-[11px]"
          >
            {{ genre }}
          </span>
          <span 
            v-if="(game.genres || []).length > 1" 
            class="text-[10px] text-slate-400 font-semibold"
          >
            +{{ (game.genres || []).length - 1 }}
          </span>
        </div>

        <span v-if="game.file_size && game.file_size > 0" class="font-mono text-[11px] text-slate-400 font-medium">
          {{ formatBytes(game.file_size) }}
        </span>
      </div>
    </div>
  </div>
</template>
