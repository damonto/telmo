<script setup lang="ts">
import { computed, ref } from 'vue'
import { Info, MessageSquare, Phone } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { RouterLink, RouterView, useRoute } from 'vue-router'

import BottomTabBar from '@/components/BottomTabBar.vue'
import { useModems } from '@/composables/useModems'

const route = useRoute()
const { t } = useI18n()

const modemId = computed(() => (route.params.id ?? 'unknown') as string)
const { getModemById } = useModems()

const modem = computed(() => getModemById(modemId.value))

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

const currentSim = ref(1)
const pendingSim = ref<number | null>(null)
const simDialogOpen = ref(false)

const openSimDialog = (sim: number) => {
  if (sim === currentSim.value) return
  pendingSim.value = sim
  simDialogOpen.value = true
}

const closeSimDialog = () => {
  pendingSim.value = null
  simDialogOpen.value = false
}

const confirmSimSwitch = () => {
  if (!pendingSim.value) return
  currentSim.value = pendingSim.value
  closeSimDialog()
}
</script>

<template>
  <div class="min-h-screen bg-linear-to-b from-background via-background to-muted/40">
    <div class="mx-auto flex w-full max-w-4xl flex-col gap-6 px-6 py-10 pb-24">
      <RouterLink
        to="/"
        class="text-sm font-medium text-muted-foreground transition hover:text-foreground"
      >
        ← {{ t('modemDetail.back') }}
      </RouterLink>

      <header class="space-y-2">
        <p class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          {{ t('modemDetail.kicker') }}
        </p>
        <div class="flex flex-wrap items-center gap-3">
          <h1 class="text-3xl font-semibold tracking-tight text-foreground">
            {{ modem?.name ?? t('modemDetail.unknown') }}
          </h1>
        </div>
        <p class="text-sm text-muted-foreground">
          {{ t('modemDetail.subtitle') }}
        </p>
      </header>

      <div
        class="inline-flex w-fit self-start rounded-full border border-white/50 bg-white/80 p-0.5 shadow-sm"
      >
        <button
          type="button"
          class="rounded-full px-2 py-1 text-[10px] font-semibold tracking-[0.16em] transition"
          :class="
            currentSim === 1
              ? 'bg-white text-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'
          "
          @click="openSimDialog(1)"
        >
          {{ t('modemDetail.sim.sim1') }}
        </button>
        <button
          type="button"
          class="rounded-full px-2 py-1 text-[10px] font-semibold tracking-[0.16em] transition"
          :class="
            currentSim === 2
              ? 'bg-white text-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground'
          "
          @click="openSimDialog(2)"
        >
          {{ t('modemDetail.sim.sim2') }}
        </button>
      </div>

      <RouterView />
    </div>
  </div>

  <BottomTabBar :items="tabItems" container-class="max-w-4xl" />

  <div
    v-if="simDialogOpen"
    class="fixed inset-0 z-30 flex items-center justify-center bg-black/30 px-4 backdrop-blur-sm"
  >
    <div class="w-full max-w-sm rounded-2xl bg-white p-6 shadow-2xl">
      <h3 class="text-sm font-semibold text-foreground">
        {{ t('modemDetail.sim.confirm', { sim: pendingSim ?? currentSim }) }}
      </h3>
      <div class="mt-6 flex gap-3">
        <button
          type="button"
          class="flex-1 rounded-full border border-border px-4 py-2 text-sm font-semibold text-muted-foreground"
          @click="closeSimDialog"
        >
          {{ t('modemDetail.actions.cancel') }}
        </button>
        <button
          type="button"
          class="flex-1 rounded-full bg-foreground px-4 py-2 text-sm font-semibold text-background"
          @click="confirmSimSwitch"
        >
          {{ t('modemDetail.actions.confirm') }}
        </button>
      </div>
    </div>
  </div>
</template>
