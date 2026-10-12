import { ref, computed } from 'vue'
import type { User, AppConfig } from '../types'

const token = ref<string | null>(localStorage.getItem('kame_jwt_token'))
const googleAccessToken = ref<string | null>(localStorage.getItem('kame_google_access_token'))
const currentUser = ref<User | null>(
  localStorage.getItem('kame_user') ? JSON.parse(localStorage.getItem('kame_user')!) : null
)
const appConfig = ref<AppConfig | null>(null)
const isLoggingIn = ref(false)
const loginError = ref<string | null>(null)

let tokenClient: any = null

export function useAuth() {
  const isAuthenticated = computed(() => !!token.value && !!currentUser.value)
  const isAdmin = computed(() => currentUser.value?.role === 'admin')

  const fetchConfig = async () => {
    try {
      const res = await fetch('/api/config')
      if (res.ok) {
        const data = await res.json()
        appConfig.value = {
          google_client_id: data.google_client_id,
          is_authenticated: isAuthenticated.value,
          user: currentUser.value,
        }
        initGoogleClient()
      }
    } catch (err) {
      console.error('Failed to fetch config:', err)
    }
  }

  const initGoogleClient = () => {
    if (!appConfig.value?.google_client_id) return

    // Load GIS script if not present
    if (!document.getElementById('google-gis-script')) {
      const script = document.createElement('script')
      script.id = 'google-gis-script'
      script.src = 'https://accounts.google.com/gsi/client'
      script.async = true
      script.defer = true
      script.onload = () => setupTokenClient()
      document.head.appendChild(script)
    } else if ((window as any).google) {
      setupTokenClient()
    }
  }

  const setupTokenClient = () => {
    const google = (window as any).google
    if (!google || !appConfig.value?.google_client_id) return

    try {
      tokenClient = google.accounts.oauth2.initTokenClient({
        client_id: appConfig.value.google_client_id,
        scope: 'openid profile email https://www.googleapis.com/auth/drive.readonly',
        callback: async (resp: any) => {
          if (resp.error) {
            loginError.value = resp.error_description || resp.error
            isLoggingIn.value = false
            return
          }
          if (resp.access_token) {
            await handleGoogleAccessToken(resp.access_token)
          }
        },
      })
    } catch (e) {
      console.error('Error initTokenClient:', e)
    }
  }

  const handleGoogleAccessToken = async (gToken: string) => {
    isLoggingIn.value = true
    loginError.value = null
    try {
      const res = await fetch('/api/auth/google', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ access_token: gToken }),
      })

      const data = await res.json()
      if (!res.ok) {
        loginError.value = data.error || 'Login failed'
        return
      }

      token.value = data.token
      googleAccessToken.value = gToken
      currentUser.value = data.user

      localStorage.setItem('kame_jwt_token', data.token)
      localStorage.setItem('kame_google_access_token', gToken)
      localStorage.setItem('kame_user', JSON.stringify(data.user))
    } catch (err: any) {
      loginError.value = err.message || 'Network error during login'
    } finally {
      isLoggingIn.value = false
    }
  }

  const loginWithGoogle = () => {
    loginError.value = null
    if (!tokenClient) {
      initGoogleClient()
      setTimeout(() => {
        if (tokenClient) {
          tokenClient.requestAccessToken({ prompt: 'consent' })
        } else {
          loginError.value = 'Google OAuth client is not ready. Please check GOOGLE_CLIENT_ID.'
        }
      }, 500)
      return
    }
    tokenClient.requestAccessToken({ prompt: 'consent' })
  }

  const logout = () => {
    token.value = null
    googleAccessToken.value = null
    currentUser.value = null
    localStorage.removeItem('kame_jwt_token')
    localStorage.removeItem('kame_google_access_token')
    localStorage.removeItem('kame_user')
  }

  return {
    token,
    googleAccessToken,
    currentUser,
    appConfig,
    isAuthenticated,
    isAdmin,
    isLoggingIn,
    loginError,
    fetchConfig,
    loginWithGoogle,
    logout,
  }
}
