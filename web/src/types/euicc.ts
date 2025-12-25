import type { ApiResponse } from '@/types/api'

export type EuiccApiResponse = {
  eid: string
  freeSpace: number
  sasUp: string
  certificates: string[]
}

export type EuiccDetailResponse = ApiResponse<EuiccApiResponse>
export type EuiccResponse = ApiResponse<EuiccApiResponse>
