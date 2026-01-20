import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL || 'http://localhost:9000',
})

api.interceptors.request.use(
  function (config) {
    const accessToken = localStorage.getItem('access_token')
    if (accessToken) {
      config.headers.Authorization = 'Bearer ' + accessToken
    }
    return config
  },
  function (error) {
    return Promise.reject(error)
  },
)

api.interceptors.response.use(
  function fulfilled(response) {
    return response
  },
  async function rejected(error) {
    const originalRequest = error.config

    // Handle 401 errors (token expired)
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const refreshToken = localStorage.getItem('refresh_token')
        if (refreshToken) {
          const response = await api.post('/api/v1/auth/refresh', {
            refreshToken: refreshToken,
          })

          const { accessToken, refreshToken: newRefreshToken } =
            response.data.data

          localStorage.setItem('access_token', accessToken)
          localStorage.setItem('refresh_token', newRefreshToken)

          // Retry the original request with new token
          originalRequest.headers.Authorization = `Bearer ${accessToken}`
          return api(originalRequest)
        }
      } catch (refreshError) {
        // Refresh failed, redirect to login
        // localStorage.removeItem('access_token')
        // localStorage.removeItem('refresh_token')
        // window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    console.log(error)
    return Promise.reject(error)
  },
)

// ============================================
// TYPE DEFINITIONS
// ============================================

export interface AuthTokens {
  accessToken: string
  refreshToken: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface UserResponse {
  id: string
  email: string
  poll_count: number
  profile_picture?: string
}

export interface CreatePollRequest {
  title: string
  options: Array<string>
  expires_at: string // ISO 8601 date string
}

export interface Vote {
  id: string
  user_id: string
  poll_id: string
}

export interface PollOption {
  id: string
  name: string
  votes: Array<Vote>
}

export interface PollViewResponse {
  id: string
  title: string
  options: Array<PollOption>
  created_at: string
  expires_at: string
  creator_id: string
  voted: boolean
}

export interface VoteRequest {
  poll_id: string
  option_id: string
}

export interface VoteResponse {
  message: string
}

// ============================================
// API FUNCTIONS
// ============================================

// Auth APIs
export const authAPI = {
  register: async (data: RegisterRequest) => {
    const response = await api.post<{ message: string; data: AuthTokens }>(
      '/api/v1/auth/register',
      data,
    )
    return response.data
  },

  login: async (data: LoginRequest) => {
    const response = await api.post<{ message: string; data: AuthTokens }>(
      '/api/v1/auth/login',
      data,
    )
    return response.data
  },

  refreshToken: async (data: RefreshTokenRequest) => {
    const response = await api.post<{ message: string; data: AuthTokens }>(
      '/api/v1/auth/refresh',
      data,
    )
    return response.data
  },
}

// User APIs
export const userAPI = {
  getUser: async (userId: string) => {
    const response = await api.get<{
      message: string
      data: UserResponse
    }>(`/api/v1/user/${userId}`)
    return response.data.data
  },
}

// Poll APIs
export const pollAPI = {
  createPoll: async (data: CreatePollRequest) => {
    const response = await api.post<{ message: string }>(
      '/api/v1/poll/create',
      data,
    )
    return response.data
  },

  deletePoll: async (pollId: string) => {
    const response = await api.post<{ message: string }>(
      `/api/v1/poll/${pollId}`,
    )
    return response.data
  },

  getPollView: async (pollId: string) => {
    const response = await api.get<PollViewResponse>(
      `/api/v1/poll/view/${pollId}`,
    )
    return response.data
  },

  getAllPolls: async () => {
    const response = await api.get<{
      message: string
      data: Array<PollViewResponse>
    }>('/api/v1/poll/all')
    return response.data.data
  },

  getPoll: async (pollId: string) => {
    const response = await api.get<PollViewResponse>(`/api/v1/poll/${pollId}`)
    return response.data
  },
}

// Vote APIs
// Note: These are commented out in your backend (app/app.go lines 89-129)
// To enable, uncomment the vote routes in your backend
export const voteAPI = {
  vote: async (data: VoteRequest) => {
    const response = await api.post<VoteResponse>('/api/v1/vote/', data)
    return response.data
  },
}

export default api
