export interface ApiError {
  code: number
  message: string
}

/**
 * Global error handler interface
 * Will be initialized by useErrorHandler
 */
let showErrorFunction: ((message: string, title?: string) => void) | null = null

/**
 * Register the global error handler
 * Called from App.vue to initialize error handling
 */
export const registerErrorHandler = (showError: (message: string, title?: string) => void) => {
  showErrorFunction = showError
}

/**
 * Error handler for API errors
 * Shows error messages to users
 */
export const handleError = (error: unknown, defaultMessage = 'An error occurred') => {
  let message = defaultMessage
  let title = 'Error'

  if (error && typeof error === 'object') {
    // Handle API error format { code, message }
    if ('message' in error && typeof error.message === 'string') {
      message = error.message
    } else if ('code' in error && 'message' in error) {
      const apiError = error as ApiError
      message = apiError.message
      // Set title based on status code
      if (apiError.code === 401) title = 'Unauthorized'
      else if (apiError.code === 404) title = 'Not Found'
      else if (apiError.code >= 500) title = 'Server Error'
    }
  } else if (typeof error === 'string') {
    message = error
  }

  if (showErrorFunction) {
    showErrorFunction(message, title)
  } else {
    // Fallback to alert if handler not registered
    alert(message)
  }

  // Log to console for debugging
  console.error('[Error]', error)
}

/**
 * Handle API response errors
 */
export const handleResponseError = (response: Response, data?: unknown) => {
  let message = `Error ${response.status}: ${response.statusText}`
  let title = 'Error'

  // Try to extract message from response data
  if (data && typeof data === 'object' && 'message' in data) {
    message = (data as ApiError).message
  }

  // Handle specific status codes
  if (response.status === 401) {
    title = 'Unauthorized'
    message = 'Unauthorized - Please login again'
    localStorage.removeItem('auth_token')
  } else if (response.status === 404) {
    title = 'Not Found'
    message =
      data && typeof data === 'object' && 'message' in data
        ? (data as ApiError).message
        : 'Resource not found'
  } else if (response.status >= 500) {
    title = 'Server Error'
    message =
      data && typeof data === 'object' && 'message' in data
        ? (data as ApiError).message
        : 'Server error - Please try again later'
  }

  if (showErrorFunction) {
    showErrorFunction(message, title)
  } else {
    // Fallback to alert if handler not registered
    alert(`${title}: ${message}`)
  }

  console.error('[API Error]', title, message, data)
}
