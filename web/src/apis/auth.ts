import { useFetch } from '@/lib/fetch'
import type { EmptyObject } from '@/types/api'
import type { AuthOtpRequirementResponse, AuthVerifyPayload, AuthVerifyResponse } from '@/types/auth'

export const useAuthApi = () => {
  const sendCode = () => {
    return useFetch<EmptyObject>('auth/otp', {
      method: 'POST',
    })
  }

  const getOtpRequirement = () => {
    return useFetch<AuthOtpRequirementResponse>('auth/otp/required').get().json()
  }

  const verifyCode = (payload: AuthVerifyPayload) => {
    return useFetch<AuthVerifyResponse>('auth/otp/verify', {
      method: 'POST',
      body: JSON.stringify(payload),
    }).json()
  }

  return {
    getOtpRequirement,
    sendCode,
    verifyCode,
  }
}
