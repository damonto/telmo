<script setup lang="ts">
import { Download } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import EsimInstallDialog from '@/components/modem/EsimInstallDialog.vue'
import EsimProfileSection from '@/components/modem/EsimProfileSection.vue'
import EsimSummaryCard from '@/components/modem/EsimSummaryCard.vue'
import ModemPhysicalCard from '@/components/modem/ModemPhysicalCard.vue'
import { useModems } from '@/composables/useModems'
import type { EsimProfile } from '@/types/modem'

const route = useRoute()
const { t } = useI18n()

const modemId = computed(() => (route.params.id ?? 'unknown') as string)
const { getModemById } = useModems()

const modem = computed(() => getModemById(modemId.value))

const physicalModem = computed(() => (modem.value && !modem.value.isESim ? modem.value : null))
const esimModem = computed(() => (modem.value && modem.value.isESim ? modem.value : null))

const esimProfiles = computed<EsimProfile[]>({
  get: () => esimModem.value?.profiles ?? [],
  set: (value) => {
    if (!esimModem.value) return
    esimModem.value.profiles = value
  },
})

const installDialogOpen = ref(false)
</script>

<template>
  <div
    v-if="!modem"
    class="rounded-2xl border border-dashed border-border p-8 text-sm text-muted-foreground"
  >
    {{ t('modemDetail.unknown') }}
  </div>

  <ModemPhysicalCard v-else-if="physicalModem" :modem="physicalModem" />

  <div v-else class="space-y-6">
    <EsimSummaryCard v-if="esimModem" :modem="esimModem" />
    <EsimProfileSection v-if="esimModem" v-model:profiles="esimProfiles" />
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
