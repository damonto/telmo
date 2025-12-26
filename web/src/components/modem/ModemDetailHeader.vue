<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import { Button } from '@/components/ui/button'

import type { Modem } from '@/types/modem'

defineProps<{
  modem: Modem | null
  isLoading: boolean
}>()

const { t } = useI18n()
</script>

<template>
  <div class="space-y-4">
    <Button as-child variant="ghost" size="sm" class="px-0 text-muted-foreground">
      <RouterLink to="/">
        ← {{ t('modemDetail.back') }}
      </RouterLink>
    </Button>

    <header class="space-y-2">
      <p class="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
        {{ t('modemDetail.kicker') }}
      </p>
      <div class="flex flex-wrap items-center gap-3">
        <h1 v-if="!isLoading" class="text-3xl font-semibold tracking-tight text-foreground">
          {{ modem?.name ?? t('modemDetail.unknown') }}
        </h1>
        <div v-else class="h-9 w-48 animate-pulse rounded bg-muted" />
      </div>
      <p class="text-sm text-muted-foreground">
        {{ t('modemDetail.subtitle') }}
      </p>
    </header>
  </div>
</template>
