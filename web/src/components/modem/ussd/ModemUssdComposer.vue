<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'

const draft = defineModel<string>({ required: true })

const props = defineProps<{
  isSending: boolean
}>()

const emit = defineEmits<{
  (event: 'submit'): void
}>()

const { t } = useI18n()

const hasDraft = computed(() => draft.value.trim().length > 0)
const isSendDisabled = computed(() => props.isSending || !hasDraft.value)
</script>

<template>
  <form
    class="mt-auto shrink-0 space-y-2 border-t border-border pt-3 pb-3"
    @submit.prevent="emit('submit')"
  >
    <div class="flex items-stretch gap-2">
      <Textarea
        v-model="draft"
        rows="1"
        class="min-h-10 flex-1 resize-none"
        :placeholder="t('modemDetail.ussd.placeholder')"
        :disabled="props.isSending"
      />
      <Button class="h-10" type="submit" :disabled="isSendDisabled">
        <span v-if="props.isSending" class="inline-flex items-center gap-2">
          <span
            class="size-4 animate-spin rounded-full border-2 border-background/60 border-t-background"
          />
          {{ t('modemDetail.ussd.send') }}
        </span>
        <span v-else>{{ t('modemDetail.ussd.send') }}</span>
      </Button>
    </div>
  </form>
</template>
