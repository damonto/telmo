<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'

const props = defineProps<{
  operatorLabel: string
  registrationState: string
  accessTechnology: string
  isScanning: boolean
}>()

const emit = defineEmits<{
  (event: 'scan'): void
}>()

const { t } = useI18n()

const isScanDisabled = computed(() => props.isScanning)
</script>

<template>
  <section class="space-y-4 rounded-2xl bg-card p-4 shadow-sm">
    <div class="flex items-center justify-between gap-4">
      <h2 class="text-base font-semibold text-foreground">
        {{ t('modemDetail.settings.networkTitle') }}
      </h2>
      <Button size="sm" type="button" :disabled="isScanDisabled" @click="emit('scan')">
        <span v-if="props.isScanning" class="inline-flex items-center gap-2">
          <Spinner class="size-4" />
          {{ t('modemDetail.settings.networkSearch') }}
        </span>
        <span v-else>{{ t('modemDetail.settings.networkSearch') }}</span>
      </Button>
    </div>
    <div class="space-y-2 text-sm">
      <div class="flex items-center justify-between gap-4">
        <span class="text-muted-foreground">{{ t('modemDetail.settings.networkOperator') }}</span>
        <span class="font-medium text-foreground">
          {{ props.operatorLabel }}
        </span>
      </div>
      <div class="flex items-center justify-between gap-4">
        <span class="text-muted-foreground">{{ t('modemDetail.settings.networkStatus') }}</span>
        <span class="font-medium text-foreground">
          {{ props.registrationState }}
        </span>
      </div>
      <div class="flex items-center justify-between gap-4">
        <span class="text-muted-foreground">{{ t('modemDetail.settings.networkAccess') }}</span>
        <span class="font-medium text-foreground">
          {{ props.accessTechnology }}
        </span>
      </div>
    </div>
  </section>
</template>
