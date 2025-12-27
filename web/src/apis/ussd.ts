import { useFetch } from '@/lib/fetch'

import type { UssdAction, UssdExecuteResponse } from '@/types/ussd'

export const useUssdApi = () => {
  const executeUssd = (id: string, action: UssdAction, code: string) => {
    return useFetch<UssdExecuteResponse>(`modems/${id}/ussd`, {
      method: 'POST',
      body: JSON.stringify({ action, code }),
    }).json()
  }

  return {
    executeUssd,
  }
}
