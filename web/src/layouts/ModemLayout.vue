<script setup lang="ts">
import { Info, MessageSquare, Phone } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterView, useRoute } from 'vue-router'

import BottomTabBar from '@/components/BottomTabBar.vue'

const route = useRoute()
const { t } = useI18n()

const modemId = computed(() => (route.params.id ?? 'unknown') as string)

const tabItems = computed(() => [
  {
    key: 'detail',
    routeName: 'modem-detail',
    to: { name: 'modem-detail', params: { id: modemId.value } },
    label: t('modemDetail.tabs.detail'),
    icon: Info,
  },
  {
    key: 'messages',
    routeName: 'modem-messages',
    to: { name: 'modem-messages', params: { id: modemId.value } },
    label: t('modemDetail.tabs.messages'),
    icon: MessageSquare,
  },
  {
    key: 'ussd',
    routeName: 'modem-ussd',
    to: { name: 'modem-ussd', params: { id: modemId.value } },
    label: t('modemDetail.tabs.ussd'),
    icon: Phone,
  },
])
</script>

<template>
  <div class="min-h-screen bg-linear-to-b from-background via-background to-muted/40">
    <div class="mx-auto flex w-full max-w-4xl flex-col gap-6 px-6 py-10 pb-24">
      <RouterView />
    </div>
  </div>

  <BottomTabBar :items="tabItems" container-class="max-w-4xl" />
</template>
