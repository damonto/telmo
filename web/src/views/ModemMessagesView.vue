<script setup lang="ts">
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import ModemDetailHeader from '@/components/modem/ModemDetailHeader.vue'
import { Card, CardContent } from '@/components/ui/card'
import { useModemDetail } from '@/composables/useModemDetail'

const route = useRoute()
const { t } = useI18n()

const modemId = computed(() => (route.params.id ?? 'unknown') as string)
const { modem, isLoading, fetchModemDetail } = useModemDetail()

watch(
  modemId,
  async (id) => {
    if (!id || id === 'unknown') return
    await fetchModemDetail(id)
  },
  { immediate: true },
)
</script>

<template>
  <div class="space-y-6">
    <ModemDetailHeader :modem="modem" :is-loading="isLoading" />

    <Card
      class="gap-0 rounded-2xl border-white/40 bg-white/80 py-0 shadow-[0_10px_30px_rgba(15,23,42,0.08)] backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/60"
    >
      <CardContent class="space-y-2 px-4 py-4">
        <h2 class="text-lg font-semibold text-foreground">
          {{ t('modemDetail.messages.title') }}
        </h2>
        <p class="text-sm text-muted-foreground">
          {{ t('modemDetail.messages.description') }}
        </p>
      </CardContent>
    </Card>
  </div>
</template>
