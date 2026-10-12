<script setup lang="ts">
withDefaults(
  defineProps<{
    padding?: 'none' | 'sm' | 'md' | 'lg'
  }>(),
  {
    padding: 'md',
  }
)
</script>

<template>
  <div
    class="rounded-2xl bg-slate-900/80 border border-slate-850 shadow-xl backdrop-blur-md transition-all"
    :class="{
      'p-0': padding === 'none',
      'p-3 sm:p-4': padding === 'sm',
      'p-5 sm:p-6': padding === 'md',
      'p-6 sm:p-8': padding === 'lg',
    }"
  >
    <slot />
  </div>
</template>
