<script setup lang="ts">
import { AlertCircle, Loader2 } from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth.store'

const authStore = useAuthStore()
</script>

<template>
  <div class="min-h-screen bg-[#070a0f] flex flex-col justify-between items-center relative overflow-hidden select-none">
    
    <!-- Ambient Cinema Glow effects -->
    <div class="absolute -top-32 left-1/2 -translate-x-1/2 w-125 h-125 bg-indigo-600/15 rounded-full blur-[140px] pointer-events-none"></div>
    <div class="absolute -bottom-32 left-1/2 -translate-x-1/2 w-125 h-125 bg-purple-600/10 rounded-full blur-[140px] pointer-events-none"></div>

    <!-- Header Spacer -->
    <div></div>

    <!-- Centered Minimal Auth Box -->
    <main class="w-full max-w-sm px-6 py-8 z-10">
      <div class="relative rounded-3xl bg-white/3 border border-white/10 p-8 shadow-2xl shadow-black/80 backdrop-blur-2xl space-y-6 text-center animate-fadeIn">
        
        <!-- App Logo & Title with optimized icon -->
        <div class="space-y-3">
          <div class="w-18 h-18 mx-auto flex items-center justify-center">
            <img 
              src="/kame-128.png" 
              alt="Kame" 
              class="w-full h-full object-contain filter drop-shadow-[0_8px_20px_rgba(99,102,241,0.35)]" 
            />
          </div>

          <div>
            <h1 class="text-2xl font-black tracking-tight text-white flex items-center justify-center gap-2">
              <span>Kame</span>
              <span class="text-[9px] font-black px-1.5 py-0.5 rounded bg-indigo-500/20 text-indigo-300 border border-indigo-500/30 uppercase tracking-widest">
                vault
              </span>
            </h1>
            <p class="text-xs text-slate-400 mt-1 font-normal">
              Sign in to access your private game vault
            </p>
          </div>
        </div>

        <!-- Error Feedback -->
        <div 
          v-if="authStore.loginError" 
          class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/25 text-rose-300 text-xs flex items-start gap-2.5 text-left animate-shake"
        >
          <AlertCircle class="w-4 h-4 shrink-0 text-rose-400 mt-0.5" />
          <div class="flex-1 leading-tight">
            {{ authStore.loginError }}
          </div>
        </div>

        <!-- Google Sign In Button -->
        <div>
          <button
            @click="authStore.loginWithGoogle"
            :disabled="authStore.isLoggingIn"
            class="group w-full flex items-center justify-center gap-3 px-5 py-3.5 rounded-2xl bg-white hover:bg-slate-100 text-slate-900 font-bold text-xs shadow-lg shadow-black/40 hover:shadow-indigo-500/20 transition-all duration-200 transform active:scale-98 cursor-pointer disabled:opacity-60"
          >
            <div v-if="authStore.isLoggingIn" class="flex items-center gap-2 text-slate-800">
              <Loader2 class="w-4 h-4 animate-spin text-indigo-600" />
              <span>Connecting...</span>
            </div>
            <template v-else>
              <svg class="w-4 h-4 shrink-0" viewBox="0 0 24 24">
                <path
                  fill="#4285F4"
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                />
                <path
                  fill="#34A853"
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                />
                <path
                  fill="#FBBC05"
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"
                />
                <path
                  fill="#EA4335"
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"
                />
              </svg>
              <span>Sign in with Google</span>
            </template>
          </button>
        </div>

        <p class="text-[10px] text-slate-400">
          Private access for authorized team members.
        </p>

      </div>
    </main>

    <!-- Minimal Footer -->
    <footer class="py-6 text-center text-xs text-slate-400 z-10 font-medium">
      Kame • Game Vault by klpod221
    </footer>

  </div>
</template>
