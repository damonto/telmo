import type { ApiResponse } from '@/types/api'

export type AuthVerifyResponse = ApiResponse<{
  token: string
}>

export type AuthOtpRequirementResponse = ApiResponse<{
  otpRequired: boolean
}>

export type AuthVerifyPayload = {
  code: string
}
