import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import CatalogView from '../views/catalog/CatalogView.vue'
import GameDetailView from '../views/game-detail/GameDetailView.vue'
import AdminView from '../views/admin/AdminView.vue'
import LoginView from '../views/auth/LoginView.vue'
import { useAuthStore } from '../stores/auth.store'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { guestOnly: true }
  },
  {
    path: '/',
    name: 'library',
    component: CatalogView,
    meta: { requiresAuth: true }
  },
  {
    path: '/games/:id',
    name: 'game-detail',
    component: GameDetailView,
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'admin',
    component: AdminView,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  // Ensure server config is loaded
  if (!authStore.appConfig) {
    await authStore.fetchConfig()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login' })
    return
  }

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    next({ name: 'library' })
    return
  }

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'library' })
    return
  }

  next()
})

export default router
