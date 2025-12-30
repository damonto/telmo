import { ref } from 'vue'

export const useFeedbackBanner = () => {
  const feedbackOpen = ref(false)
  const feedbackMessage = ref('')

  const showFeedback = (message: string) => {
    feedbackMessage.value = message
    feedbackOpen.value = true
  }

  const clearFeedback = () => {
    feedbackOpen.value = false
    feedbackMessage.value = ''
  }

  return {
    feedbackOpen,
    feedbackMessage,
    showFeedback,
    clearFeedback,
  }
}
