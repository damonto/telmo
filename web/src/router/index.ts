import { createRouter, createWebHistory } from 'vue-router'

import { getStoredToken } from '@/lib/auth-storage'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/auth',
      name: 'auth',
      component: () => import('@/views/AuthVerifyView.vue'),
    },
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
    },
    {
      path: '/messages',
      name: 'messages',
      component: () => import('@/views/MessagesView.vue'),
    },
    {
      path: '/modems/:id',
      component: () => import('@/layouts/ModemLayout.vue'),
      children: [
        {
          path: '',
          name: 'modem-detail',
          component: () => import('@/views/ModemDetailView.vue'),
        },
        {
          path: 'messages',
          name: 'modem-messages',
          component: () => import('@/views/ModemMessagesView.vue'),
        },
        {
          path: 'messages/:participant',
          name: 'modem-message-thread',
          component: () => import('@/views/ModemMessageThreadView.vue'),
        },
        {
          path: 'ussd',
          name: 'modem-ussd',
          component: () => import('@/views/ModemUssdView.vue'),
        },
        {
          path: 'settings',
          name: 'modem-settings',
          component: () => import('@/views/ModemSettingsView.vue'),
        },
      ],
    },
  ],
})

const AUTH_ROUTE_NAME = 'auth'

router.beforeEach((to) => {
  const token = getStoredToken()
  if (!token && to.name !== AUTH_ROUTE_NAME) {
    return { name: AUTH_ROUTE_NAME }
  }

  if (token && to.name === AUTH_ROUTE_NAME) {
    return { name: 'home' }
  }
})

export default router
