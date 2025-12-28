<script setup lang="ts">
import { Download } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import { useModemApi } from '@/apis/modem'
import EsimDownloadConfirmationModal from '@/components/esim/EsimDownloadConfirmationModal.vue'
import EsimDownloadPreviewModal from '@/components/esim/EsimDownloadPreviewModal.vue'
import EsimDownloadProgressModal from '@/components/esim/EsimDownloadProgressModal.vue'
import EsimDownloadResultModal from '@/components/esim/EsimDownloadResultModal.vue'
import EsimInstallDialog from '@/components/esim/EsimInstallDialog.vue'
import EsimProfileSection from '@/components/esim/EsimProfileSection.vue'
import EsimSummaryCard from '@/components/esim/EsimSummaryCard.vue'
import SuccessBanner from '@/components/feedback/SuccessBanner.vue'
import ModemDetailCard from '@/components/modem/ModemDetailCard.vue'
import ModemDetailHeader from '@/components/modem/ModemDetailHeader.vue'
import SimSlotSwitcher from '@/components/modem/SimSlotSwitcher.vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
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
const modemApi = useModemApi()

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

const installDialogOpen = ref(false)
const detailDialogOpen = ref(false)
const confirmationCode = ref('')
const feedbackOpen = ref(false)
const feedbackMessage = ref('')
const resultState = ref<'completed' | 'error' | null>(null)
const resultErrorMessage = ref('')
const resultErrorType = ref<'none' | 'failed' | 'disconnected'>('none')
const resultName = ref('')

const {
  downloadState,
  downloadStage,
  progress,
  errorType,
  errorMessage,
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
const resultTone = computed(() => (resultState.value === 'error' ? 'error' : 'success'))
const resultTitle = computed(() => {
  if (resultState.value === 'error') {
    return t('modemDetail.esim.downloadErrorTitle')
  }
  if (resultState.value === 'completed') {
    return t('modemDetail.esim.downloadCompletedTitle')
  }
  return ''
})
const resultMessage = computed(() => {
  if (resultState.value === 'error') {
    if (resultErrorMessage.value) return resultErrorMessage.value
    return resultErrorType.value === 'disconnected'
      ? t('modemDetail.esim.downloadDisconnected')
      : t('modemDetail.esim.downloadErrorFallback')
  }
  if (resultState.value === 'completed') {
    const fallbackName = t('modemDetail.esim.downloadCompletedFallbackName')
    const name = resultName.value || fallbackName
    return t('modemDetail.esim.downloadCompletedMessage', { name })
  }
  return ''
})

const confirmationTitle = computed(() => t('modemDetail.esim.downloadConfirmationTitle'))
const confirmationHint = computed(() => t('modemDetail.esim.downloadConfirmationHint'))
const confirmationPlaceholder = computed(() =>
  t('modemDetail.esim.downloadConfirmationPlaceholder'),
)

const refreshModem = async () => {
  if (!modemId.value || modemId.value === 'unknown') return
  await fetchModemDetail(modemId.value)
}

const showSuccess = (message: string) => {
  feedbackMessage.value = message
  feedbackOpen.value = true
}

const getSimLabel = (identifier: string) => {
  const index = simSlots.value.findIndex((slot) => slot.identifier === identifier)
  if (index === 0) return t('modemDetail.sim.sim1')
  if (index === 1) return t('modemDetail.sim.sim2')
  if (index >= 0) return `SIM ${index + 1}`
  return ''
}

watch(downloadState, (value) => {
  if (value === 'confirmation') {
    confirmationCode.value = ''
  }
  if (value === 'connecting') {
    resultState.value = null
    resultErrorMessage.value = ''
    resultErrorType.value = 'none'
    resultName.value = ''
  }
  if (value === 'error') {
    resultState.value = 'error'
    resultErrorMessage.value = errorMessage.value
    resultErrorType.value = errorType.value
  }
  if (value === 'completed') {
    resultState.value = 'completed'
    resultName.value = downloadedName.value
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

const handleSimSwitch = async (identifier: string) => {
  if (!modemId.value || modemId.value === 'unknown') {
    throw new Error('Modem ID is unavailable')
  }
  await modemApi.switchSimSlot(modemId.value, identifier)
  await refreshModem()
  const simLabel = getSimLabel(identifier)
  if (simLabel) {
    showSuccess(t('modemDetail.sim.switchSuccess', { sim: simLabel }))
  } else {
    showSuccess(t('modemDetail.sim.switchSuccessFallback'))
  }
}
</script>

<template>
  <ModemDetailHeader
    :modem="modem"
    :is-loading="isLoading"
    :show-details-action="isEsimModem"
    @open-details="detailDialogOpen = true"
  />

  <div
    v-if="!modem && !isLoading"
    class="rounded-2xl border border-dashed border-border p-8 text-sm text-muted-foreground"
  >
    {{ t('modemDetail.unknown') }}
  </div>

  <!-- SIM Slot Switcher -->
  <SimSlotSwitcher
    v-if="modem"
    v-model="currentSimIdentifier"
    :slots="simSlots"
    :on-switch="handleSimSwitch"
  />

  <!-- eSIM modem: show original layout -->
  <div v-if="modem && isEsimModem" class="space-y-4">
    <EsimSummaryCard :modem="modem" :euicc="euicc" />
    <EsimProfileSection
      v-model:profiles="esimProfiles"
      :loading="isEsimProfilesLoading"
      :modem-id="modemId"
      :refresh-modem="refreshModem"
      @success="showSuccess"
    />
  </div>

  <!-- Physical modem: show detail card -->
  <div v-if="modem && isPhysicalModem" class="space-y-4">
    <ModemDetailCard :modem="modem" :euicc="null" />
  </div>

  <Dialog v-model:open="detailDialogOpen">
    <DialogContent v-if="modem && isEsimModem" class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>{{ t('modemDetail.tabs.detail') }}</DialogTitle>
      </DialogHeader>
      <div class="max-h-[70vh] overflow-y-auto pr-2">
        <ModemDetailCard :modem="modem" :euicc="euicc" />
      </div>
    </DialogContent>
  </Dialog>

  <Button
    v-if="modem && isEsimModem"
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

  <SuccessBanner v-model:open="feedbackOpen" :message="feedbackMessage" />
</template>
