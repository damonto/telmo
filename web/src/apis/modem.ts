import { useFetch } from '@/lib/fetch'

import type {
  ModemDetailResponse,
  ModemListResponse,
  ModemSettings,
  ModemSettingsResponse,
} from '@/types/modem'

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

  /**
   * Switch active SIM slot by identifier
   * PUT /api/v1/modems/:id/sim-slots/:identifier
   */
  const switchSimSlot = (id: string, identifier: string) => {
    return useFetch<string>(`modems/${id}/sim-slots/${identifier}`, {
      method: 'PUT',
    })
  }

  /**
   * Update MSISDN
   * PUT /api/v1/modems/:id/msisdn
   */
  const updateMsisdn = (id: string, number: string) => {
    return useFetch<string>(`modems/${id}/msisdn`, {
      method: 'PUT',
      body: JSON.stringify({ number }),
    })
  }

  /**
   * Fetch modem settings
   * GET /api/v1/modems/:id/settings
   */
  const getSettings = (id: string) => {
    return useFetch<ModemSettingsResponse>(`modems/${id}/settings`).get().json()
  }

  /**
   * Update modem settings
   * PUT /api/v1/modems/:id/settings
   */
  const updateSettings = (id: string, payload: ModemSettings) => {
    return useFetch<string>(`modems/${id}/settings`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  }

  return {
    getModems,
    getModem,
    switchSimSlot,
    updateMsisdn,
    getSettings,
    updateSettings,
  }
}
