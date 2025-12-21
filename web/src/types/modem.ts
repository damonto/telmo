export type EsimProfile = {
  id: string
  name: string
  iccid: string
  enabled: boolean
  regionCode: string
  logoUrl?: string
}

type ModemBase = {
  id: string
  name: string
  manufacturer: string
  carrierName: string
  roamingCarrierName?: string
  regionCode: string
  isRoaming: boolean
  signalDbm: number
  isESim: boolean
  tech: '4G' | '5G'
  logoUrl?: string
}

export type PhysicalModem = ModemBase & {
  isESim: false
  iccid: string
}

export type EsimModem = ModemBase & {
  isESim: true
  imei: string
  eid: string
  storageRemaining: string
  profiles: EsimProfile[]
}

export type Modem = PhysicalModem | EsimModem
