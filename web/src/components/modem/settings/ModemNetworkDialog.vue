<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import ModemNetworkOption from '@/components/modem/settings/ModemNetworkOption.vue'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { RadioGroup } from '@/components/ui/radio-group'
import type { NetworkResponse } from '@/types/network'

const open = defineModel<boolean>('open', { required: true })
const selectedNetwork = defineModel<string>('selectedNetwork', { required: true })

const props = defineProps<{
  networks: NetworkResponse[]
  isLoading: boolean
  isRegistering: boolean
  hasAvailableNetworks: boolean
  hasSelection: boolean
}>()

const emit = defineEmits<{
  (event: 'register'): void
}>()

const { t } = useI18n()
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>{{ t('modemDetail.settings.networkDialogTitle') }}</DialogTitle>
        <DialogDescription>
          {{ t('modemDetail.settings.networkDialogDescription') }}
        </DialogDescription>
      </DialogHeader>

      <div class="max-h-[60vh] overflow-y-auto pr-1">
        <div v-if="props.isLoading" class="flex items-center justify-center py-10">
          <span
            class="size-6 animate-spin rounded-full border-2 border-muted-foreground/40 border-t-muted-foreground"
            aria-hidden="true"
          />
          <span class="sr-only">{{ t('modemDetail.actions.loading') }}</span>
        </div>

        <div v-else-if="props.hasAvailableNetworks" class="space-y-2">
          <RadioGroup v-model="selectedNetwork" class="gap-2">
            <ModemNetworkOption
              v-for="network in props.networks"
              :key="network.operatorCode"
              :network="network"
              :is-selected="selectedNetwork === network.operatorCode"
            />
          </RadioGroup>
        </div>

        <p v-else class="text-sm text-muted-foreground">
          {{ t('modemDetail.settings.networkEmpty') }}
        </p>
      </div>

      <DialogFooter>
        <Button
          variant="ghost"
          type="button"
          :disabled="props.isRegistering"
          @click="open = false"
        >
          {{ t('modemDetail.actions.cancel') }}
        </Button>
        <Button
          type="button"
          :disabled="!props.hasSelection || props.isRegistering"
          @click="emit('register')"
        >
          <span v-if="props.isRegistering" class="inline-flex items-center gap-2">
            <span
              class="size-4 animate-spin rounded-full border-2 border-background/60 border-t-background"
            />
            {{ t('modemDetail.settings.networkRegister') }}
          </span>
          <span v-else>{{ t('modemDetail.settings.networkRegister') }}</span>
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
