<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    variant?: 'indigo' | 'purple' | 'amber' | 'emerald' | 'rose' | 'slate'
    size?: 'sm' | 'md'
  }>(),
  {
    variant: 'slate',
    size: 'sm',
  }
)

const variantClasses = computed(() => {
  switch (props.variant) {
    case 'indigo':
      return 'bg-indigo-500/10 text-indigo-400 border-indigo-500/25'
    case 'purple':
      return 'bg-purple-500/10 text-purple-400 border-purple-500/25'
    case 'amber':
      return 'bg-amber-500/15 text-amber-300 border-amber-500/35'
    case 'emerald':
      return 'bg-emerald-500/15 text-emerald-300 border-emerald-500/35'
    case 'rose':
      return 'bg-rose-500/15 text-rose-300 border-rose-500/35'
    case 'slate':
    default:
      return 'bg-slate-900 text-slate-300 border-slate-800'
  }
})
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 border font-semibold select-none rounded-lg"
    :class="[
      variantClasses,
      size === 'sm' ? 'px-2 py-0.5 text-[11px]' : 'px-3 py-1 text-xs rounded-xl font-bold'
    ]"
  >
    <slot />
  </span>
</template>
