import { useFetch } from '@/lib/fetch'

import type { ModemDetailResponse, ModemListResponse } from '@/types/modem'

/**
 * Modem API
 * Centralized API definitions
 */
export const useModemApi = () => {
  /**
   * Fetch all modems
   * GET /api/v1/modems
   */
  const getModems = () => {
    return useFetch<ModemListResponse>('modems').get().json()
  }

  /**
   * Fetch single modem by ID
   * GET /api/v1/modems/:id
   */
  const getModem = (id: string) => {
    return useFetch<ModemDetailResponse>(`modems/${id}`).get().json()
  }

  return {
    getModems,
    getModem,
  }
}
