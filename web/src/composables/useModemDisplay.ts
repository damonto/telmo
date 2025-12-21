import { SignalHigh, SignalLow, SignalMedium, SignalZero } from 'lucide-vue-next'

export const useModemDisplay = () => {
  const signalIcon = (dbm: number) => {
    if (dbm >= -70) return SignalHigh
    if (dbm >= -85) return SignalMedium
    if (dbm >= -100) return SignalLow
    return SignalZero
  }

  const signalTone = (dbm: number) => {
    if (dbm >= -70) return 'text-emerald-500'
    if (dbm >= -85) return 'text-lime-500'
    if (dbm >= -100) return 'text-amber-500'
    return 'text-rose-500'
  }

  const formatSignal = (dbm: number) => `${dbm} dBm`

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
