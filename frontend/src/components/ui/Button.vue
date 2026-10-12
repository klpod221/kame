<script setup lang="ts">
import { computed } from 'vue'
import { Loader2 } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
    size?: 'sm' | 'md' | 'lg'
    loading?: boolean
    disabled?: boolean
    type?: 'button' | 'submit' | 'reset'
  }>(),
  {
    variant: 'secondary',
    size: 'md',
    loading: false,
    disabled: false,
    type: 'button',
  }
)

const variantClasses = computed(() => {
  switch (props.variant) {
    case 'primary':
      return 'bg-linear-to-r from-indigo-600 via-indigo-500 to-purple-600 hover:from-indigo-500 hover:to-purple-500 text-white shadow-lg shadow-indigo-600/30 font-bold active:scale-98'
    case 'danger':
      return 'bg-rose-500/10 hover:bg-rose-500/20 text-rose-300 border border-rose-500/30 font-semibold active:scale-98'
    case 'ghost':
      return 'bg-transparent hover:bg-slate-800/60 text-slate-400 hover:text-white font-medium'
    case 'secondary':
    default:
      return 'bg-slate-900 hover:bg-slate-800 text-slate-200 hover:text-white border border-slate-800 hover:border-slate-700 font-semibold active:scale-98 shadow-xs'
  }
})

const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'px-2.5 py-1 text-xs rounded-lg gap-1.5'
    case 'lg':
      return 'px-6 py-3.5 text-sm sm:text-base rounded-2xl gap-2.5 font-black'
    case 'md':
    default:
      return 'px-4 py-2.5 text-xs sm:text-sm rounded-xl gap-2'
  }
})
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="inline-flex items-center justify-center transition-all duration-150 cursor-pointer disabled:opacity-50 disabled:pointer-events-none select-none"
    :class="[variantClasses, sizeClasses]"
  >
    <Loader2 v-if="loading" class="w-4 h-4 animate-spin shrink-0" />
    <slot />
  </button>
</template>
