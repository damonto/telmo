<script setup lang="ts">
import { Info } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { Button } from '@/components/ui/button'

import type { Modem } from '@/types/modem'

const props = withDefaults(
  defineProps<{
    modem: Modem | null
    isLoading: boolean
    showDetailsAction?: boolean
  }>(),
  {
    showDetailsAction: false,
  },
)

const emit = defineEmits<{
  (event: 'open-details'): void
}>()

const { t } = useI18n()
const router = useRouter()

const handleBack = () => {
  router.back()
}
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-3">
      <Button
        variant="ghost"
        size="sm"
        type="button"
        class="px-0 text-muted-foreground"
        @click="handleBack"
      >
        ← {{ t('modemDetail.back') }}
      </Button>
      <Button
        v-if="props.showDetailsAction"
        variant="ghost"
        size="icon"
        type="button"
        :aria-label="t('modemDetail.tabs.detail')"
        :title="t('modemDetail.tabs.detail')"
        @click="emit('open-details')"
      >
        <Info class="size-4 text-muted-foreground" />
      </Button>
    </div>

    <header class="space-y-2">
      <div class="flex flex-wrap items-center gap-3">
        <h1
          v-if="!props.isLoading"
          class="text-3xl font-semibold tracking-tight text-foreground"
        >
          {{ props.modem?.name ?? t('modemDetail.unknown') }}
        </h1>
        <div v-else class="h-9 w-48 animate-pulse rounded bg-muted" />
      </div>
      <p class="text-sm text-muted-foreground">
        {{ t('modemDetail.subtitle') }}
      </p>
    </header>
  </div>
</template>
