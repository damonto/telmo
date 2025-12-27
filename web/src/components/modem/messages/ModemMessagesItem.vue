<script setup lang="ts">
import { computed } from 'vue'
import { Trash2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'

import { Button } from '@/components/ui/button'
import type { ConversationItem } from '@/composables/useModemMessages'

const props = defineProps<{
  item: ConversationItem
  modemId: string
}>()

const emit = defineEmits<{
  (event: 'delete', item: ConversationItem): void
}>()

const { t } = useI18n()

const threadRoute = computed(() => ({
  name: 'modem-message-thread',
  params: { id: props.modemId, participant: props.item.participantValue },
}))
</script>

<template>
  <div class="rounded-lg bg-card px-4 py-3 shadow-sm">
    <div class="flex flex-col gap-1">
      <div class="flex items-center justify-between gap-3">
        <RouterLink :to="threadRoute" class="min-w-0 flex-1">
          <p class="truncate text-sm font-semibold text-foreground">
            {{ props.item.participantLabel }}
          </p>
        </RouterLink>
        <span class="shrink-0 text-xs font-medium text-muted-foreground">
          {{ props.item.timestampLabel }}
        </span>
      </div>
      <div class="flex items-center justify-between gap-3">
        <RouterLink :to="threadRoute" class="min-w-0 flex-1">
          <p class="truncate text-xs text-muted-foreground">
            {{ props.item.preview }}
          </p>
        </RouterLink>
        <Button
          variant="ghost"
          size="icon"
          type="button"
          :aria-label="t('modemDetail.actions.delete')"
          @click.stop.prevent="emit('delete', props.item)"
        >
          <Trash2 class="size-4 text-muted-foreground" />
        </Button>
      </div>
    </div>
  </div>
</template>
