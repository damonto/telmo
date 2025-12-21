<script setup lang="ts">
import { computed } from 'vue'
import { RouterView } from 'vue-router'
import { House, MessageSquare } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import BottomTabBar from '@/components/BottomTabBar.vue'

const { t } = useI18n()

const navItems = computed(() => [
  {
    key: 'modems',
    routeName: 'home',
    to: { name: 'home' },
    label: t('nav.modems'),
    icon: House,
  },
  {
    key: 'messages',
    routeName: 'messages',
    to: { name: 'messages' },
    label: t('nav.messages'),
    icon: MessageSquare,
  },
])
</script>

<template>
  <div class="min-h-screen bg-linear-to-b from-background via-background to-muted/40">
    <div class="mx-auto flex w-full max-w-6xl flex-col gap-8 px-6 py-10 pb-24">
      <RouterView />
    </div>
    <BottomTabBar :items="navItems" />
  </div>
</template>
