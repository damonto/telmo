<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Download } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import EsimDownloadConfirmationModal from '@/components/esim/EsimDownloadConfirmationModal.vue'
import EsimDownloadPreviewModal from '@/components/esim/EsimDownloadPreviewModal.vue'
import EsimDownloadProgressModal from '@/components/esim/EsimDownloadProgressModal.vue'
import EsimDownloadResultModal from '@/components/esim/EsimDownloadResultModal.vue'
import EsimInstallDialog from '@/components/esim/EsimInstallDialog.vue'
import EsimProfileSection from '@/components/esim/EsimProfileSection.vue'
import EsimSummaryCard from '@/components/esim/EsimSummaryCard.vue'
import ModemDetailCard from '@/components/modem/ModemDetailCard.vue'
import ModemDetailHeader from '@/components/modem/ModemDetailHeader.vue'
import SimSlotSwitcher from '@/components/modem/SimSlotSwitcher.vue'
import { Button } from '@/components/ui/button'
import { useEsimDownload } from '@/composables/useEsimDownload'
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
  fetchEsimProfiles,
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
const confirmationCode = ref('')

const {
  downloadState,
  downloadStage,
  progress,
  errorType,
  previewProfile,
  downloadedName,
  startDownload,
  confirmPreview,
  submitConfirmationCode,
  cancelDownload,
  closeDialog,
} = useEsimDownload(modemId, {
  onCompleted: () => {
    if (!modemId.value || modemId.value === 'unknown') return
    void fetchEsimProfiles(modemId.value)
  },
})

const isProgressModalOpen = computed(
  () => downloadState.value === 'connecting' || downloadState.value === 'progress',
)
const isPreviewModalOpen = computed(() => downloadState.value === 'preview')
const isConfirmationModalOpen = computed(() => downloadState.value === 'confirmation')
const isResultModalOpen = computed(
  () => downloadState.value === 'completed' || downloadState.value === 'error',
)

const stageLabel = computed(() => {
  if (downloadStage.value === 'initializing') return t('modemDetail.esim.downloadStageInitializing')
  if (downloadStage.value === 'connecting') return t('modemDetail.esim.downloadStageConnecting')
  if (downloadStage.value === 'installing') return t('modemDetail.esim.downloadStageInstalling')
  return t('modemDetail.esim.downloadConnecting')
})

const progressTitle = computed(() => t('modemDetail.esim.downloadTitle'))
const resultTone = computed(() => (downloadState.value === 'error' ? 'error' : 'success'))
const resultTitle = computed(() =>
  downloadState.value === 'error'
    ? t('modemDetail.esim.downloadErrorTitle')
    : t('modemDetail.esim.downloadCompletedTitle'),
)
const resultMessage = computed(() => {
  if (downloadState.value === 'error') {
    return errorType.value === 'disconnected'
      ? t('modemDetail.esim.downloadDisconnected')
      : t('modemDetail.esim.downloadErrorFallback')
  }
  const fallbackName = t('modemDetail.esim.downloadCompletedFallbackName')
  const name = downloadedName.value || fallbackName
  return t('modemDetail.esim.downloadCompletedMessage', { name })
})

const confirmationTitle = computed(() => t('modemDetail.esim.downloadConfirmationTitle'))
const confirmationHint = computed(() => t('modemDetail.esim.downloadConfirmationHint'))
const confirmationPlaceholder = computed(() => t('modemDetail.esim.downloadConfirmationPlaceholder'))

const refreshModem = async () => {
  if (!modemId.value || modemId.value === 'unknown') return
  await fetchModemDetail(modemId.value)
}

watch(downloadState, (value) => {
  if (value === 'confirmation') {
    confirmationCode.value = ''
  }
})

// Fetch modem detail when route changes or on mount
watch(
  modemId,
  async (id) => {
    if (!id || id === 'unknown') return
    await fetchModemDetail(id)
  },
  { immediate: true },
)

const handleConfirmationSubmit = () => {
  submitConfirmationCode(confirmationCode.value)
}

const handlePreviewConfirm = () => {
  confirmPreview(true)
}

const handlePreviewCancel = () => {
  confirmPreview(false)
}

const handleResultConfirm = () => {
  closeDialog()
}
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
    <EsimProfileSection
      v-model:profiles="esimProfiles"
      :loading="isEsimProfilesLoading"
      :modem-id="modemId"
      :refresh-modem="refreshModem"
    />
  </div>

  <!-- Physical modem: show detail card -->
  <div v-if="physicalModem" class="space-y-4">
    <ModemDetailCard :modem="physicalModem" />
  </div>

  <Button
    v-if="esimModem"
    type="button"
    size="icon-lg"
    class="fixed bottom-24 right-6 z-20 rounded-full shadow-xl transition hover:-translate-y-0.5"
    @click="installDialogOpen = true"
    :aria-label="t('modemDetail.esim.installButton')"
    :title="t('modemDetail.esim.installButton')"
  >
    <Download class="size-5" />
  </Button>

  <EsimInstallDialog v-model:open="installDialogOpen" @confirm="startDownload" />

  <EsimDownloadProgressModal
    :open="isProgressModalOpen"
    :title="progressTitle"
    :stage-label="stageLabel"
    :progress="progress"
    :cancel-label="t('modemDetail.actions.cancel')"
    @cancel="cancelDownload"
  />

  <EsimDownloadPreviewModal
    :open="isPreviewModalOpen"
    :title="t('modemDetail.esim.downloadPreviewTitle')"
    :hint="t('modemDetail.esim.downloadPreviewHint')"
    :profile="previewProfile"
    :confirm-label="t('modemDetail.actions.confirm')"
    :cancel-label="t('modemDetail.actions.cancel')"
    @confirm="handlePreviewConfirm"
    @cancel="handlePreviewCancel"
  />

  <EsimDownloadConfirmationModal
    v-model:code="confirmationCode"
    :open="isConfirmationModalOpen"
    :title="confirmationTitle"
    :hint="confirmationHint"
    :placeholder="confirmationPlaceholder"
    :confirm-label="t('modemDetail.actions.confirm')"
    :cancel-label="t('modemDetail.actions.cancel')"
    @submit="handleConfirmationSubmit"
    @cancel="cancelDownload"
  />

  <EsimDownloadResultModal
    :open="isResultModalOpen"
    :title="resultTitle"
    :message="resultMessage"
    :confirm-label="t('modemDetail.actions.confirm')"
    :tone="resultTone"
    @confirm="handleResultConfirm"
  />
</template>
