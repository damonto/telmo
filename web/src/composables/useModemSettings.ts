import { computed, ref, watch, type ComputedRef } from 'vue'
import { useI18n } from 'vue-i18n'

import { useModemApi } from '@/apis/modem'
import { useNetworkApi } from '@/apis/network'
import type { Modem, ModemSettings } from '@/types/modem'
import type { NetworkResponse } from '@/types/network'

export const useModemSettings = (modemId: ComputedRef<string>) => {
  const { t } = useI18n()
  const modemApi = useModemApi()
  const networkApi = useNetworkApi()

  const modem = ref<Modem | null>(null)
  const isModemLoading = ref(false)

  const msisdnInput = ref('')
  const isMsisdnUpdating = ref(false)

  const settingsAlias = ref('')
  const settingsMss = ref('')
  const settingsCompatible = ref(false)
  const isSettingsLoading = ref(false)
  const isSettingsUpdating = ref(false)

  const networkDialogOpen = ref(false)
  const availableNetworks = ref<NetworkResponse[]>([])
  const selectedNetwork = ref('')
  const isNetworkLoading = ref(false)
  const isNetworkRegistering = ref(false)

  const feedbackOpen = ref(false)
  const feedbackMessage = ref('')

  const msisdnValue = computed(() => msisdnInput.value.trim())
  const isMsisdnValid = computed(() => msisdnValue.value.length > 0)

  const mssValue = computed(() => Number.parseInt(settingsMss.value, 10))
  const isMssValid = computed(() => {
    return Number.isInteger(mssValue.value) && mssValue.value >= 64 && mssValue.value <= 254
  })

  const currentOperatorLabel = computed(() => {
    const unknown = t('modemDetail.settings.networkUnknown')
    const name = modem.value?.registeredOperator?.name?.trim() ?? ''
    const code = modem.value?.registeredOperator?.code?.trim() ?? ''
    if (!name && !code) return unknown
    if (name && code) return `${name} (${code})`
    return name || code || unknown
  })

  const currentRegistrationState = computed(() => {
    const value = modem.value?.registrationState?.trim()
    return value && value.length > 0 ? value : t('modemDetail.settings.networkUnknown')
  })

  const currentAccessTechnology = computed(() => {
    const value = modem.value?.accessTechnology?.trim()
    return value && value.length > 0 ? value : t('modemDetail.settings.networkUnknown')
  })

  const hasAvailableNetworks = computed(() => availableNetworks.value.length > 0)
  const hasNetworkSelection = computed(() => selectedNetwork.value.trim().length > 0)

  const fetchModem = async (id: string) => {
    if (isModemLoading.value) return
    isModemLoading.value = true
    try {
      const { data } = await modemApi.getModem(id)
      modem.value = data.value?.data ?? null
      msisdnInput.value = modem.value?.number ?? ''
    } finally {
      isModemLoading.value = false
    }
  }

  const fetchSettings = async (id: string) => {
    if (isSettingsLoading.value) return
    isSettingsLoading.value = true
    try {
      const { data } = await modemApi.getSettings(id)
      const payload: ModemSettings | undefined = data.value?.data
      settingsAlias.value = payload?.alias ?? ''
      settingsMss.value = payload?.mss ? String(payload.mss) : ''
      settingsCompatible.value = payload?.compatible ?? false
    } finally {
      isSettingsLoading.value = false
    }
  }

  const handleMsisdnUpdate = async () => {
    const targetId = modemId.value
    if (!targetId || targetId === 'unknown') return
    if (!isMsisdnValid.value || isMsisdnUpdating.value) return
    isMsisdnUpdating.value = true
    try {
      await modemApi.updateMsisdn(targetId, msisdnValue.value)
      await fetchModem(targetId)
      feedbackMessage.value = t('modemDetail.settings.msisdnSuccess')
      feedbackOpen.value = true
    } catch (err) {
      console.error('[useModemSettings] Failed to update MSISDN:', err)
    } finally {
      isMsisdnUpdating.value = false
    }
  }

  const openNetworkDialog = async () => {
    const targetId = modemId.value
    if (!targetId || targetId === 'unknown') return
    if (isNetworkLoading.value) return
    networkDialogOpen.value = true
    selectedNetwork.value = ''
    isNetworkLoading.value = true
    try {
      const { data } = await networkApi.scanNetworks(targetId)
      availableNetworks.value = data.value?.data ?? []
    } catch (err) {
      console.error('[useModemSettings] Failed to scan networks:', err)
      availableNetworks.value = []
    } finally {
      isNetworkLoading.value = false
    }
  }

  const handleNetworkRegister = async () => {
    const targetId = modemId.value
    if (!targetId || targetId === 'unknown') return
    if (!hasNetworkSelection.value || isNetworkRegistering.value) return
    isNetworkRegistering.value = true
    try {
      await networkApi.registerNetwork(targetId, selectedNetwork.value)
      await fetchModem(targetId)
      networkDialogOpen.value = false
      feedbackMessage.value = t('modemDetail.settings.networkSuccess')
      feedbackOpen.value = true
    } catch (err) {
      console.error('[useModemSettings] Failed to register network:', err)
    } finally {
      isNetworkRegistering.value = false
    }
  }

  const handleSettingsUpdate = async () => {
    const targetId = modemId.value
    if (!targetId || targetId === 'unknown') return
    if (!isMssValid.value || isSettingsUpdating.value) return
    isSettingsUpdating.value = true
    try {
      const payload: ModemSettings = {
        alias: settingsAlias.value.trim(),
        compatible: settingsCompatible.value,
        mss: mssValue.value,
      }
      await modemApi.updateSettings(targetId, payload)
      await fetchSettings(targetId)
      feedbackMessage.value = t('modemDetail.settings.deviceSuccess')
      feedbackOpen.value = true
    } catch (err) {
      console.error('[useModemSettings] Failed to update settings:', err)
    } finally {
      isSettingsUpdating.value = false
    }
  }

  watch(
    modemId,
    async (id) => {
      if (!id || id === 'unknown') return
      await Promise.all([fetchModem(id), fetchSettings(id)])
    },
    { immediate: true },
  )

  return {
    modem,
    isModemLoading,
    msisdnInput,
    isMsisdnUpdating,
    isMsisdnValid,
    settingsAlias,
    settingsMss,
    settingsCompatible,
    isSettingsLoading,
    isSettingsUpdating,
    isMssValid,
    currentOperatorLabel,
    currentRegistrationState,
    currentAccessTechnology,
    networkDialogOpen,
    availableNetworks,
    selectedNetwork,
    isNetworkLoading,
    isNetworkRegistering,
    hasAvailableNetworks,
    hasNetworkSelection,
    feedbackOpen,
    feedbackMessage,
    openNetworkDialog,
    handleMsisdnUpdate,
    handleNetworkRegister,
    handleSettingsUpdate,
  }
}
