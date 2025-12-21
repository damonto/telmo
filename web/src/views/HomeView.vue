<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'

import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { useModemDisplay } from '@/composables/useModemDisplay'
import { useModems } from '@/composables/useModems'
import type { Modem } from '@/types/modem'

const { t } = useI18n()
const { flagClass, formatSignal, signalIcon, signalTone } = useModemDisplay()

const { modems } = useModems()

const modemCount = computed(() => modems.value.length)
const hasModems = computed(() => modems.value.length > 0)
const subtitle = computed(() => t('home.subtitle', { count: modemCount.value }))

const techTone = (tech: Modem['tech']) => {
  if (tech === '5G') {
    return 'bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900'
  }
  return 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200'
}
</script>

<template>
  <header class="flex flex-col gap-3">
    <p class="text-xs font-semibold uppercase tracking-[0.3em] text-muted-foreground">
      {{ t('home.kicker') }}
    </p>
    <div class="flex flex-col gap-2">
      <h1 class="text-3xl font-semibold tracking-tight text-foreground md:text-4xl">
        {{ t('home.title') }}
      </h1>
      <p class="text-sm text-muted-foreground md:text-base">
        {{ subtitle }}
      </p>
    </div>
  </header>

  <div v-if="hasModems" class="grid gap-3 md:grid-cols-2">
    <RouterLink
      v-for="modem in modems"
      :key="modem.id"
      :to="{ name: 'modem-detail', params: { id: modem.id } }"
      class="group block"
    >
      <Card
        class="h-full gap-0 rounded-2xl border-white/40 bg-white/80 py-0 shadow-[0_10px_30px_rgba(15,23,42,0.08)] backdrop-blur-xl transition duration-300 group-hover:-translate-y-0.5 dark:border-white/10 dark:bg-slate-950/60"
      >
        <CardContent class="flex items-center gap-3 px-3 py-2.5">
          <div
            class="flex size-11 shrink-0 items-center justify-center rounded-2xl bg-white/80 ring-1 ring-border/60 dark:bg-slate-950/60"
          >
            <img
              v-if="modem.logoUrl"
              :src="modem.logoUrl"
              :alt="`${modem.name} logo`"
              class="size-6 object-contain"
            />
            <span
              v-else-if="flagClass(modem.regionCode)"
              :class="flagClass(modem.regionCode)"
              class="rounded-sm text-[18px]"
              :aria-label="modem.regionCode"
              :title="modem.regionCode"
            />
            <span
              v-else
              class="rounded-sm text-xs font-semibold text-muted-foreground"
              :aria-label="modem.regionCode"
              :title="modem.regionCode"
            >
              {{ modem.regionCode }}
            </span>
          </div>
          <div class="flex min-w-0 flex-1 flex-col gap-1.5">
            <div class="flex items-center justify-between gap-2">
              <p class="truncate text-sm font-semibold text-foreground">
                {{ modem.name }}
              </p>
              <div class="flex items-center gap-1.5">
                <span
                  class="rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-[0.16em]"
                  :class="techTone(modem.tech)"
                >
                  {{ modem.tech }}
                </span>
                <Badge variant="outline" class="text-[10px] tracking-[0.2em]">
                  {{ modem.isESim ? t('labels.esim') : t('labels.psim') }}
                </Badge>
              </div>
            </div>
            <div class="flex items-center justify-between gap-3">
              <div class="min-w-0">
                <p class="truncate text-sm text-foreground">
                  <span class="text-foreground">{{ modem.carrierName }}</span>
                  <span
                    v-if="modem.isRoaming && modem.roamingCarrierName"
                    class="text-[11px] text-muted-foreground"
                  >
                    ({{ modem.roamingCarrierName }})
                  </span>
                </p>
              </div>
              <div class="flex items-center gap-1.5">
                <component
                  :is="signalIcon(modem.signalDbm)"
                  class="size-5 shrink-0 -mt-1"
                  :class="signalTone(modem.signalDbm)"
                  :title="`${t('labels.signal')}: ${formatSignal(modem.signalDbm)}`"
                />
                <span
                  v-if="modem.isRoaming"
                  class="flex size-5 shrink-0 items-center justify-center rounded-full bg-amber-100 text-[10px] font-semibold text-amber-700 dark:bg-amber-500/20 dark:text-amber-300"
                  :aria-label="t('labels.roaming')"
                  :title="t('labels.roaming')"
                >
                  R
                </span>
                <span
                  class="inline-flex h-5 min-w-12 items-center justify-end text-right font-mono text-[11px] text-muted-foreground tabular-nums"
                >
                  {{ formatSignal(modem.signalDbm) }}
                </span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </RouterLink>
  </div>

  <div v-else class="rounded-lg border border-dashed border-border p-10 text-center">
    <p class="text-sm text-muted-foreground">
      {{ t('home.noModems') }}
    </p>
  </div>
</template>
