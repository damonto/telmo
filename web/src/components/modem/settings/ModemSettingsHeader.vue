<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import { Button } from '@/components/ui/button'

const { t } = useI18n()
</script>

<template>
  <header class="space-y-2 pb-3">
    <Button as-child variant="ghost" size="sm" class="px-0 text-muted-foreground">
      <RouterLink to="/">
        &larr; {{ t('modemDetail.back') }}
      </RouterLink>
    </Button>
    <div class="space-y-1">
      <h1 class="text-2xl font-semibold text-foreground">
        {{ t('modemDetail.settings.title') }}
      </h1>
      <p class="text-sm text-muted-foreground">
        {{ t('modemDetail.settings.subtitle') }}
      </p>
    </div>
  </header>
</template>
