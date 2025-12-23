import { createFetch } from '@vueuse/core'

import { handleError } from './error-handler'

/**
 * Custom fetch instance with global configuration
 * Unified error handling - no need to handle errors in callers
 */
export const useFetch = createFetch({
  baseUrl: import.meta.env.VITE_API_BASE_URL || 'http://10.10.10.101:9527/api/v1',
  options: {
    async beforeFetch({ options }) {
      // Add authentication token if available
      const token = localStorage.getItem('auth_token')
      if (token) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${token}`,
        }
      }

      // Add custom headers
      options.headers = {
        ...options.headers,
        'Content-Type': 'application/json',
      }

      return { options }
    },

    async afterFetch({ response, data }) {
      // Unified logging in development
      if (import.meta.env.DEV) {
        console.log(`[API] ${response.url} → ${response.status}`, data)
      }

      // Unified error handling for non-2xx responses
      if (!response.ok) {
        const errorData = data as { code: number; message: string }
        handleError(errorData)
      }

      return { response, data }
    },

    async onFetchError({ response, error }) {
      // Unified network error handling
      const errorMessage = error?.message || 'Network error occurred'
      handleError(errorMessage)
      console.error('[API] Network error:', errorMessage)

      // Handle 401 unauthorized
      if (response?.status === 401) {
        localStorage.removeItem('auth_token')
        // Optionally redirect to login
        // window.location.href = '/login'
      }

      // Don't throw - let callers handle error via error.value
      return { response, error }
    },
  },
  fetchOptions: {
    mode: 'cors',
  },
})
