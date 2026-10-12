import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Game } from '../types'
import { gameService } from '../services/game.service'
import { useAuthStore } from './auth.store'
import { useToastStore } from './toast.store'

export const useGameStore = defineStore('game', () => {
  const games = ref<Game[]>([])
  const total = ref(0)
  const genres = ref<string[]>([])
  const isLoading = ref(false)
  const searchQuery = ref('')
  const selectedGenre = ref('')
  const sortBy = ref<'rating' | 'newest' | 'title'>('rating')
  const activeGame = ref<Game | null>(null)
  const isLoadingDetail = ref(false)

  let activeAbortController: AbortController | null = null

  const sortedGames = computed(() => {
    const list = [...games.value]
    if (sortBy.value === 'rating') {
      return list.sort((a, b) => b.rating - a.rating)
    }
    if (sortBy.value === 'newest') {
      return list.sort((a, b) => b.release_year - a.release_year)
    }
    if (sortBy.value === 'title') {
      return list.sort((a, b) => a.title.localeCompare(b.title))
    }
    return list
  })

  const fetchGames = async () => {
    // Abort previous in-flight search request to prevent race conditions
    if (activeAbortController) {
      activeAbortController.abort()
    }
    activeAbortController = new AbortController()

    isLoading.value = true
    try {
      const data = await gameService.listGames({
        q: searchQuery.value.trim(),
        genre: selectedGenre.value.trim(),
        page: 1,
        limit: 100,
      }, activeAbortController.signal)

      games.value = data.games || []
      total.value = data.total || 0
    } catch (err: any) {
      if (err.name !== 'AbortError') {
        console.error('Failed to fetch games:', err)
      }
    } finally {
      isLoading.value = false
    }
  }

  const fetchGenres = async () => {
    try {
      genres.value = await gameService.listGenres()
    } catch (err) {
      console.error('Failed to fetch genres:', err)
    }
  }

  const fetchGameById = async (id: string): Promise<Game | null> => {
    // Check cached games list
    const cached = games.value.find(g => g.id === id)
    if (cached) {
      activeGame.value = cached
      return cached
    }

    isLoadingDetail.value = true
    try {
      const game = await gameService.getGame(id)
      activeGame.value = game
      return game
    } catch (err: any) {
      const toast = useToastStore()
      toast.error(err.message || 'Failed to load game details', 'Error')
      return null
    } finally {
      isLoadingDetail.value = false
    }
  }

  const downloadGame = (game: Game) => {
    const authStore = useAuthStore()
    const toast = useToastStore()

    if (!authStore.isAuthenticated) {
      authStore.loginWithGoogle()
      return
    }

    const downloadUrl = gameService.getDownloadUrl(game.id, authStore.googleAccessToken)
    const a = document.createElement('a')
    a.href = downloadUrl
    a.download = `${game.title}.zip`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)

    toast.info(`Preparing direct stream for "${game.title}"...`, 'Download Started')
  }

  const copyCurl = async (game: Game) => {
    const authStore = useAuthStore()
    const toast = useToastStore()

    if (!authStore.isAuthenticated) {
      authStore.loginWithGoogle()
      return
    }

    const cmd = gameService.getCurlCommand(game.drive_file_id, game.title, authStore.googleAccessToken)
    await navigator.clipboard.writeText(cmd)
    toast.success('High-speed cURL command copied to clipboard!', 'cURL Copied')
  }

  const openInDrive = (game: Game) => {
    const url = gameService.getDriveWebUrl(game.drive_file_id)
    window.open(url, '_blank')
  }

  return {
    games,
    total,
    genres,
    isLoading,
    searchQuery,
    selectedGenre,
    sortBy,
    sortedGames,
    activeGame,
    isLoadingDetail,
    fetchGames,
    fetchGenres,
    fetchGameById,
    downloadGame,
    copyCurl,
    openInDrive,
  }
})
