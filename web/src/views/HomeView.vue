<script setup lang="ts">
import { RefreshCw } from 'lucide-vue-next'
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { useModemDisplay } from '@/composables/useModemDisplay'
import { useModems } from '@/composables/useModems'

const { t } = useI18n()
const { flagClass, formatSignal, signalIcon, signalTone } = useModemDisplay()

const { modems, isLoading, fetchModems } = useModems()

const modemCount = computed(() => modems.value.length)
const hasModems = computed(() => modems.value.length > 0)
const subtitle = computed(() => t('home.subtitle', { count: modemCount.value }))

/**
 * Determine technology type from access technology string
 */
const getTech = (accessTechnology: string | null): '4G' | '5G' => {
  if (!accessTechnology) return '4G'

  const upperTech = accessTechnology.toUpperCase()
  if (upperTech.includes('5G') || upperTech.includes('NR') || upperTech.includes('SA')) {
    return '5G'
  }

  return '4G'
}

const techTone = (tech: '4G' | '5G') => {
  if (tech === '5G') {
    return 'bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900'
  }
  return 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-200'
}

const handleRefresh = () => {
  fetchModems()
}

// Fetch modems on component mount
onMounted(() => {
  fetchModems()
})
</script>

<template>
  <header class="flex flex-col gap-3">
    <div class="flex items-start justify-between gap-4">
      <div class="flex flex-col gap-3">
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
      </div>
      <button
        type="button"
        class="shrink-0 rounded-lg bg-white/80 p-2.5 shadow-sm ring-1 ring-border/60 transition hover:bg-white dark:bg-slate-950/60 dark:hover:bg-slate-900"
        :disabled="isLoading"
        :title="t('home.refresh')"
        @click="handleRefresh"
      >
        <RefreshCw class="size-5" :class="{ 'animate-spin': isLoading }" />
      </button>
    </div>
  </header>

  <!-- Loading State -->
  <div v-if="isLoading" class="grid gap-3 md:grid-cols-2">
    <Card
      v-for="i in 4"
      :key="`skeleton-${i}`"
      class="h-full gap-0 rounded-2xl border-white/40 bg-white/80 py-0 shadow-[0_10px_30px_rgba(15,23,42,0.08)] backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/60"
    >
      <CardContent class="flex items-center gap-3 px-3 py-2.5">
        <div class="h-11 w-11 shrink-0 animate-pulse rounded-2xl bg-muted/80" />
        <div class="flex flex-1 flex-col gap-1.5">
          <div class="h-4 w-32 animate-pulse rounded bg-muted/60" />
          <div class="h-3.5 w-48 animate-pulse rounded bg-muted/40" />
        </div>
      </CardContent>
    </Card>
  </div>

  <!-- Modems List -->
  <div v-else-if="hasModems" class="grid gap-3 md:grid-cols-2">
    <RouterLink
      v-for="modem in modems"
      :key="modem.id"
      :to="{ name: 'modem-detail', params: { id: modem.id } }"
      class="group block"
    >
      <Card
        class="h-full gap-0 rounded-2xl border-white/40 bg-white/80 py-0 shadow-[0_10px_30px_rgba(15,23,42,0.08)] backdrop-blur-xl transition duration-300 group-hover:-translate-y-0.5 dark:border-white/10 dark:bg-slate-950/60"
      >
        <CardContent class="flex items-center gap-2 px-2.5 py-2">
          <div
            class="flex size-9 shrink-0 items-center justify-center rounded-xl bg-white/80 ring-1 ring-border/60 dark:bg-slate-950/60"
          >
            <span
              v-if="flagClass(modem.sim.regionCode)"
              :class="flagClass(modem.sim.regionCode)"
              class="rounded-sm text-[16px]"
              :aria-label="modem.sim.regionCode"
              :title="modem.sim.regionCode"
            />
            <span
              v-else
              class="rounded-sm text-[10px] font-semibold text-muted-foreground"
              :aria-label="modem.sim.regionCode"
              :title="modem.sim.regionCode"
            >
              {{ modem.sim.regionCode }}
            </span>
          </div>
          <div class="flex min-w-0 flex-1 flex-col gap-1">
            <div class="flex items-center justify-between gap-2">
              <p class="truncate text-sm font-semibold text-foreground">
                {{ modem.name }}
              </p>
              <div class="flex items-center gap-1.5">
                <span
                  class="rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-[0.16em]"
                  :class="techTone(getTech(modem.accessTechnology))"
                >
                  {{ getTech(modem.accessTechnology) }}
                </span>
                <Badge variant="outline" class="text-[10px] tracking-[0.2em]">
                  {{ modem.isEsim ? t('labels.esim') : t('labels.psim') }}
                </Badge>
              </div>
            </div>
            <p class="truncate text-sm text-foreground">
              <span class="text-foreground">{{ modem.sim.operatorName }}</span>
              <span
                v-if="modem.registeredOperator.code && modem.registrationState === 'Roaming'"
                class="text-[11px] text-muted-foreground"
              >
                ({{ modem.registeredOperator.name || modem.registeredOperator.code }})
              </span>
            </p>
            <div class="mt-auto flex items-center justify-between gap-3">
              <p class="truncate text-xs text-muted-foreground">
                {{ modem.number || t('home.noNumber') }}
              </p>
              <div class="flex items-center gap-1">
                <component
                  :is="signalIcon(modem.signalQuality)"
                  class="size-4 shrink-0"
                  :class="signalTone(modem.signalQuality)"
                  :title="`${t('labels.signal')}: ${formatSignal(modem.signalQuality)}`"
                />
                <span
                  v-if="modem.registrationState === 'Roaming'"
                  class="flex size-4 shrink-0 items-center justify-center rounded-full bg-amber-100 text-[9px] font-semibold text-amber-700 dark:bg-amber-500/20 dark:text-amber-300"
                  :aria-label="t('labels.roaming')"
                  :title="t('labels.roaming')"
                >
                  R
                </span>
                <span
                  class="inline-flex h-4 min-w-6 items-center justify-end text-right font-mono text-[10px] text-muted-foreground tabular-nums"
                >
                  {{ formatSignal(modem.signalQuality) }}
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
