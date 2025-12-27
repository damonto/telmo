<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import ModemMessageBubble from '@/components/modem/messages/ModemMessageBubble.vue'
import ModemMessageThreadSkeletonList from '@/components/modem/messages/ModemMessageThreadSkeletonList.vue'
import type { ThreadMessageItem } from '@/composables/useModemMessageThread'

const props = defineProps<{
  items: ThreadMessageItem[]
  isLoading: boolean
  participantLabel: string
}>()

const { t } = useI18n()

const hasItems = computed(() => props.items.length > 0)
const emptyLabel = computed(() =>
  t('modemDetail.messages.threadEmpty', { participant: props.participantLabel }),
)
</script>

<template>
  <ModemMessageThreadSkeletonList v-if="props.isLoading" />

  <div v-else-if="hasItems" class="space-y-3">
    <ModemMessageBubble v-for="item in props.items" :key="item.key" :item="item" />
  </div>

  <p v-else class="text-sm text-muted-foreground">
    {{ emptyLabel }}
  </p>
</template>
