<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import ModemNotificationsItem from '@/components/modem/notifications/ModemNotificationsItem.vue'
import ModemNotificationsSkeletonList from '@/components/modem/notifications/ModemNotificationsSkeletonList.vue'
import type { NotificationItem } from '@/composables/useModemNotifications'

const props = defineProps<{
  items: NotificationItem[]
  isLoading: boolean
  resendingSequence: string | null
}>()

const emit = defineEmits<{
  (event: 'resend', item: NotificationItem): void
  (event: 'delete', item: NotificationItem): void
}>()

const { t } = useI18n()

const hasItems = computed(() => props.items.length > 0)
</script>

<template>
  <section class="space-y-3">
    <ModemNotificationsSkeletonList v-if="props.isLoading" />

    <div v-else-if="hasItems" class="space-y-3">
      <ModemNotificationsItem
        v-for="item in props.items"
        :key="item.key"
        :item="item"
        :is-resending="props.resendingSequence === item.sequenceNumber"
        @resend="emit('resend', $event)"
        @delete="emit('delete', $event)"
      />
    </div>

    <div
      v-else
      class="rounded-lg border border-dashed border-border p-4 text-sm text-muted-foreground"
    >
      {{ t('modemDetail.notifications.empty') }}
    </div>
  </section>
</template>
