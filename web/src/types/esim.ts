import type { ApiResponse } from '@/types/api'

export type EsimProfileApiResponse = {
  name: string
  serviceProviderName: string
  iccid: string
  icon: string
  profileState: number
  regionCode?: string
}

export type EsimProfilesResponse = ApiResponse<EsimProfileApiResponse[]>

export type EsimProfile = {
  id: string
  name: string
  iccid: string
  enabled: boolean
  regionCode: string
  logoUrl?: string
}
