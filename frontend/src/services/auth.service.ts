import { api } from './api'
import type { User, AppConfig } from '../types'

export interface LoginResponse {
  token: string
  user: User
}

export const authService = {
  async getConfig(): Promise<AppConfig> {
    return api.get<AppConfig>('/api/config')
  },

  async loginWithGoogleToken(accessToken: string): Promise<LoginResponse> {
    return api.post<LoginResponse>('/api/auth/google', { access_token: accessToken })
  },

  async getMe(): Promise<User> {
    return api.get<User>('/api/auth/me')
  },
}
