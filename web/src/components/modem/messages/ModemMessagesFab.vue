<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'

const props = defineProps<{
  disabled: boolean
}>()

const emit = defineEmits<{
  (event: 'click'): void
}>()

const { t } = useI18n()
</script>

<template>
  <Button
    type="button"
    size="icon-lg"
    class="fixed bottom-24 right-6 z-20 rounded-full shadow-xl transition hover:-translate-y-0.5"
    :aria-label="t('modemDetail.messages.newConversation')"
    :title="t('modemDetail.messages.newConversation')"
    :disabled="props.disabled"
    @click="emit('click')"
  >
    <Plus class="size-5" />
  </Button>
</template>
