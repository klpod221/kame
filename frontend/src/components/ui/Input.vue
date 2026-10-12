<script setup lang="ts">
import { X } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    modelValue?: string | number | null
    placeholder?: string
    type?: string
    disabled?: boolean
    error?: string
    clearable?: boolean
    size?: 'sm' | 'md'
  }>(),
  {
    modelValue: '',
    type: 'text',
    disabled: false,
    clearable: false,
    size: 'md',
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', val: any): void
  (e: 'clear'): void
}>()

const handleInput = (e: Event) => {
  emit('update:modelValue', (e.target as HTMLInputElement).value)
}

const handleClear = () => {
  emit('update:modelValue', '')
  emit('clear')
}
</script>

<template>
  <div class="relative w-full">
    <div class="relative flex items-center">
      <div v-if="$slots.prefix" class="absolute left-3.5 flex items-center pointer-events-none text-slate-400">
        <slot name="prefix" />
      </div>

      <input
        :type="type"
        :value="modelValue ?? ''"
        :placeholder="placeholder"
        :disabled="disabled"
        @input="handleInput"
        class="w-full text-slate-100 placeholder-slate-500 bg-slate-900 border rounded-xl transition-all duration-150 focus:outline-none disabled:opacity-50"
        :class="[
          $slots.prefix ? 'pl-10' : 'pl-3.5',
          clearable && modelValue ? 'pr-9' : 'pr-3.5',
          size === 'sm' ? 'py-1.5 text-xs' : 'py-2.5 text-xs sm:text-sm',
          error 
            ? 'border-rose-500/80 focus:border-rose-500 focus:ring-1 focus:ring-rose-500/30' 
            : 'border-slate-850 focus:border-indigo-500/80 focus:ring-1 focus:ring-indigo-500/30'
        ]"
      />

      <button
        v-if="clearable && modelValue && !disabled"
        type="button"
        @click="handleClear"
        class="absolute right-3 p-0.5 rounded-full text-slate-400 hover:text-white transition-colors cursor-pointer"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>

    <p v-if="error" class="mt-1 text-[11px] text-rose-400 font-medium">
      {{ error }}
    </p>
  </div>
</template>
