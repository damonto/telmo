<script setup lang="ts">
import type { SlotInfo } from '@/types/modem'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

interface Props {
  slots: SlotInfo[]
  modelValue: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const { t } = useI18n()

const pendingIdentifier = ref<string | null>(null)
const dialogOpen = ref(false)

const hasMultipleSlots = computed(() => props.slots.length > 1)

const openDialog = (identifier: string) => {
  if (identifier === props.modelValue) return
  pendingIdentifier.value = identifier
  dialogOpen.value = true
}

const closeDialog = () => {
  pendingIdentifier.value = null
  dialogOpen.value = false
}

const confirmSwitch = () => {
  if (!pendingIdentifier.value) return
  // TODO: Call API to switch SIM slot
  emit('update:modelValue', pendingIdentifier.value)
  closeDialog()
}

const getSlotLabel = (slot: SlotInfo, index: number) => {
  return slot.active ? 'Active' : `SIM ${index + 1}`
}

const getPendingSlotInfo = computed(() => {
  if (!pendingIdentifier.value) return null
  return props.slots.find((slot) => slot.identifier === pendingIdentifier.value)
})
</script>

<template>
  <div v-if="hasMultipleSlots && slots.length > 0">
    <!-- SIM Slot Switcher -->
    <div
      class="inline-flex w-fit self-start rounded-full border border-white/50 bg-white/80 p-0.5 shadow-sm"
    >
      <button
        v-for="(slot, index) in slots"
        :key="slot.identifier"
        type="button"
        class="rounded-full px-2 py-1 text-[10px] font-semibold tracking-[0.16em] transition"
        :class="
          modelValue === slot.identifier
            ? 'bg-white text-foreground shadow-sm'
            : 'text-muted-foreground hover:text-foreground'
        "
        @click="openDialog(slot.identifier)"
      >
        {{ getSlotLabel(slot, index) }}
      </button>
    </div>

    <!-- Confirmation Dialog -->
    <div
      v-if="dialogOpen"
      class="fixed inset-0 z-30 flex items-center justify-center bg-black/30 px-4 backdrop-blur-sm"
    >
      <div class="w-full max-w-sm rounded-2xl bg-white p-6 shadow-2xl">
        <h3 class="text-sm font-semibold text-foreground">
          Switch to {{ getPendingSlotInfo?.operatorName || 'SIM' }}?
        </h3>
        <div class="mt-6 flex gap-3">
          <button
            type="button"
            class="flex-1 rounded-full border border-border px-4 py-2 text-sm font-semibold text-muted-foreground"
            @click="closeDialog"
          >
            {{ t('modemDetail.actions.cancel') }}
          </button>
          <button
            type="button"
            class="flex-1 rounded-full bg-foreground px-4 py-2 text-sm font-semibold text-background"
            @click="confirmSwitch"
          >
            {{ t('modemDetail.actions.confirm') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
