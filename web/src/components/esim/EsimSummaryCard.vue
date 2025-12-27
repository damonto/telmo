<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { EuiccApiResponse } from '@/types/euicc'
import type { Modem } from '@/types/modem'

const props = defineProps<{
  modem: Modem
  euicc: EuiccApiResponse | null
}>()

const { t } = useI18n()

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
  <section class="grid gap-3 text-sm">
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
      <span class="text-sm font-semibold dark:text-emerald-400">
        {{ storageRemaining }}
      </span>
    </div>
  </section>
</template>
