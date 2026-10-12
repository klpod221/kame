export interface Game {
  id: string
  title: string
  description: string
  genres: string[]
  cover_url: string
  background_url: string
  drive_file_id: string
  file_size: number
  release_year: number
  rating: number
  created_at: string
  updated_at: string
}

export interface User {
  email: string
  name: string
  picture: string
  role: 'admin' | 'user'
}

export interface AppConfig {
  google_client_id: string
  is_authenticated: boolean
  user: User | null
}

export interface GameCandidate {
  source: 'rawg' | 'igdb'
  source_id: string
  title: string
  description: string
  genres: string[]
  cover_url: string
  background_url: string
  release_year: number
  rating: number
}
