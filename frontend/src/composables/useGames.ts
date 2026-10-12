import { ref } from 'vue'
import type { Game, GameCandidate } from '../types'
import { useAuth } from './useAuth'

const games = ref<Game[]>([])
const total = ref(0)
const genres = ref<string[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)

export function useGames() {
  const { token, googleAccessToken } = useAuth()

  const loadGames = async (search = '', genre = '', page = 1, limit = 24) => {
    isLoading.value = true
    error.value = null
    try {
      const params = new URLSearchParams()
      if (search) params.append('q', search)
      if (genre) params.append('genre', genre)
      params.append('page', page.toString())
      params.append('limit', limit.toString())

      const res = await fetch(`/api/games?${params.toString()}`)
      if (!res.ok) throw new Error('Failed to fetch games')
      const data = await res.json()
      games.value = data.games || []
      total.value = data.total || 0
    } catch (err: any) {
      error.value = err.message
    } finally {
      isLoading.value = false
    }
  }

  const loadGenres = async () => {
    try {
      const res = await fetch('/api/genres')
      if (res.ok) {
        const data = await res.json()
        genres.value = data.genres || []
      }
    } catch (err) {
      console.error('Failed to load genres:', err)
    }
  }

  const createGame = async (gameData: Partial<Game>) => {
    if (!token.value) throw new Error('Authentication required')
    const res = await fetch('/api/admin/games', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify(gameData),
    })

    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Failed to create game')
    await loadGames()
    await loadGenres()
    return data
  }

  const deleteGame = async (id: string) => {
    if (!token.value) throw new Error('Authentication required')
    const res = await fetch(`/api/admin/games/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token.value}`,
      },
    })

    if (!res.ok) {
      const data = await res.json()
      throw new Error(data.error || 'Failed to delete game')
    }
    await loadGames()
  }

  const scrapeMetadata = async (query: string): Promise<GameCandidate[]> => {
    if (!token.value) throw new Error('Authentication required')
    const res = await fetch('/api/admin/scrape', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify({ query }),
    })

    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'Scraping failed')
    return data.candidates || []
  }

  const getDownloadUrl = (driveFileId: string): string => {
    const cleanId = driveFileId.trim()
    const gToken = googleAccessToken.value || ''
    return `https://www.googleapis.com/drive/v3/files/${cleanId}?alt=media&acknowledgeAbuse=true&access_token=${encodeURIComponent(gToken)}`
  }

  return {
    games,
    total,
    genres,
    isLoading,
    error,
    loadGames,
    loadGenres,
    createGame,
    deleteGame,
    scrapeMetadata,
    getDownloadUrl,
  }
}
