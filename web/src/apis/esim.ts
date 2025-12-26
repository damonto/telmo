import { useFetch } from '@/lib/fetch'

import type { ApiResponse, EmptyObject } from '@/types/api'
import type { EsimProfilesResponse } from '@/types/esim'

export const useEsimApi = () => {
  const getEsims = (id: string) => {
    return useFetch<EsimProfilesResponse>(`modems/${id}/esims`).get().json()
  }

  const updateEsimNickname = (id: string, iccid: string, nickname: string) => {
    return useFetch<ApiResponse<EmptyObject>>(`modems/${id}/esims/${iccid}/nickname`, {
      method: 'PUT',
      body: JSON.stringify({ nickname }),
    }).json()
  }

  const enableEsim = (id: string, iccid: string) => {
    return useFetch<ApiResponse<EmptyObject>>(`modems/${id}/esims/${iccid}/enabling`, {
      method: 'POST',
    }).json()
  }

  const deleteEsim = (id: string, iccid: string) => {
    return useFetch<ApiResponse<EmptyObject>>(`modems/${id}/esims/${iccid}`, {
      method: 'DELETE',
    }).json()
  }

  return {
    getEsims,
    updateEsimNickname,
    enableEsim,
    deleteEsim,
  }
}
