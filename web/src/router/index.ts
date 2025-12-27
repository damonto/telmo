import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
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

export default router
