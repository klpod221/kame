import { api } from './api'
import type { Game } from '../types'

export interface GameListResponse {
  games: Game[]
  total: number
  page: number
  limit: number
}

export interface GenreListResponse {
  genres: string[]
}

export const gameService = {
  async listGames(params?: { q?: string; genre?: string; page?: number; limit?: number }, signal?: AbortSignal): Promise<GameListResponse> {
    return api.get<GameListResponse>('/api/games', {
      params,
      signal,
    })
  },

  async getGame(id: string): Promise<Game> {
    return api.get<Game>(`/api/games/${id}`)
  },

  async listGenres(): Promise<string[]> {
    const res = await api.get<GenreListResponse>('/api/genres')
    return res.genres || []
  },

  getDownloadUrl(gameId: string, googleToken?: string | null): string {
    const tokenParam = googleToken ? `?token=${encodeURIComponent(googleToken)}` : ''
    return `/api/games/${gameId}/download${tokenParam}`
  },

  getDriveWebUrl(driveFileId: string): string {
    return `https://drive.google.com/uc?export=download&id=${driveFileId.trim()}`
  },

  getCurlCommand(driveFileId: string, title: string, googleToken?: string | null): string {
    const token = googleToken || 'YOUR_GOOGLE_ACCESS_TOKEN'
    const safeName = title.replace(/[^a-zA-Z0-9_\-.]/g, '_')
    return `curl -L -C - -H "Authorization: Bearer ${token}" "https://www.googleapis.com/drive/v3/files/${driveFileId.trim()}?alt=media&acknowledgeAbuse=true" -o "${safeName}.zip"`
  },
}
