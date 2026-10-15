<script setup lang="ts">
import { useAdminStore } from '../../../stores/admin.store'
import { formatBytes } from '../../../utils/format'

const adminStore = useAdminStore()
</script>

<template>
  <div class="space-y-4 animate-fadeIn">
    
    <div class="flex items-center justify-between pb-3 border-b border-white/5">
      <div>
        <h1 class="text-xl sm:text-2xl font-black text-white tracking-tight">System Info</h1>
        <p class="text-xs text-slate-400 mt-1 font-normal">Deployment diagnostics, scrapers status, and storage telemetry</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      
      <!-- Deployment Container Specs -->
      <div class="p-5 rounded-2xl bg-white/2 border border-white/5 space-y-3 backdrop-blur-xl">
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-400">Environment & Container</h3>
        <div class="space-y-2 text-xs">
          <div class="flex justify-between py-1.5 border-b border-white/5">
            <span class="text-slate-400">Application</span>
            <span class="font-bold text-white">Kame Game Vault</span>
          </div>
          <div class="flex justify-between py-1.5 border-b border-white/5">
            <span class="text-slate-400">Architecture</span>
            <span class="font-mono text-indigo-400">Go 1.24 + SQLite WAL + Vue 3 SPA</span>
          </div>
          <div class="flex justify-between py-1.5 border-b border-white/5">
            <span class="text-slate-400">Docker Port</span>
            <span class="font-mono text-emerald-400 font-bold">:8324</span>
          </div>
          <div class="flex justify-between py-1.5 border-b border-white/5">
            <span class="text-slate-400">Database Path</span>
            <span class="font-mono text-slate-300 text-[11px] truncate max-w-50">{{ adminStore.stats?.db_path || 'data/kame.db' }}</span>
          </div>
          <div class="flex justify-between py-1.5">
            <span class="text-slate-400">Database Size</span>
            <span class="font-mono text-white font-bold">{{ formatBytes(adminStore.stats?.db_size_bytes ?? 0) }}</span>
          </div>
        </div>
      </div>

      <!-- Scrapers & Integrations -->
      <div class="p-5 rounded-2xl bg-white/2 border border-white/5 space-y-3 backdrop-blur-xl">
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-400">Integrations & Upstream</h3>
        <div class="space-y-2 text-xs">
          <div class="flex justify-between py-1.5 border-b border-white/5">
            <span class="text-slate-400">RAWG Video Games API</span>
            <span :class="(adminStore.stats?.has_rawg || (adminStore.stats as any)?.has_rawg_key) ? 'text-emerald-400 font-bold' : 'text-slate-500'">
              {{ (adminStore.stats?.has_rawg || (adminStore.stats as any)?.has_rawg_key) ? 'Connected' : 'Not configured' }}
            </span>
          </div>
          <div class="flex justify-between py-1.5 border-b border-white/5">
            <span class="text-slate-400">IGDB Fallback Scraper</span>
            <span :class="adminStore.stats?.has_igdb ? 'text-emerald-400 font-bold' : 'text-slate-500'">
              {{ adminStore.stats?.has_igdb ? 'Active' : 'Not configured' }}
            </span>
          </div>
          <div class="flex justify-between py-1.5 border-b border-white/5">
            <span class="text-slate-400">Access Mode</span>
            <span :class="adminStore.stats?.whitelist_enabled ? 'text-amber-400 font-bold' : 'text-emerald-400 font-bold'">
              {{ adminStore.stats?.whitelist_enabled ? 'Whitelist Restricted' : 'Open (Any Google Account)' }}
            </span>
          </div>
          <div class="flex justify-between py-1.5">
            <span class="text-slate-400">Author & System Architect</span>
            <span class="font-bold text-indigo-400">klpod221</span>
          </div>
        </div>
      </div>

    </div>

  </div>
</template>
