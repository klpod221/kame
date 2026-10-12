<script setup lang="ts">
import { useToastStore } from '../../stores/toast.store'
import { CheckCircle2, AlertCircle, Info, AlertTriangle, X } from 'lucide-vue-next'

const toastStore = useToastStore()
</script>

<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none">
      <TransitionGroup
        enter-active-class="transform ease-out duration-300 transition"
        enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
        enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
        leave-active-class="transition ease-in duration-200"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-for="toast in toastStore.toasts"
          :key="toast.id"
          class="pointer-events-auto flex items-start gap-3 p-4 rounded-2xl border shadow-xl backdrop-blur-xl text-xs sm:text-sm font-medium"
          :class="{
            'bg-emerald-950/90 border-emerald-500/40 text-emerald-200': toast.type === 'success',
            'bg-rose-950/90 border-rose-500/40 text-rose-200': toast.type === 'error',
            'bg-amber-950/90 border-amber-500/40 text-amber-200': toast.type === 'warning',
            'bg-slate-900/90 border-slate-700/80 text-slate-200': toast.type === 'info',
          }"
        >
          <CheckCircle2 v-if="toast.type === 'success'" class="w-5 h-5 text-emerald-400 shrink-0 mt-0.5" />
          <AlertCircle v-else-if="toast.type === 'error'" class="w-5 h-5 text-rose-400 shrink-0 mt-0.5" />
          <AlertTriangle v-else-if="toast.type === 'warning'" class="w-5 h-5 text-amber-400 shrink-0 mt-0.5" />
          <Info v-else class="w-5 h-5 text-indigo-400 shrink-0 mt-0.5" />

          <div class="flex-1 min-w-0">
            <h4 v-if="toast.title" class="font-bold text-white text-xs mb-0.5">
              {{ toast.title }}
            </h4>
            <p class="leading-relaxed wrap-break-word text-xs">
              {{ toast.message }}
            </p>
          </div>

          <button
            @click="toastStore.removeToast(toast.id)"
            class="text-slate-400 hover:text-white transition-colors cursor-pointer shrink-0 -mr-1 -mt-1 p-1"
          >
            <X class="w-4 h-4" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
