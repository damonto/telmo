import { SignalHigh, SignalLow, SignalMedium, SignalZero } from 'lucide-vue-next'

export const useModemDisplay = () => {
  /**
   * Get signal icon based on quality percentage (0-100)
   */
  const signalIcon = (percentage: number) => {
    if (percentage >= 70) return SignalHigh
    if (percentage >= 50) return SignalMedium
    if (percentage >= 30) return SignalLow
    return SignalZero
  }

  /**
   * Get signal color class based on quality percentage (0-100)
   */
  const signalTone = (percentage: number) => {
    if (percentage >= 70) return 'text-emerald-500'
    if (percentage >= 50) return 'text-lime-500'
    if (percentage >= 30) return 'text-amber-500'
    return 'text-rose-500'
  }

  /**
   * Format signal quality as percentage
   */
  const formatSignal = (percentage: number) => `${Math.round(percentage)}%`

  const flagClass = (regionCode: string) => {
    const normalized = regionCode.trim().toLowerCase()
    if (!/^[a-z]{2}$/.test(normalized)) {
      return null
    }
    return `fi fi-${normalized}`
  }

  return {
    flagClass,
    formatSignal,
    signalIcon,
    signalTone,
  }
}
