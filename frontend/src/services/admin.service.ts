import { api } from './api'
import type { Game, User, GameCandidate } from '../types'

export interface AdminStats {
  total_games: number
  total_users: number
  total_genres: number
  db_size_bytes: number
  db_path: string
  has_rawg?: boolean
  has_igdb?: boolean
  whitelist_enabled?: boolean
}

export const adminService = {
  async getStats(): Promise<AdminStats> {
    return api.get<AdminStats>('/api/admin/stats')
  },

  async listUsers(): Promise<User[]> {
    const res = await api.get<{ users: User[] }>('/api/admin/users')
    return res.users || []
  },

  async updateUserRole(email: string, role: string): Promise<void> {
    await api.put(`/api/admin/users/${encodeURIComponent(email)}/role`, { role })
  },

  async deleteUser(email: string): Promise<void> {
    await api.delete(`/api/admin/users/${encodeURIComponent(email)}`)
  },

  async createGame(gameData: Partial<Game>): Promise<Game> {
    return api.post<Game>('/api/admin/games', gameData)
  },

  async updateGame(id: string, gameData: Partial<Game>): Promise<Game> {
    return api.put<Game>(`/api/admin/games/${id}`, gameData)
  },

  async deleteGame(id: string): Promise<void> {
    await api.delete(`/api/admin/games/${id}`)
  },

  async scrapeMetadata(query: string): Promise<GameCandidate[]> {
    const res = await api.post<{ candidates: GameCandidate[] }>('/api/admin/scrape', { query })
    return res.candidates || []
  },
}
