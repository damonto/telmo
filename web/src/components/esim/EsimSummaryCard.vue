<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useModemDisplay } from '@/composables/useModemDisplay'
import type { EuiccApiResponse } from '@/types/euicc'
import type { Modem } from '@/types/modem'

const props = defineProps<{
  modem: Modem
  euicc: EuiccApiResponse | null
}>()

const { t } = useI18n()
const { formatSignal, signalIcon, signalTone } = useModemDisplay()

// Format bytes to human-readable size
const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${Math.round((bytes / Math.pow(k, i)) * 100) / 100} ${sizes[i]}`
}

const storageRemaining = computed(() => {
  return props.euicc ? formatBytes(props.euicc.freeSpace) : 'N/A'
})

const eid = computed(() => {
  return props.euicc?.eid || 'N/A'
})
</script>

<template>
  <details class="group rounded-2xl border border-border bg-card">
    <summary
      class="cursor-pointer list-none px-5 py-4 text-sm transition hover:bg-muted/40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
    >
      <div class="grid gap-3">
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.imei') }}
          </span>
          <span class="font-mono text-sm text-foreground">{{ props.modem.id }}</span>
        </div>
        <div class="flex items-start justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.eid') }}
          </span>
          <span class="break-all font-mono text-sm text-foreground">{{ eid }}</span>
        </div>
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.storageRemaining') }}
          </span>
          <span class="text-sm font-semibold text-emerald-600 dark:text-emerald-400">
            {{ storageRemaining }}
          </span>
        </div>
      </div>
    </summary>
    <div class="border-t border-border px-5 py-4 text-sm">
      <div class="grid gap-3">
        <div v-if="euicc?.sasUp" class="flex items-center justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.sasAccreditation') }}
          </span>
          <span class="text-right text-foreground">{{ euicc.sasUp }}</span>
        </div>
        <div
          v-if="euicc?.certificates && euicc.certificates.length > 0"
          class="flex flex-col gap-2"
        >
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.certificates') }}
          </span>
          <div class="flex flex-col gap-1">
            <span
              v-for="(cert, index) in euicc.certificates"
              :key="index"
              class="rounded-md bg-muted px-2 py-1 text-xs text-foreground"
            >
              {{ cert }}
            </span>
          </div>
        </div>
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.manufacturer') }}
          </span>
          <span class="text-foreground">{{ props.modem.manufacturer }}</span>
        </div>
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.carrier') }}
          </span>
          <span class="text-foreground">{{ props.modem.sim.operatorName }}</span>
        </div>
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.iccid') }}
          </span>
          <span class="font-mono text-foreground">{{ props.modem.sim.identifier }}</span>
        </div>
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.active') }}
          </span>
          <span class="text-foreground">{{ props.modem.sim.active ? 'Yes' : 'No' }}</span>
        </div>
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.roamingCarrier') }}
          </span>
          <span class="text-muted-foreground">
            {{ props.modem.registeredOperator.name || '—' }}
          </span>
        </div>
        <div class="flex items-center justify-between gap-4">
          <span class="text-xs font-medium text-muted-foreground">
            {{ t('modemDetail.fields.signal') }}
          </span>
          <div class="flex items-center gap-2">
            <component
              :is="signalIcon(props.modem.signalQuality)"
              class="size-5"
              :class="signalTone(props.modem.signalQuality)"
            />
            <span
              v-if="props.modem.registrationState === 'Roaming'"
              class="flex size-5 items-center justify-center rounded-full bg-amber-100 text-[10px] font-semibold text-amber-700 dark:bg-amber-500/20 dark:text-amber-300"
            >
              R
            </span>
            <span class="font-mono text-xs text-muted-foreground">
              {{ formatSignal(props.modem.signalQuality) }}
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
