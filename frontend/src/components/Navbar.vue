<script setup lang="ts">
import { Gamepad2, Search, PlusCircle, LogIn, LogOut, ShieldCheck, User as UserIcon } from 'lucide-vue-next'
import { useAuth } from '../composables/useAuth'

defineProps<{
  searchQuery: string
  selectedGenre: string
  genres: string[]
}>()

const emit = defineEmits<{
  (e: 'update:searchQuery', val: string): void
  (e: 'update:selectedGenre', val: string): void
  (e: 'openAddModal'): void
}>()

const { currentUser, isAuthenticated, isAdmin, isLoggingIn, loginError, loginWithGoogle, logout } = useAuth()
</script>

<template>
  <header class="sticky top-0 z-40 w-full backdrop-blur-md bg-slate-950/80 border-b border-slate-800/80 px-4 lg:px-8 py-3">
    <div class="max-w-7xl mx-auto flex items-center justify-between gap-4">
      
      <!-- Brand Logo -->
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-linear-to-br from-indigo-500 to-purple-600 flex items-center justify-center shadow-lg shadow-indigo-500/20">
          <Gamepad2 class="w-6 h-6 text-white" />
        </div>
        <div>
          <span class="text-xl font-bold tracking-tight bg-linear-to-r from-white via-slate-200 to-indigo-300 bg-clip-text text-transparent">
            KAME
          </span>
          <span class="hidden sm:inline-block ml-2 text-xs font-semibold px-2 py-0.5 rounded-full bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
            Catalog
          </span>
        </div>
      </div>

      <!-- Search & Filters -->
      <div class="flex-1 max-w-md mx-2">
        <div class="relative">
          <Search class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input
            type="text"
            :value="searchQuery"
            @input="emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
            placeholder="Search games..."
            class="w-full pl-10 pr-4 py-2 text-sm rounded-xl bg-slate-900/90 border border-slate-800 text-slate-200 placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all"
          />
        </div>
      </div>

      <!-- Action Buttons & Auth -->
      <div class="flex items-center gap-3">
        <!-- Add Game (Admin only) -->
        <button
          v-if="isAdmin"
          @click="emit('openAddModal')"
          class="flex items-center gap-2 px-3.5 py-2 text-sm font-medium rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white shadow-lg shadow-indigo-600/25 transition-all cursor-pointer"
        >
          <PlusCircle class="w-4 h-4" />
          <span class="hidden md:inline">Add Game</span>
        </button>

        <!-- Auth Status -->
        <div v-if="isAuthenticated && currentUser" class="flex items-center gap-3">
          <div class="flex items-center gap-2.5 pl-2 border-l border-slate-800">
            <img
              v-if="currentUser.picture"
              :src="currentUser.picture"
              :alt="currentUser.name"
              class="w-8 h-8 rounded-full border border-indigo-500/30 object-cover"
            />
            <div v-else class="w-8 h-8 rounded-full bg-slate-800 flex items-center justify-center">
              <UserIcon class="w-4 h-4 text-slate-400" />
            </div>

            <div class="hidden lg:block text-left">
              <div class="text-xs font-semibold text-slate-200 flex items-center gap-1">
                {{ currentUser.name }}
                <ShieldCheck v-if="isAdmin" class="w-3.5 h-3.5 text-indigo-400" />
              </div>
              <div class="text-[10px] text-slate-400 capitalize">
                {{ currentUser.role }}
              </div>
            </div>

            <button
              @click="logout"
              title="Logout"
              class="p-2 text-slate-400 hover:text-rose-400 hover:bg-rose-500/10 rounded-lg transition-colors cursor-pointer"
            >
              <LogOut class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Login Button -->
        <div v-else>
          <button
            @click="loginWithGoogle"
            :disabled="isLoggingIn"
            class="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700/60 shadow-sm transition-all cursor-pointer disabled:opacity-50"
          >
            <LogIn class="w-4 h-4 text-indigo-400" />
            <span>{{ isLoggingIn ? 'Signing in...' : 'Sign in with Google' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Error Toast if Login Failed -->
    <div v-if="loginError" class="max-w-7xl mx-auto mt-2 px-4 py-1.5 rounded-lg bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs text-center">
      {{ loginError }}
    </div>
  </header>
</template>
