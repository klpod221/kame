<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { X } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    modelValue?: boolean
    title?: string
    maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '4xl'
  }>(),
  {
    modelValue: true,
    maxWidth: 'lg',
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', val: boolean): void
  (e: 'close'): void
}>()

const close = () => {
  emit('update:modelValue', false)
  emit('close')
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    close()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <Teleport to="body">
    <div class="fixed inset-0 z-50 overflow-y-auto">
      <!-- Backdrop Overlay -->
      <div
        class="fixed inset-0 bg-black/75 backdrop-blur-sm transition-opacity animate-fadeIn"
        @click="close"
      ></div>

      <!-- Modal Dialog Container -->
      <div class="flex min-h-full items-center justify-center p-4 sm:p-6 text-center">
        <div
          class="w-full transform rounded-3xl bg-slate-900 border border-slate-800 text-left align-middle shadow-2xl transition-all animate-fadeIn relative z-10 flex flex-col max-h-[90vh]"
          :class="{
            'max-w-sm': maxWidth === 'sm',
            'max-w-md': maxWidth === 'md',
            'max-w-lg': maxWidth === 'lg',
            'max-w-xl': maxWidth === 'xl',
            'max-w-2xl': maxWidth === '2xl',
            'max-w-4xl': maxWidth === '4xl',
          }"
        >
          <!-- Header Bar -->
          <div v-if="title || $slots.header" class="flex items-center justify-between px-6 py-5 border-b border-slate-800/80 shrink-0">
            <slot name="header">
              <h3 class="text-base sm:text-lg font-black text-white tracking-tight">
                {{ title }}
              </h3>
            </slot>
            <button
              @click="close"
              class="p-1.5 rounded-xl text-slate-400 hover:text-white hover:bg-slate-800 transition-colors cursor-pointer"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Body Content -->
          <div class="px-6 py-5 overflow-y-auto flex-1">
            <slot />
          </div>

          <!-- Footer Bar -->
          <div v-if="$slots.footer" class="px-6 py-4 border-t border-slate-800/80 bg-slate-950/40 rounded-b-3xl shrink-0">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
