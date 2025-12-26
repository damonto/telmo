import { createFetch } from '@vueuse/core'

import { handleError, handleResponseError } from './error-handler'

const rawBaseUrl = import.meta.env.VITE_API_BASE_URL
const baseUrl =
  rawBaseUrl && rawBaseUrl.trim().length > 0 ? rawBaseUrl.replace(/\/$/, '') : '/api/v1'

/**
 * Custom fetch instance with global configuration
 * Unified error handling - no need to handle errors in callers
 */
export const useFetch = createFetch({
  baseUrl,
  options: {
    updateDataOnError: true,
    async beforeFetch({ options }) {
      const headers = new Headers(options.headers)

      // Add authentication token if available
      const token = localStorage.getItem('auth_token')
      if (token) {
        headers.set('Authorization', `Bearer ${token}`)
      }

      const hasBody = options.body !== undefined && options.body !== null
      if (hasBody && !(options.body instanceof FormData) && !headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json')
      }

      options.headers = headers

      return { options }
    },

    afterFetch({ response, data }) {
      // Unified logging in development
      if (import.meta.env.DEV) {
        console.log(`[API] ${response.url} → ${response.status}`, data)
      }

      return { response, data }
    },

    onFetchError({ response, error, data }) {
      if (response) {
        handleResponseError(response, data)
        console.error('[API] Response error:', response.status, data)
      } else {
        // Unified network error handling
        handleError(error, 'Network error occurred')
        console.error('[API] Network error:', error)
      }

      // Don't throw - let callers handle error via error.value
      return { response, error, data: null }
    },
  },
  fetchOptions: {
    mode: 'cors',
  },
})
