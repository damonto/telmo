import { useFetch } from '@/lib/fetch'

import type { NotificationsResponse } from '@/types/notification'

export const useNotificationApi = () => {
  const getNotifications = (id: string) => {
    return useFetch<NotificationsResponse>(`modems/${id}/notifications`).get().json()
  }

  const resendNotification = (id: string, sequence: string) => {
    return useFetch<string>(`modems/${id}/notifications/${sequence}/resend`, {
      method: 'POST',
    })
  }

  const deleteNotification = (id: string, sequence: string) => {
    return useFetch<string>(`modems/${id}/notifications/${sequence}`, {
      method: 'DELETE',
    })
  }

  return {
    getNotifications,
    resendNotification,
    deleteNotification,
  }
}
