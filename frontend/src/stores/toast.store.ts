import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface ToastItem {
  id: string
  type: 'success' | 'error' | 'info' | 'warning'
  title?: string
  message: string
  duration?: number
}

export const useToastStore = defineStore('toast', () => {
  const toasts = ref<ToastItem[]>([])

  const removeToast = (id: string) => {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  const show = (toast: Omit<ToastItem, 'id'>) => {
    const id = Math.random().toString(36).substring(2, 9)
    const item: ToastItem = {
      ...toast,
      id,
      duration: toast.duration ?? 4000,
    }
    toasts.value.push(item)

    if (item.duration && item.duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, item.duration)
    }
  }

  const success = (message: string, title?: string) => show({ type: 'success', message, title })
  const error = (message: string, title?: string) => show({ type: 'error', message, title })
  const info = (message: string, title?: string) => show({ type: 'info', message, title })
  const warning = (message: string, title?: string) => show({ type: 'warning', message, title })

  return {
    toasts,
    show,
    success,
    error,
    info,
    warning,
    removeToast,
  }
})
