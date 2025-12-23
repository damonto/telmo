// Backend API Response Types

// Common API response wrapper
export type ApiResponse<T = unknown> = {
  data: T
}

// Modem data types
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

// API response types with wrapper
export type ModemListResponse = ApiResponse<ModemApiResponse[]>
export type ModemDetailResponse = ApiResponse<ModemApiResponse>
export type EuiccDetailResponse = ApiResponse<EuiccApiResponse>

// Frontend uses ModemApiResponse directly as Modem type
export type Modem = ModemApiResponse

// eSIM Profile type (for future eSIM management features)
export type EsimProfile = {
  id: string
  name: string
  iccid: string
  enabled: boolean
  regionCode: string
  logoUrl?: string
}

// eSIM Chip Information
export type EuiccApiResponse = {
  eid: string
  freeSpace: number
  sasUp: string
  certificates: string[]
}

export type EuiccResponse = ApiResponse<EuiccApiResponse>
