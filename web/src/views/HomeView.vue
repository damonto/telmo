<script setup lang="ts">
import { RefreshCw } from 'lucide-vue-next'
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
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

const techVariant = (tech: '4G' | '5G') => (tech === '5G' ? 'default' : 'secondary')

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
      <Button
        type="button"
        variant="outline"
        size="icon"
        class="shrink-0"
        :disabled="isLoading"
        :title="t('home.refresh')"
        @click="handleRefresh"
      >
        <RefreshCw class="size-5" :class="{ 'animate-spin': isLoading }" />
      </Button>
    </div>
  </header>

  <!-- Loading State -->
  <div v-if="isLoading" class="grid gap-3 md:grid-cols-2">
    <Card v-for="i in 4" :key="`skeleton-${i}`" class="h-full">
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
      <Card class="h-full transition duration-300 group-hover:-translate-y-0.5 py-2">
        <CardContent class="flex items-center gap-3">
          <div
            class="flex size-9 shrink-0 items-center justify-center rounded-lg border bg-background"
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
                <span class="text-foreground">{{ modem.sim.operatorName }}</span>
                <span
                  v-if="modem.registeredOperator.code && modem.registrationState === 'Roaming'"
                  class="text-[11px] text-muted-foreground"
                >
                  ({{ modem.registeredOperator.name || modem.registeredOperator.code }})
                </span>
              </p>
              <div class="flex items-center gap-1.5">
                <Badge :variant="techVariant(getTech(modem.accessTechnology))">
                  {{ getTech(modem.accessTechnology) }}
                </Badge>
                <Badge variant="outline">
                  {{ modem.supportsEsim ? t('labels.esim') : t('labels.psim') }}
                </Badge>
              </div>
            </div>
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
                <Badge
                  v-if="modem.registrationState === 'Roaming'"
                  variant="secondary"
                  class="h-4 px-1.5 text-[9px] font-semibold"
                  :aria-label="t('labels.roaming')"
                  :title="t('labels.roaming')"
                >
                  R
                </Badge>
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

  <Card v-else>
    <CardContent class="py-10 text-center">
      <p class="text-sm text-muted-foreground">
        {{ t('home.noModems') }}
      </p>
    </CardContent>
  </Card>
</template>
