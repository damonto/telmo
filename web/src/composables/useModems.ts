import { ref } from 'vue'

import { modemFixtures } from '@/data/modems'
import type { Modem } from '@/types/modem'

const modems = ref<Modem[]>(modemFixtures)

const getModemById = (id: string) => modems.value.find((modem) => modem.id === id) ?? null

export const useModems = () => ({
  getModemById,
  modems,
})
