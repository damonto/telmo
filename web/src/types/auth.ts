import type { ApiResponse } from '@/types/api'

export type AuthVerifyResponse = ApiResponse<{
  token: string
}>

export type AuthVerifyPayload = {
  code: string
}
