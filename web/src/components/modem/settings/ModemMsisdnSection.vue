<script setup lang="ts">
import { computed } from 'vue'
import { Save } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Spinner } from '@/components/ui/spinner'

const msisdn = defineModel<string>({ required: true })

const props = defineProps<{
  isLoading: boolean
  isUpdating: boolean
  isValid: boolean
}>()

const emit = defineEmits<{
  (event: 'update'): void
}>()

const { t } = useI18n()

const isInputDisabled = computed(() => props.isLoading || props.isUpdating)
const isActionDisabled = computed(() => !props.isValid || props.isUpdating)
</script>

<template>
  <section class="space-y-4 rounded-2xl bg-card p-4 shadow-sm">
    <h2 class="text-base font-semibold text-foreground">
      {{ t('modemDetail.settings.msisdnTitle') }}
    </h2>
    <div class="flex items-stretch gap-2">
      <Input
        v-model="msisdn"
        type="tel"
        inputmode="tel"
        autocomplete="tel"
        class="flex-1"
        :disabled="isInputDisabled"
        :placeholder="t('modemDetail.settings.msisdnPlaceholder')"
      />
      <Button
        size="icon"
        type="button"
        :disabled="isActionDisabled"
      :aria-label="t('modemDetail.actions.update')"
      @click="emit('update')"
    >
        <Spinner v-if="props.isUpdating" class="size-4" />
        <Save v-else class="size-4" />
      </Button>
    </div>
  </section>
</template>
