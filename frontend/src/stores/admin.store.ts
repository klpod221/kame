import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Game, User, GameCandidate } from '../types'
import { adminService, type AdminStats, type RescrapeStatus } from '../services/admin.service'
import { useGameStore } from './game.store'
import { useToastStore } from './toast.store'

export const useAdminStore = defineStore('admin', () => {
  const stats = ref<AdminStats | null>(null)
  const users = ref<User[]>([])
  const isLoading = ref(false)
  const isScraping = ref(false)
  const rescrapeStatus = ref<RescrapeStatus | null>(null)
  let pollingTimer: any = null

  const isRescraping = computed(() => rescrapeStatus.value?.status === 'running')

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

  const fetchRescrapeStatus = async () => {
    try {
      const status = await adminService.getRescrapeStatus()
      rescrapeStatus.value = status
      if (status.status === 'running') {
        startRescrapePolling()
      } else {
        stopRescrapePolling()
      }
      return status
    } catch (err: any) {
      console.error('Failed to fetch rescrape status:', err)
      return null
    }
  }

  const startRescrapePolling = () => {
    if (pollingTimer) return
    pollingTimer = setInterval(async () => {
      try {
        const s = await adminService.getRescrapeStatus()
        rescrapeStatus.value = s
        if (s.status !== 'running') {
          stopRescrapePolling()
          if (s.status === 'completed') {
            toast.success(`Re-scrape complete! Processed ${s.processed}/${s.total} games (${s.success} updated, ${s.failed} skipped/failed).`, 'Steam & Metadata Scraper')
            await gameStore.fetchGames()
            await gameStore.fetchGenres()
            await fetchStats()
          } else if (s.status === 'cancelled') {
            toast.info('Re-scrape worker stopped by admin', 'Worker Cancelled')
          }
        }
      } catch (err) {
        stopRescrapePolling()
      }
    }, 1500)
  }

  const stopRescrapePolling = () => {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
  }

  const triggerRescrapeAll = async () => {
    try {
      const res = await adminService.triggerRescrapeAll()
      toast.info(`Started background re-scraping for ${res.total} games`, 'Worker Running')
      await fetchRescrapeStatus()
      startRescrapePolling()
    } catch (err: any) {
      toast.error(err.message || 'Failed to start re-scrape job', 'Error')
      throw err
    }
  }

  const cancelRescrape = async () => {
    try {
      await adminService.cancelRescrape()
      toast.info('Stopping re-scrape worker...', 'Cancelled')
      await fetchRescrapeStatus()
    } catch (err: any) {
      toast.error(err.message || 'Failed to cancel re-scrape job', 'Error')
    }
  }

  const toggleGameLock = async (game: Game) => {
    const newLocked = !game.is_locked
    try {
      await adminService.toggleGameLock(game.id, newLocked)
      game.is_locked = newLocked
      const found = gameStore.games.find(g => g.id === game.id)
      if (found) {
        found.is_locked = newLocked
      }
      toast.success(
        newLocked ? `Locked "${game.title}". It will be skipped during Re-scrape All.` : `Unlocked "${game.title}".`,
        'Tracking Status'
      )
    } catch (err: any) {
      toast.error(err.message || 'Failed to update lock status', 'Error')
    }
  }

  const rescrapeSingleGame = async (gameId: string) => {
    try {
      const updated = await adminService.rescrapeSingleGame(gameId)
      const idx = gameStore.games.findIndex(g => g.id === gameId)
      if (idx !== -1) {
        gameStore.games[idx] = updated
      }
      toast.success(`Successfully re-scraped metadata for "${updated.title}"`, 'Metadata Updated')
      return updated
    } catch (err: any) {
      toast.error(err.message || 'Failed to re-scrape game metadata', 'Scraper Error')
      throw err
    }
  }

  return {
    stats,
    users,
    isLoading,
    isScraping,
    rescrapeStatus,
    isRescraping,
    fetchStats,
    fetchUsers,
    updateUserRole,
    deleteUser,
    createGame,
    updateGame,
    deleteGame,
    scrapeMetadata,
    fetchRescrapeStatus,
    triggerRescrapeAll,
    cancelRescrape,
    toggleGameLock,
    rescrapeSingleGame,
  }
})
