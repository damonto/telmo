import { useFetch } from '@/lib/fetch'

import type { EuiccDetailResponse } from '@/types/euicc'

export const useEuiccApi = () => {
  const getEuicc = (id: string) => {
    return useFetch<EuiccDetailResponse>(`modems/${id}/euicc`).get().json()
  }

  return {
    getEuicc,
  }
}
