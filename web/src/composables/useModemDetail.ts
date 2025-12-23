import { computed, ref } from 'vue'

import { useModemApi } from '@/apis/modem'
import type { EuiccApiResponse, Modem } from '@/types/modem'

/**
 * Composable for fetching and managing single modem details
 */
export const useModemDetail = () => {
  const modemApi = useModemApi()

  // State
  const modem = ref<Modem | null>(null)
  const euicc = ref<EuiccApiResponse | null>(null)
  const isLoading = ref(false)
  const isEuiccLoading = ref(false)
  const error = ref<string | null>(null)

  /**
   * Fetch modem details by ID
   */
  const fetchModemDetail = async (id: string) => {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      const { data } = await modemApi.getModem(id)

      if (data.value) {
        modem.value = data.value.data

        // If modem supports eSIM, fetch eUICC information
        if (modem.value?.supportsEsim) {
          await fetchEuicc(id)
        }
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch modem details'
      console.error('[useModemDetail] Error:', err)
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Fetch eUICC information
   */
  const fetchEuicc = async (id: string) => {
    isEuiccLoading.value = true

    try {
      const { data } = await modemApi.getEuicc(id)

      if (data.value) {
        euicc.value = data.value.data
      }
    } catch (err) {
      console.error('[useModemDetail] Failed to fetch eUICC info:', err)
      euicc.value = null
    } finally {
      isEuiccLoading.value = false
    }
  }

  return {
    // State
    modem,
    euicc,
    isLoading,
    isEuiccLoading,
    error,

    // Computed
    hasModem: computed(() => modem.value !== null),
    isPhysicalModem: computed(() => modem.value && !modem.value.supportsEsim),
    isEsimModem: computed(() => modem.value && modem.value.supportsEsim),

    // Actions
    fetchModemDetail,
    fetchEuicc,
  }
}
