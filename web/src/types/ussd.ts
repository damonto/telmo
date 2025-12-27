import type { ApiResponse } from '@/types/api'

export type UssdAction = 'initialize' | 'reply'

export type UssdReply = {
  reply: string
}

export type UssdExecuteResponse = ApiResponse<UssdReply>
