import { onMounted, ref } from 'vue'

import { registerErrorHandler } from '@/lib/error-handler'

export interface ErrorOptions {
  title?: string
  message: string
}

// Global error state (singleton)
const errorTitle = ref('')
const errorMessage = ref('')
const isErrorOpen = ref(false)

/**
 * Show error dialog (internal function)
 */
const showAlertDialog = (message: string, title = 'Error') => {
  errorTitle.value = title
  errorMessage.value = message
  isErrorOpen.value = true
}

/**
 * Clear error dialog (internal function)
 */
const clearAlertDialog = () => {
  isErrorOpen.value = false
}

/**
 * Global error handler composable
 * Shows error dialogs using ErrorAlert component
 */
export const useErrorHandler = () => {
  /**
   * Show error dialog
   */
  const showError = (options: ErrorOptions | string, title?: string) => {
    if (typeof options === 'string') {
      showAlertDialog(options, title || 'Error')
    } else {
      showAlertDialog(options.message, options.title || 'Error')
    }
  }

  /**
   * Clear error dialog
   */
  const clearError = () => {
    clearAlertDialog()
  }

  // Register global error handler on mount (only once)
  onMounted(() => {
    registerErrorHandler((message: string, title?: string) => {
      showAlertDialog(message, title || 'Error')
    })
  })

  return {
    // State
    errorTitle,
    errorMessage,
    isErrorOpen,

    // Actions
    showError,
    clearError,
  }
}

// Export internal functions for direct access from lib/error-handler
export { clearAlertDialog }
