// Backend API Response Types (Direct mapping)

export type SimInfo = {
  operatorName: string
  operatorIdentifier: string
  regionCode: string
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
  sim: SimInfo
  accessTechnology: string | null
  registrationState: string
  registeredOperator: RegisteredOperator
  signalQuality: number
  isEsim: boolean
}

export type ModemListResponse = ModemApiResponse[]

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
