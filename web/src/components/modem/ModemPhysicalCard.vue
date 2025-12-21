<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { Card, CardContent } from '@/components/ui/card'
import { useModemDisplay } from '@/composables/useModemDisplay'
import type { PhysicalModem } from '@/types/modem'

const props = defineProps<{
  modem: PhysicalModem
}>()

const { t } = useI18n()
const { flagClass, formatSignal, signalIcon, signalTone } = useModemDisplay()
</script>

<template>
  <Card
    class="gap-0 rounded-2xl border-white/40 bg-white/80 py-0 shadow-[0_10px_30px_rgba(15,23,42,0.08)] backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/60"
  >
    <CardContent class="space-y-4 px-4 py-4 text-sm">
      <div class="flex items-center justify-between gap-4">
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          {{ t('modemDetail.fields.moduleName') }}
        </span>
        <span class="font-semibold text-foreground">{{ props.modem.name }}</span>
      </div>
      <div class="flex items-center justify-between gap-4">
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          {{ t('modemDetail.fields.manufacturer') }}
        </span>
        <span class="text-foreground">{{ props.modem.manufacturer }}</span>
      </div>
      <div class="flex items-center justify-between gap-4">
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          {{ t('modemDetail.fields.carrier') }}
        </span>
        <span class="text-foreground">{{ props.modem.carrierName }}</span>
      </div>
      <div class="flex items-center justify-between gap-4">
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          {{ t('modemDetail.fields.roamingCarrier') }}
        </span>
        <span class="text-muted-foreground">
          {{ props.modem.roamingCarrierName ?? '—' }}
        </span>
      </div>
      <div class="flex items-center justify-between gap-4">
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          {{ t('modemDetail.fields.signal') }}
        </span>
        <div class="flex items-center gap-2">
          <component
            :is="signalIcon(props.modem.signalDbm)"
            class="size-5"
            :class="signalTone(props.modem.signalDbm)"
          />
          <span
            v-if="props.modem.isRoaming"
            class="flex size-5 items-center justify-center rounded-full bg-amber-100 text-[10px] font-semibold text-amber-700 dark:bg-amber-500/20 dark:text-amber-300"
          >
            R
          </span>
          <span class="font-mono text-xs text-muted-foreground">
            {{ formatSignal(props.modem.signalDbm) }}
          </span>
        </div>
      </div>
      <div class="flex items-center justify-between gap-4">
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          {{ t('modemDetail.fields.flag') }}
        </span>
        <span class="rounded-sm text-[18px]">
          <span
            v-if="flagClass(props.modem.regionCode)"
            :class="flagClass(props.modem.regionCode)"
          />
          <span v-else class="text-xs font-semibold text-muted-foreground">
            {{ props.modem.regionCode }}
          </span>
        </span>
      </div>
      <div class="flex items-center justify-between gap-4">
        <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          {{ t('modemDetail.fields.iccid') }}
        </span>
        <span class="font-mono text-xs text-foreground">
          {{ props.modem.iccid }}
        </span>
      </div>
    </CardContent>
  </Card>
</template>
