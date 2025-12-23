<script setup lang="ts">
import { Download } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import EsimInstallDialog from '@/components/modem/EsimInstallDialog.vue'
import EsimProfileSection from '@/components/modem/EsimProfileSection.vue'
import EsimSummaryCard from '@/components/modem/EsimSummaryCard.vue'
import ModemDetailCard from '@/components/modem/ModemDetailCard.vue'
import SimSlotSwitcher from '@/components/modem/SimSlotSwitcher.vue'
import { useModemDetail } from '@/composables/useModemDetail'
import type { EsimProfile } from '@/types/modem'

const route = useRoute()
const { t } = useI18n()

const modemId = computed(() => (route.params.id ?? 'unknown') as string)
const { modem, euicc, isLoading, isPhysicalModem, isEsimModem, fetchModemDetail } = useModemDetail()

// SIM slot switching logic
const currentSimIdentifier = ref('')

const simSlots = computed(() => modem.value?.slots ?? [])

// Initialize current SIM identifier when modem data loads
watch(
  modem,
  (newModem) => {
    if (newModem && !currentSimIdentifier.value) {
      // Set to the first active slot, or first slot if none is active
      const activeSlot = newModem.slots.find((slot) => slot.active)
      currentSimIdentifier.value = activeSlot?.identifier ?? newModem.slots[0]?.identifier ?? ''
    }
  },
  { immediate: true },
)

// Determine modem type
const physicalModem = computed(() => (isPhysicalModem.value ? modem.value : null))
const esimModem = computed(() => (isEsimModem.value ? modem.value : null))

// TODO: Implement profiles API when available
const esimProfiles = computed<EsimProfile[]>({
  get: () => [],
  set: () => {
    // Placeholder for future implementation
  },
})

const installDialogOpen = ref(false)

// Fetch modem detail when route changes or on mount
const loadModemDetail = async () => {
  await fetchModemDetail(modemId.value)
}

onMounted(() => {
  loadModemDetail()
})

watch(modemId, () => {
  loadModemDetail()
})
</script>

<template>
  <div
    v-if="isLoading"
    class="flex items-center justify-center rounded-2xl border border-dashed border-border p-8 text-sm text-muted-foreground"
  >
    Loading modem details...
  </div>

  <div
    v-else-if="!modem"
    class="rounded-2xl border border-dashed border-border p-8 text-sm text-muted-foreground"
  >
    {{ t('modemDetail.unknown') }}
  </div>

  <!-- SIM Slot Switcher -->
  <SimSlotSwitcher v-if="modem" v-model="currentSimIdentifier" :slots="simSlots" />

  <!-- eSIM modem: show original layout -->
  <div v-if="esimModem" class="space-y-4">
    <EsimSummaryCard :modem="esimModem" :euicc="euicc" />
    <EsimProfileSection v-model:profiles="esimProfiles" />
  </div>

  <!-- Physical modem: show detail card -->
  <div v-if="physicalModem" class="space-y-4">
    <ModemDetailCard :modem="physicalModem" />
  </div>

  <button
    v-if="esimModem"
    type="button"
    class="fixed bottom-24 right-6 z-20 flex size-12 items-center justify-center rounded-full bg-foreground text-background shadow-xl transition hover:-translate-y-0.5"
    @click="installDialogOpen = true"
    :aria-label="t('modemDetail.esim.installButton')"
    :title="t('modemDetail.esim.installButton')"
  >
    <Download class="size-5" />
  </button>

  <EsimInstallDialog v-model:open="installDialogOpen" />
</template>
