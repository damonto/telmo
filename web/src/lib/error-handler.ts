export interface ApiError {
  code: number
  message: string
}

/**
 * Global error handler interface
 * Will be initialized by useErrorHandler
 */
let showErrorFunction: ((message: string, title?: string) => void) | null = null

const resolveErrorInfo = (error: unknown, defaultMessage: string) => {
  let message = defaultMessage
  let title = 'Error'

  if (error && typeof error === 'object') {
    if ('code' in error && 'message' in error && typeof (error as ApiError).message === 'string') {
      const apiError = error as ApiError
      message = apiError.message
      if (apiError.code === 401) title = 'Unauthorized'
      else if (apiError.code === 404) title = 'Not Found'
      else if (apiError.code >= 500) title = 'Server Error'
    } else if ('message' in error && typeof error.message === 'string') {
      message = error.message
    }
  } else if (typeof error === 'string') {
    message = error
  }

  return { message, title }
}

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
  const { message, title } = resolveErrorInfo(error, defaultMessage)

  if (showErrorFunction) {
    showErrorFunction(message, title)
  }

  // Log to console for debugging
  console.error('[Error]', error)
}

/**
 * Handle API response errors
 */
export const handleResponseError = (response: Response, data?: unknown) => {
  let message = `Error ${response.status}: ${response.statusText}`

  // Try to extract message from response data
  if (data && typeof data === 'object' && 'message' in data && typeof data.message === 'string') {
    message = data.message
  }

  if (response.status === 401) {
    localStorage.removeItem('auth_token')
  }

  handleError({ code: response.status, message })
  console.error('[API Error]', message, data)
}
