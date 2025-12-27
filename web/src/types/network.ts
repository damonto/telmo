import type { ApiResponse } from '@/types/api'

export type NetworkResponse = {
  status: string
  operatorName: string
  operatorShortName: string
  operatorCode: string
  accessTechnologies: string[]
}

export type NetworksResponse = ApiResponse<NetworkResponse[]>
