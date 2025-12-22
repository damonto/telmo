import { useFetch } from '@/lib/fetch'

import type { ModemListResponse } from '@/types/modem'

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

  return {
    getModems,
  }
}
