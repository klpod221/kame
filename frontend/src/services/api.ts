export interface RequestConfig extends RequestInit {
  params?: Record<string, string | number | boolean | undefined>
}

export class ApiError extends Error {
  status: number
  data?: any

  constructor(status: number, message: string, data?: any) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }
}

class ApiClient {
  private getHeaders(customHeaders?: HeadersInit): Headers {
    const headers = new Headers(customHeaders)
    if (!headers.has('Content-Type')) {
      headers.set('Content-Type', 'application/json')
    }

    const token = localStorage.getItem('kame_jwt_token')
    if (token && !headers.has('Authorization')) {
      headers.set('Authorization', `Bearer ${token}`)
    }

    return headers
  }

  private buildUrl(endpoint: string, params?: Record<string, string | number | boolean | undefined>): string {
    const url = new URL(endpoint, window.location.origin)
    if (params) {
      Object.entries(params).forEach(([key, val]) => {
        if (val !== undefined && val !== null && val !== '') {
          url.searchParams.append(key, String(val))
        }
      })
    }
    return url.pathname + url.search
  }

  private async request<T>(endpoint: string, config: RequestConfig = {}): Promise<T> {
    const { params, headers, ...customConfig } = config
    const targetUrl = this.buildUrl(endpoint, params)

    const response = await fetch(targetUrl, {
      ...customConfig,
      headers: this.getHeaders(headers),
    })

    // Handle 401 Unauthorized (session expired)
    if (response.status === 401) {
      localStorage.removeItem('kame_jwt_token')
      localStorage.removeItem('kame_google_access_token')
      localStorage.removeItem('kame_user')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
      throw new ApiError(401, 'Session expired. Please log in again.')
    }

    // Handle 429 Too Many Requests
    if (response.status === 429) {
      throw new ApiError(429, 'Rate limit exceeded. Please wait a moment before trying again.')
    }

    let responseData: any = null
    const contentType = response.headers.get('content-type')
    if (contentType && contentType.includes('application/json')) {
      responseData = await response.json().catch(() => null)
    }

    if (!response.ok) {
      const errorMessage = responseData?.error || `Request failed with status ${response.status}`
      throw new ApiError(response.status, errorMessage, responseData)
    }

    return responseData as T
  }

  public get<T>(endpoint: string, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, { ...config, method: 'GET' })
  }

  public post<T>(endpoint: string, body?: any, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, {
      ...config,
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  public put<T>(endpoint: string, body?: any, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, {
      ...config,
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  public delete<T>(endpoint: string, config?: RequestConfig): Promise<T> {
    return this.request<T>(endpoint, { ...config, method: 'DELETE' })
  }
}

export const api = new ApiClient()
