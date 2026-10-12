import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Game, User, GameCandidate } from '../types'
import { adminService, type AdminStats } from '../services/admin.service'
import { useGameStore } from './game.store'
import { useToastStore } from './toast.store'

export const useAdminStore = defineStore('admin', () => {
  const stats = ref<AdminStats | null>(null)
  const users = ref<User[]>([])
  const isLoading = ref(false)
  const isScraping = ref(false)

  const toast = useToastStore()
  const gameStore = useGameStore()

  const fetchStats = async () => {
    try {
      stats.value = await adminService.getStats()
    } catch (err: any) {
      console.error('Failed to fetch admin stats:', err)
    }
  }

  const fetchUsers = async () => {
    isLoading.value = true
    try {
      users.value = await adminService.listUsers()
    } catch (err: any) {
      toast.error(err.message || 'Failed to fetch users', 'Admin Error')
    } finally {
      isLoading.value = false
    }
  }

  const updateUserRole = async (email: string, role: string) => {
    try {
      await adminService.updateUserRole(email, role)
      toast.success(`Role for ${email} updated to ${role.toUpperCase()}`, 'Success')
      await fetchUsers()
      await fetchStats()
    } catch (err: any) {
      toast.error(err.message || 'Failed to update user role', 'Error')
      throw err
    }
  }

  const deleteUser = async (email: string) => {
    try {
      await adminService.deleteUser(email)
      toast.success(`User ${email} removed`, 'Success')
      await fetchUsers()
      await fetchStats()
    } catch (err: any) {
      toast.error(err.message || 'Failed to remove user', 'Error')
      throw err
    }
  }

  const createGame = async (gameData: Partial<Game>): Promise<Game> => {
    try {
      const created = await adminService.createGame(gameData)
      toast.success(`Game "${created.title}" added to vault!`, 'Game Created')
      await gameStore.fetchGames()
      await gameStore.fetchGenres()
      await fetchStats()
      return created
    } catch (err: any) {
      toast.error(err.message || 'Failed to create game', 'Error')
      throw err
    }
  }

  const updateGame = async (id: string, gameData: Partial<Game>): Promise<Game> => {
    try {
      const updated = await adminService.updateGame(id, gameData)
      toast.success(`"${updated.title}" updated successfully`, 'Changes Saved')
      await gameStore.fetchGames()
      await gameStore.fetchGenres()
      await fetchStats()
      return updated
    } catch (err: any) {
      toast.error(err.message || 'Failed to update game', 'Error')
      throw err
    }
  }

  const deleteGame = async (id: string, title?: string) => {
    try {
      await adminService.deleteGame(id)
      toast.success(`"${title || 'Game'}" deleted from vault`, 'Game Removed')
      await gameStore.fetchGames()
      await gameStore.fetchGenres()
      await fetchStats()
    } catch (err: any) {
      toast.error(err.message || 'Failed to delete game', 'Error')
      throw err
    }
  }

  const scrapeMetadata = async (query: string): Promise<GameCandidate[]> => {
    isScraping.value = true
    try {
      const candidates = await adminService.scrapeMetadata(query)
      if (candidates.length === 0) {
        toast.info('No matches found from RAWG or IGDB', 'Scraper')
      }
      return candidates
    } catch (err: any) {
      toast.error(err.message || 'Scraper failed to query upstream APIs', 'Scraper Error')
      return []
    } finally {
      isScraping.value = false
    }
  }

  return {
    stats,
    users,
    isLoading,
    isScraping,
    fetchStats,
    fetchUsers,
    updateUserRole,
    deleteUser,
    createGame,
    updateGame,
    deleteGame,
    scrapeMetadata,
  }
})
