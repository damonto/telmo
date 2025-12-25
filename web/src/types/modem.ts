import type { ApiResponse } from '@/types/api'

export type SlotInfo = {
  active: boolean
  operatorName: string
  operatorIdentifier: string
  regionCode: string
  identifier: string
}

export type RegisteredOperator = {
  name: string
  code: string
}

export type ModemApiResponse = {
  manufacturer: string
  id: string
  firmwareRevision: string
  hardwareRevision: string
  name: string
  number: string
  sim: SlotInfo
  slots: SlotInfo[]
  accessTechnology: string | null
  registrationState: string
  registeredOperator: RegisteredOperator
  signalQuality: number
  supportsEsim: boolean
}

export type ModemListResponse = ApiResponse<ModemApiResponse[]>
export type ModemDetailResponse = ApiResponse<ModemApiResponse>

export type Modem = ModemApiResponse
