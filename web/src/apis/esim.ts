import { useFetch } from '@/lib/fetch'

import type { EsimProfilesResponse } from '@/types/esim'

export const useEsimApi = () => {
  const getEsims = (id: string) => {
    return useFetch<EsimProfilesResponse>(`modems/${id}/esims`).get().json()
  }

  const updateEsimNickname = (id: string, iccid: string, nickname: string) => {
    return useFetch<string>(`modems/${id}/esims/${iccid}/nickname`, {
      method: 'PUT',
      body: JSON.stringify({ nickname }),
    }).text()
  }

  const enableEsim = (id: string, iccid: string) => {
    return useFetch<string>(`modems/${id}/esims/${iccid}/enabling`, {
      method: 'POST',
    }).text()
  }

  const deleteEsim = (id: string, iccid: string) => {
    return useFetch<string>(`modems/${id}/esims/${iccid}`, {
      method: 'DELETE',
    }).text()
  }

  return {
    getEsims,
    updateEsimNickname,
    enableEsim,
    deleteEsim,
  }
}
