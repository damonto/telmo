<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import type { SlotInfo } from '@/types/modem'

const props = defineProps<{
  slots: SlotInfo[]
}>()

const selectedIdentifier = defineModel<string>({ required: true })

const { t } = useI18n()

const pendingIdentifier = ref<string | null>(null)
const dialogOpen = ref(false)

const hasMultipleSlots = computed(() => props.slots.length > 1)

const openDialog = (identifier: string) => {
  if (identifier === selectedIdentifier.value) return
  pendingIdentifier.value = identifier
  dialogOpen.value = true
}

const handleSelect = (identifier: string) => {
  if (identifier === selectedIdentifier.value) return
  openDialog(identifier)
}

const closeDialog = () => {
  pendingIdentifier.value = null
  dialogOpen.value = false
}

const confirmSwitch = () => {
  if (!pendingIdentifier.value) return
  // TODO: Call API to switch SIM slot
  selectedIdentifier.value = pendingIdentifier.value
  closeDialog()
}

const getSlotLabel = (slot: SlotInfo, index: number) => {
  return slot.active ? 'Active' : `SIM ${index + 1}`
}

const getPendingSlotInfo = computed(() => {
  if (!pendingIdentifier.value) return null
  return props.slots.find((slot) => slot.identifier === pendingIdentifier.value)
})

const confirmTitle = computed(() => {
  const operatorName = getPendingSlotInfo.value?.operatorName || 'SIM'
  return `Switch to ${operatorName}?`
})
</script>

<template>
  <div v-if="hasMultipleSlots && slots.length > 0">
    <!-- SIM Slot Switcher -->
    <RadioGroup
      :model-value="selectedIdentifier"
      class="flex flex-wrap gap-4"
      @update:model-value="handleSelect"
    >
      <div v-for="(slot, index) in slots" :key="slot.identifier" class="flex items-center gap-2">
        <RadioGroupItem :id="`sim-slot-${slot.identifier}`" :value="slot.identifier" />
        <Label
          :for="`sim-slot-${slot.identifier}`"
          class="text-[10px] font-semibold uppercase tracking-[0.16em]"
        >
          {{ getSlotLabel(slot, index) }}
        </Label>
      </div>
    </RadioGroup>

    <!-- Confirmation Dialog -->
    <AlertDialog v-model:open="dialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ confirmTitle }}</AlertDialogTitle>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="closeDialog">
            {{ t('modemDetail.actions.cancel') }}
          </AlertDialogCancel>
          <AlertDialogAction @click="confirmSwitch">
            {{ t('modemDetail.actions.confirm') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
