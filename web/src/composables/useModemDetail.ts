import { computed, ref } from 'vue'

import { useModemApi } from '@/apis/modem'
import type { EsimProfile, EsimProfileApiResponse, EuiccApiResponse, Modem } from '@/types/modem'

export const useModemDetail = () => {
  const modemApi = useModemApi()

  const modem = ref<Modem | null>(null)
  const euicc = ref<EuiccApiResponse | null>(null)
  const esimProfiles = ref<EsimProfile[]>([])
  const isLoading = ref(false)
  const isEuiccLoading = ref(false)
  const isEsimProfilesLoading = ref(false)
  const error = ref<string | null>(null)

  const mapEsimProfile = (profile: EsimProfileApiResponse): EsimProfile => {
    return {
      id: profile.iccid,
      name: profile.name,
      iccid: profile.iccid,
      enabled: profile.profileState === 1,
      regionCode: profile.regionCode ?? '',
      logoUrl: profile.icon.length > 0 ? profile.icon : undefined,
    }
  }

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

  const fetchEsimProfiles = async (id: string) => {
    isEsimProfilesLoading.value = true
    try {
      const { data } = await modemApi.getEsims(id)
      if (data.value?.data) {
        esimProfiles.value = data.value.data.map(mapEsimProfile)
      } else {
        esimProfiles.value = []
      }
    } catch (err) {
      console.error('[useModemDetail] Failed to fetch eSIM profiles:', err)
      esimProfiles.value = []
    } finally {
      isEsimProfilesLoading.value = false
    }
  }

  const fetchModemDetail = async (id: string) => {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null
    modem.value = null
    euicc.value = null
    esimProfiles.value = []

    try {
      const { data } = await modemApi.getModem(id)

      if (data.value?.data) {
        modem.value = data.value.data

        if (modem.value?.supportsEsim) {
          void fetchEuicc(id)
          void fetchEsimProfiles(id)
        }
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch modem details'
      console.error('[useModemDetail] Error:', err)
    } finally {
      isLoading.value = false
    }
  }

  return {
    modem,
    euicc,
    esimProfiles,
    isLoading,
    isEuiccLoading,
    isEsimProfilesLoading,
    error,
    hasModem: computed(() => modem.value !== null),
    isPhysicalModem: computed(() => Boolean(modem.value && !modem.value.supportsEsim)),
    isEsimModem: computed(() => Boolean(modem.value && modem.value.supportsEsim)),
    fetchModemDetail,
    fetchEuicc,
    fetchEsimProfiles,
  }
}
