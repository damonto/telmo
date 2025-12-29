<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import ModemDeviceSettingsSection from '@/components/modem/settings/ModemDeviceSettingsSection.vue'
import ModemMsisdnSection from '@/components/modem/settings/ModemMsisdnSection.vue'
import ModemNetworkDialog from '@/components/modem/settings/ModemNetworkDialog.vue'
import ModemNetworkSection from '@/components/modem/settings/ModemNetworkSection.vue'
import SuccessBanner from '@/components/feedback/SuccessBanner.vue'
import ModemSettingsHeader from '@/components/modem/settings/ModemSettingsHeader.vue'
import { useModemSettings } from '@/composables/useModemSettings'

const route = useRoute()

const modemId = computed(() => (route.params.id ?? 'unknown') as string)

const {
  isModemLoading,
  msisdnInput,
  isMsisdnUpdating,
  isMsisdnValid,
  settingsAlias,
  settingsMss,
  settingsCompatible,
  isSettingsLoading,
  isSettingsUpdating,
  isMssValid,
  currentOperatorLabel,
  currentRegistrationState,
  currentAccessTechnology,
  networkDialogOpen,
  availableNetworks,
  selectedNetwork,
  isNetworkLoading,
  isNetworkRegistering,
  hasAvailableNetworks,
  hasNetworkSelection,
  feedbackOpen,
  feedbackMessage,
  openNetworkDialog,
  handleMsisdnUpdate,
  handleNetworkRegister,
  handleSettingsUpdate,
} = useModemSettings(modemId)
</script>

<template>
  <div class="space-y-3">
    <ModemSettingsHeader />

    <ModemMsisdnSection
      v-model="msisdnInput"
      :is-loading="isModemLoading"
      :is-updating="isMsisdnUpdating"
      :is-valid="isMsisdnValid"
      @update="handleMsisdnUpdate"
    />

    <ModemNetworkSection
      :operator-label="currentOperatorLabel"
      :registration-state="currentRegistrationState"
      :access-technology="currentAccessTechnology"
      :is-scanning="isNetworkLoading"
      @scan="openNetworkDialog"
    />

    <ModemDeviceSettingsSection
      v-model:alias="settingsAlias"
      v-model:mss="settingsMss"
      v-model:compatible="settingsCompatible"
      :is-loading="isSettingsLoading"
      :is-updating="isSettingsUpdating"
      :is-valid="isMssValid"
      @update="handleSettingsUpdate"
    />
  </div>

  <ModemNetworkDialog
    v-model:open="networkDialogOpen"
    v-model:selected-network="selectedNetwork"
    :networks="availableNetworks"
    :is-loading="isNetworkLoading"
    :is-registering="isNetworkRegistering"
    :has-available-networks="hasAvailableNetworks"
    :has-selection="hasNetworkSelection"
    @register="handleNetworkRegister"
  />

  <SuccessBanner v-model:open="feedbackOpen" :message="feedbackMessage" />
</template>
