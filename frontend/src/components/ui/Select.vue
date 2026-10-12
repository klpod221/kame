<script setup lang="ts">
import { ChevronDown } from 'lucide-vue-next'

withDefaults(
  defineProps<{
    modelValue: string | number
    disabled?: boolean
  }>(),
  {
    disabled: false,
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', val: string): void
}>()

const handleChange = (e: Event) => {
  emit('update:modelValue', (e.target as HTMLSelectElement).value)
}
</script>

<template>
  <div class="relative inline-flex items-center">
    <select
      :value="modelValue"
      :disabled="disabled"
      @change="handleChange"
      class="appearance-none bg-slate-900 border border-slate-800 text-slate-200 font-semibold text-xs rounded-xl pl-3 pr-8 py-2 focus:outline-none focus:border-indigo-500/80 cursor-pointer disabled:opacity-50"
    >
      <slot />
    </select>
    <ChevronDown class="w-3.5 h-3.5 text-slate-400 absolute right-2.5 pointer-events-none" />
  </div>
</template>
