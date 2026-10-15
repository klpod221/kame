<script setup lang="ts">
import { ref } from 'vue'
import { Shield, ShieldCheck, Trash2, AlertTriangle, UserCheck } from 'lucide-vue-next'
import { useAdminStore } from '../../../stores/admin.store'
import { useAuthStore } from '../../../stores/auth.store'
import { formatDate } from '../../../utils/format'
import Badge from '../../../components/ui/Badge.vue'
import Modal from '../../../components/ui/Modal.vue'
import Button from '../../../components/ui/Button.vue'
import type { User } from '../../../types'

const adminStore = useAdminStore()
const authStore = useAuthStore()

// Modal states
const rolePending = ref<{ user: User; newRole: string } | null>(null)
const deletePending = ref<User | null>(null)

const openRoleModal = (user: User) => {
  const newRole = user.role === 'admin' ? 'member' : 'admin'
  rolePending.value = { user, newRole }
}

const confirmRoleChange = async () => {
  if (rolePending.value) {
    const { user, newRole } = rolePending.value
    rolePending.value = null
    await adminStore.updateUserRole(user.email, newRole)
  }
}

const openDeleteModal = (user: User) => {
  deletePending.value = user
}

const confirmDeleteUser = async () => {
  if (deletePending.value) {
    const user = deletePending.value
    deletePending.value = null
    await adminStore.deleteUser(user.email)
  }
}
</script>

<template>
  <div class="space-y-4 animate-fadeIn">
    
    <div class="flex items-center justify-between pb-3 border-b border-white/5">
      <div>
        <h1 class="text-xl sm:text-2xl font-black text-white tracking-tight flex items-center gap-2.5">
          <span>User Management</span>
          <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-white/5 border border-white/10 text-slate-400">
            {{ adminStore.users.length }}
          </span>
        </h1>
        <p class="text-xs text-slate-400 mt-1 font-normal">Manage team access and administrator permissions</p>
      </div>
    </div>

    <!-- Users Table Container -->
    <div class="rounded-2xl bg-white/2 border border-white/5 overflow-hidden backdrop-blur-xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs">
          <thead class="bg-white/3 border-b border-white/5 text-slate-400 font-semibold uppercase tracking-wider text-[10px]">
            <tr>
              <th class="py-3.5 px-4">User</th>
              <th class="py-3.5 px-4">Role</th>
              <th class="py-3.5 px-4">Last Login</th>
              <th class="py-3.5 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5 text-slate-300">
            <tr 
              v-for="user in adminStore.users" 
              :key="user.email"
              class="hover:bg-white/3 transition-colors"
            >
              <!-- Name & Email -->
              <td class="py-3 px-4">
                <div class="flex items-center gap-3">
                  <img
                    v-if="user.picture"
                    :src="user.picture"
                    :alt="user.name"
                    class="w-9 h-9 rounded-full border border-white/10 object-cover shrink-0"
                  />
                  <div v-else class="w-9 h-9 rounded-full bg-indigo-600/30 text-indigo-400 flex items-center justify-center font-bold text-xs shrink-0">
                    {{ user.name?.charAt(0) || user.email.charAt(0) }}
                  </div>
                  <div>
                    <p class="font-bold text-white text-xs flex items-center gap-1.5">
                      {{ user.name || 'Unnamed' }}
                      <span v-if="user.email === authStore.currentUser?.email" class="text-[10px] text-indigo-400 font-normal">(You)</span>
                    </p>
                    <p class="text-[11px] text-slate-400 font-mono mt-0.5">{{ user.email }}</p>
                  </div>
                </div>
              </td>

              <!-- Role Badge -->
              <td class="py-3 px-4">
                <Badge :variant="user.role === 'admin' ? 'purple' : 'slate'" class="bg-white/5 border-white/10 text-slate-300">
                  <ShieldCheck v-if="user.role === 'admin'" class="w-3 h-3 text-purple-400" />
                  <Shield v-else class="w-3 h-3 text-slate-400" />
                  <span class="capitalize">{{ user.role }}</span>
                </Badge>
              </td>

              <!-- Last Login -->
              <td class="py-3 px-4 text-slate-400 text-[11px]">
                {{ formatDate(String(user.last_login_at)) }}
              </td>

              <!-- Actions -->
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button
                    v-if="user.email !== authStore.currentUser?.email"
                    @click="openRoleModal(user)"
                    class="px-2.5 py-1 rounded-lg bg-white/5 hover:bg-white/10 text-slate-300 text-[11px] font-semibold border border-white/10 transition-colors cursor-pointer"
                  >
                    {{ user.role === 'admin' ? 'Demote to Member' : 'Promote to Admin' }}
                  </button>

                  <button
                    v-if="user.email !== authStore.currentUser?.email"
                    @click="openDeleteModal(user)"
                    class="p-1.5 rounded-lg bg-white/5 text-slate-400 hover:text-rose-400 hover:bg-rose-500/10 border border-white/5 transition-colors cursor-pointer"
                    title="Remove user"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 1. CHANGE ROLE CONFIRMATION MODAL -->
    <Modal
      v-if="rolePending"
      title="Change User Access Level"
      maxWidth="md"
      @close="rolePending = null"
    >
      <div class="space-y-4 text-xs">
        <div class="p-3.5 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 flex items-start gap-3">
          <UserCheck class="w-5 h-5 text-indigo-400 shrink-0 mt-0.5" />
          <div class="space-y-1">
            <p class="font-bold text-white text-sm">Role Update</p>
            <p class="text-slate-300 leading-relaxed">
              Are you sure you want to change the role of <span class="font-bold text-white">{{ rolePending?.user.name || rolePending?.user.email }}</span> to <span class="font-bold text-indigo-300 uppercase">{{ rolePending?.newRole }}</span>?
            </p>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-white/3 border border-white/5 space-y-1">
          <div class="flex justify-between text-slate-400">
            <span>User:</span>
            <span class="font-mono text-slate-200">{{ rolePending?.user.email }}</span>
          </div>
          <div class="flex justify-between text-slate-400">
            <span>New Access Level:</span>
            <span class="font-bold text-indigo-400 uppercase">{{ rolePending?.newRole }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2.5">
          <Button
            variant="ghost"
            size="sm"
            @click="rolePending = null"
            class="text-slate-400 hover:text-white"
          >
            Cancel
          </Button>
          <Button
            variant="primary"
            size="sm"
            @click="confirmRoleChange"
            class="font-semibold shadow-lg shadow-indigo-600/30"
          >
            Confirm Change
          </Button>
        </div>
      </template>
    </Modal>

    <!-- 2. DELETE USER CONFIRMATION MODAL -->
    <Modal
      v-if="deletePending"
      title="Remove User Account"
      maxWidth="md"
      @close="deletePending = null"
    >
      <div class="space-y-4 text-xs">
        <div class="p-3.5 rounded-2xl bg-rose-500/10 border border-rose-500/20 flex items-start gap-3">
          <AlertTriangle class="w-5 h-5 text-rose-400 shrink-0 mt-0.5" />
          <div class="space-y-1">
            <p class="font-bold text-white text-sm">Permanent Revocation</p>
            <p class="text-slate-300 leading-relaxed">
              Are you sure you want to remove <span class="font-bold text-white">{{ deletePending?.name || deletePending?.email }}</span>? This will revoke all their catalog privileges immediately.
            </p>
          </div>
        </div>

        <div class="p-3 rounded-xl bg-white/3 border border-white/5 space-y-1">
          <div class="flex justify-between text-slate-400">
            <span>Email Address:</span>
            <span class="font-mono text-rose-300">{{ deletePending?.email }}</span>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2.5">
          <Button
            variant="ghost"
            size="sm"
            @click="deletePending = null"
            class="text-slate-400 hover:text-white"
          >
            Cancel
          </Button>
          <Button
            variant="danger"
            size="sm"
            @click="confirmDeleteUser"
            class="font-semibold bg-rose-600 hover:bg-rose-500 shadow-lg shadow-rose-600/30"
          >
            Remove Account
          </Button>
        </div>
      </template>
    </Modal>

  </div>
</template>
