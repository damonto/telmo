import { useFetch } from '@/lib/fetch'
import type { EmptyObject } from '@/types/api'
import type { AuthVerifyPayload, AuthVerifyResponse } from '@/types/auth'

/**
 * Auth API
 * Centralized API definitions
 */
export const useAuthApi = () => {
  /**
   * Send verification code
   * TODO: replace endpoint once the auth API is finalized.
   * POST /api/v1/auth/otp
   */
  const sendCode = () => {
    return useFetch<EmptyObject>('auth/otp', {
      method: 'POST',
    })
  }

  /**
   * Verify verification code
   * TODO: replace endpoint once the auth API is finalized.
   * POST /api/v1/auth/otp/verify
   */
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
