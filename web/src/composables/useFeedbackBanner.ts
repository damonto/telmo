import { toast } from 'vue-sonner'

export const useFeedbackBanner = () => {
  const showFeedback = (message: string) => {
    const trimmed = message.trim()
    if (!trimmed) return
    toast.success(trimmed)
  }

  return {
    showFeedback,
  }
}
