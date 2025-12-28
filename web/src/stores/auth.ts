import { defineStore } from 'pinia'

import { useAuthApi } from '@/apis/auth'
import { clearStoredToken, getStoredToken, setStoredToken } from '@/lib/auth-storage'

const RESEND_COOLDOWN_MS = 60_000
const CODE_LENGTH = 6

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: getStoredToken(),
    isSending: false,
    isVerifying: false,
    resendAvailableAt: 0,
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
  },
  actions: {
    setToken(token: string) {
      this.token = token
      setStoredToken(token)
    },
    clearToken() {
      this.token = null
      clearStoredToken()
    },
    async sendCode() {
      if (this.isSending) return
      if (this.resendAvailableAt > Date.now()) return

      this.isSending = true
      try {
        const { error } = await useAuthApi().sendCode()
        if (error.value) return

        this.resendAvailableAt = Date.now() + RESEND_COOLDOWN_MS
      } finally {
        this.isSending = false
      }
    },
    async verifyCode(code: string) {
      if (this.isVerifying) return null
      if (code.trim().length !== CODE_LENGTH) return null

      this.isVerifying = true
      try {
        const { data, error } = await useAuthApi().verifyCode({ code })
        if (error.value) return null

        const token = data.value?.data?.token
        if (!token) return null

        this.setToken(token)
        return token
      } finally {
        this.isVerifying = false
      }
    },
  },
})
