<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'

const props = defineProps<{
  count: number
  isLoading: boolean
  modemId: string
}>()

const { t } = useI18n()

const badgeLabel = computed(() => (props.isLoading ? '...' : String(props.count)))
const backRoute = computed(() =>
  props.modemId && props.modemId !== 'unknown'
    ? { name: 'modem-detail', params: { id: props.modemId } }
    : { name: 'home' },
)
</script>

<template>
  <header class="space-y-2 pb-3">
    <div class="flex items-center justify-between gap-3">
      <div class="space-y-1">
        <Button as-child variant="ghost" size="sm" class="px-0 text-muted-foreground">
          <RouterLink :to="backRoute">
            &larr; {{ t('modemDetail.back') }}
          </RouterLink>
        </Button>
        <div class="space-y-1">
          <h1 class="text-2xl font-semibold text-foreground">
            {{ t('modemDetail.notifications.title') }}
          </h1>
          <p class="text-sm text-muted-foreground">
            {{ t('modemDetail.notifications.subtitle') }}
          </p>
        </div>
      </div>
      <Badge variant="outline" class="text-[10px] uppercase tracking-[0.2em]">
        {{ badgeLabel }}
      </Badge>
    </div>
  </header>
</template>
