import { useFetch } from '@/lib/fetch'

import type { NetworksResponse } from '@/types/network'

export const useNetworkApi = () => {
  const scanNetworks = (id: string) => {
    return useFetch<NetworksResponse>(`modems/${id}/networks`).get().json()
  }

  const registerNetwork = (id: string, operatorCode: string) => {
    const encoded = encodeURIComponent(operatorCode)
    return useFetch<string>(`modems/${id}/networks/${encoded}`, {
      method: 'PUT',
    }).text()
  }

  return {
    scanNetworks,
    registerNetwork,
  }
}
