<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string
    rows?: number | string
    disabled?: boolean
    error?: string
  }>(),
  {
    rows: 4,
    disabled: false,
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', val: string): void
}>()

const handleInput = (e: Event) => {
  emit('update:modelValue', (e.target as HTMLTextAreaElement).value)
}
</script>

<template>
  <div class="w-full">
    <textarea
      :value="modelValue"
      :rows="rows"
      :placeholder="placeholder"
      :disabled="disabled"
      @input="handleInput"
      class="w-full px-3.5 py-2.5 text-xs sm:text-sm text-slate-100 placeholder-slate-500 bg-slate-900 border rounded-xl transition-all duration-150 focus:outline-none disabled:opacity-50"
      :class="[
        error 
          ? 'border-rose-500/80 focus:border-rose-500 focus:ring-1 focus:ring-rose-500/30' 
          : 'border-slate-850 focus:border-indigo-500/80 focus:ring-1 focus:ring-indigo-500/30'
      ]"
    ></textarea>
    <p v-if="error" class="mt-1 text-[11px] text-rose-400 font-medium">
      {{ error }}
    </p>
  </div>
</template>
