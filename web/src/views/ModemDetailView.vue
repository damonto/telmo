<script setup lang="ts">
import { Download } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import EsimInstallDialog from '@/components/modem/EsimInstallDialog.vue'
import EsimProfileSection from '@/components/modem/EsimProfileSection.vue'
import EsimSummaryCard from '@/components/modem/EsimSummaryCard.vue'
import ModemDetailHeader from '@/components/modem/ModemDetailHeader.vue'
import ModemDetailCard from '@/components/modem/ModemDetailCard.vue'
import SimSlotSwitcher from '@/components/modem/SimSlotSwitcher.vue'
import { useModemDetail } from '@/composables/useModemDetail'

const route = useRoute()
const { t } = useI18n()

const modemId = computed(() => (route.params.id ?? 'unknown') as string)
const {
  modem,
  euicc,
  esimProfiles,
  isLoading,
  isEsimProfilesLoading,
  isPhysicalModem,
  isEsimModem,
  fetchModemDetail,
} = useModemDetail()

// SIM slot switching logic
const currentSimIdentifier = ref('')

const simSlots = computed(() => modem.value?.slots ?? [])

// Initialize current SIM identifier when modem data loads
watch(
  modem,
  (newModem) => {
    if (!newModem) {
      currentSimIdentifier.value = ''
      return
    }

    // Set to the first active slot, or first slot if none is active
    const activeSlot = newModem.slots.find((slot) => slot.active)
    currentSimIdentifier.value = activeSlot?.identifier ?? newModem.slots[0]?.identifier ?? ''
  },
  { immediate: true },
)

// Determine modem type
const physicalModem = computed(() => (isPhysicalModem.value ? modem.value : null))
const esimModem = computed(() => (isEsimModem.value ? modem.value : null))

const installDialogOpen = ref(false)

// Fetch modem detail when route changes or on mount
watch(
  modemId,
  async (id) => {
    if (!id || id === 'unknown') return
    await fetchModemDetail(id)
  },
  { immediate: true },
)
</script>

<template>
  <ModemDetailHeader :modem="modem" :is-loading="isLoading" />

  <div
    v-if="!modem && !isLoading"
    class="rounded-2xl border border-dashed border-border p-8 text-sm text-muted-foreground"
  >
    {{ t('modemDetail.unknown') }}
  </div>

  <!-- SIM Slot Switcher -->
  <SimSlotSwitcher v-if="modem" v-model="currentSimIdentifier" :slots="simSlots" />

  <!-- eSIM modem: show original layout -->
  <div v-if="esimModem" class="space-y-4">
    <EsimSummaryCard :modem="esimModem" :euicc="euicc" />
    <EsimProfileSection v-model:profiles="esimProfiles" :loading="isEsimProfilesLoading" />
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
