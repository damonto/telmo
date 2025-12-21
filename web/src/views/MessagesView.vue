<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { Card, CardContent } from '@/components/ui/card'

const { t } = useI18n()
</script>

<template>
  <header class="flex flex-col gap-3">
    <p class="text-xs font-semibold uppercase tracking-[0.3em] text-muted-foreground">
      {{ t('messages.kicker') }}
    </p>
    <div class="flex flex-col gap-2">
      <h1 class="text-3xl font-semibold tracking-tight text-foreground md:text-4xl">
        {{ t('messages.title') }}
      </h1>
      <p class="text-sm text-muted-foreground md:text-base">
        {{ t('messages.subtitle') }}
      </p>
    </div>
  </header>

  <Card
    class="gap-0 rounded-2xl border-white/40 bg-white/80 py-0 shadow-[0_10px_30px_rgba(15,23,42,0.08)] backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/60"
  >
    <CardContent class="space-y-3 px-4 py-4">
      <div>
        <p class="text-lg font-semibold text-foreground">
          {{ t('messages.emptyTitle') }}
        </p>
        <p class="text-sm text-muted-foreground">
          {{ t('messages.emptyDescription') }}
        </p>
      </div>
      <div class="rounded-lg border border-dashed border-border p-6 text-sm text-muted-foreground">
        {{ t('messages.emptyPlaceholder') }}
      </div>
    </CardContent>
  </Card>
</template>
