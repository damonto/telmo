import { useFetch } from '@/lib/fetch'
import type { EmptyObject } from '@/types/api'
import type { AuthVerifyPayload, AuthVerifyResponse } from '@/types/auth'

export const useAuthApi = () => {
  const sendCode = () => {
    return useFetch<EmptyObject>('auth/otp', {
      method: 'POST',
    })
  }

  const verifyCode = (payload: AuthVerifyPayload) => {
    return useFetch<AuthVerifyResponse>('auth/otp/verify', {
      method: 'POST',
      body: JSON.stringify(payload),
    }).json()
  }

  return {
    sendCode,
    verifyCode,
  }
}
