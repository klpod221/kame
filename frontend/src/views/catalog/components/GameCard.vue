<script setup lang="ts">
import { ref } from 'vue'
import { Star, Calendar, Gamepad2, ArrowRight, Download } from 'lucide-vue-next'
import type { Game } from '../../../types'
import { formatBytes } from '../../../utils/format'

defineProps<{
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
    class="group relative aspect-2/3 w-full rounded-md overflow-hidden bg-[#101722] border border-[#2a475e]/40 transition-all duration-300 ease-out cursor-pointer select-none shadow-[0_4px_16px_rgba(0,0,0,0.6)] hover:shadow-[0_16px_36px_rgba(0,0,0,0.95),0_0_20px_rgba(102,192,244,0.25)] hover:border-[#66c0f4]/80 hover:scale-[1.03]"
  >
    <!-- Skeleton shimmer placeholder while image loads -->
    <div 
      v-if="!isImageLoaded && !hasImageError && game.cover_url"
      class="absolute inset-0 bg-slate-800/80 animate-pulse"
    ></div>

    <!-- Steam Full-bleed Cover Image with subtle zoom on hover -->
    <img
      v-if="game.cover_url && !hasImageError"
      :src="game.cover_url"
      :alt="game.title"
      loading="lazy"
      @load="onImageLoad"
      @error="onImageError"
      class="w-full h-full object-cover transition-transform duration-500 ease-out group-hover:scale-108"
      :class="isImageLoaded ? 'opacity-100' : 'opacity-0'"
    />

    <!-- Fallback when no cover image or on load error -->
    <div 
      v-else 
      class="w-full h-full flex flex-col items-center justify-center text-slate-500 bg-[#141c28] p-4 text-center gap-2 border border-white/5"
    >
      <Gamepad2 class="w-10 h-10 text-slate-600 group-hover:text-[#66c0f4] transition-colors" />
      <span class="text-xs font-bold text-slate-200 line-clamp-3 uppercase tracking-wider">
        {{ game.title }}
      </span>
    </div>

    <!-- Always visible subtle bottom vignette for contrast -->
    <div class="absolute inset-x-0 bottom-0 h-16 bg-linear-to-t from-black/60 to-transparent pointer-events-none group-hover:opacity-0 transition-opacity duration-200"></div>

    <!-- IN-CARD HOVER OVERLAY (Displays directly ON TOP of the card) -->
    <div 
      class="absolute inset-0 bg-linear-to-t from-[#080d16]/98 via-[#080d16]/85 to-[#080d16]/40 backdrop-blur-[5px] p-3.5 sm:p-4 flex flex-col justify-between opacity-0 group-hover:opacity-100 transition-all duration-300 ease-out z-20 pointer-events-none group-hover:pointer-events-auto"
    >
      <!-- Top badges: Release Year & Rating -->
      <div class="flex items-center justify-between gap-2 transform -translate-y-2 group-hover:translate-y-0 transition-transform duration-300">
        <span 
          v-if="game.release_year" 
          class="flex items-center gap-1 text-[10px] font-bold px-2 py-0.5 rounded bg-black/60 text-slate-300 border border-white/10 backdrop-blur-md shadow-xs"
        >
          <Calendar class="w-3 h-3 text-[#66c0f4]" />
          {{ game.release_year }}
        </span>
        <span v-else></span>

        <span 
          v-if="game.rating > 0" 
          class="flex items-center gap-1 text-[10px] font-extrabold px-2 py-0.5 rounded bg-amber-500/20 text-amber-300 border border-amber-500/40 backdrop-blur-md shadow-xs"
        >
          <Star class="w-3 h-3 fill-amber-400 text-amber-400" />
          {{ game.rating.toFixed(1) }}
        </span>
      </div>

      <!-- Bottom content: Title, Genres, Size & Action -->
      <div class="flex flex-col gap-2 transform translate-y-3 group-hover:translate-y-0 transition-transform duration-300">
        <!-- Game Title -->
        <h4 class="text-xs sm:text-sm font-black text-white line-clamp-2 tracking-tight leading-snug drop-shadow-md">
          {{ game.title }}
        </h4>

        <!-- Genres chips & File Size -->
        <div class="flex items-center justify-between gap-1.5 flex-wrap text-[10px]">
          <div class="flex items-center gap-1 flex-wrap">
            <span 
              v-for="g in (game.genres || []).slice(0, 2)" 
              :key="g"
              class="px-1.5 py-0.5 rounded bg-white/10 text-slate-200 font-medium"
            >
              {{ g }}
            </span>
            <span v-if="(game.genres || []).length > 2" class="text-slate-400 font-bold">
              +{{ (game.genres || []).length - 2 }}
            </span>
          </div>

          <span v-if="game.file_size > 0" class="font-mono text-slate-400 ml-auto shrink-0 font-medium">
            {{ formatBytes(game.file_size) }}
          </span>
        </div>

        <!-- Short Description -->
        <p v-if="game.description" class="text-[10px] text-slate-300 line-clamp-2 leading-relaxed opacity-90">
          {{ game.description }}
        </p>

        <!-- Steam Action Bar Button -->
        <div class="pt-2 border-t border-white/10 flex items-center justify-between text-[11px] font-bold text-[#66c0f4]">
          <span class="flex items-center gap-1.5">
            <Download class="w-3.5 h-3.5 text-[#66c0f4]" />
            <span>View & Download</span>
          </span>
          <ArrowRight class="w-3.5 h-3.5 group-hover:translate-x-1 transition-transform" />
        </div>
      </div>
    </div>

  </div>
</template>
