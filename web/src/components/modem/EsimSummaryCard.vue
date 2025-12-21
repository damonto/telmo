<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { useModemDisplay } from '@/composables/useModemDisplay'
import type { EsimModem } from '@/types/modem'

const props = defineProps<{
  modem: EsimModem
}>()

const { t } = useI18n()
const { formatSignal, signalIcon, signalTone } = useModemDisplay()
</script>

<template>
  <details
    class="group rounded-2xl border border-white/40 bg-white/80 shadow-[0_10px_30px_rgba(15,23,42,0.08)] backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/60"
  >
    <summary class="cursor-pointer list-none px-4 py-4">
      <div class="space-y-3 text-sm">
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
            {{ t('modemDetail.fields.imei') }}
          </span>
          <span class="font-mono text-sm text-foreground">
            {{ props.modem.imei }}
          </span>
        </div>
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
            {{ t('modemDetail.fields.eid') }}
          </span>
          <span class="font-mono text-sm text-foreground">
            {{ props.modem.eid }}
          </span>
        </div>
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
            {{ t('modemDetail.fields.storageRemaining') }}
          </span>
          <span class="text-sm font-semibold text-foreground">
            {{ props.modem.storageRemaining }}
          </span>
        </div>
      </div>
    </summary>
    <div class="border-t border-white/40 px-4 py-4 text-sm">
      <div class="grid gap-3">
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
      </div>
    </div>
  </details>
</template>

<style scoped>
summary::marker {
  content: '';
}

summary::-webkit-details-marker {
  display: none;
}
</style>
