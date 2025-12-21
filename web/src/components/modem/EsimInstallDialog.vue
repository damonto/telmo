<script setup lang="ts">
import { ScanQrCode } from 'lucide-vue-next'
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

type InstallPayload = {
  smdp: string
  activationCode: string
  confirmationCode: string
}

const emit = defineEmits<{
  (event: 'confirm', payload: InstallPayload): void
}>()

const open = defineModel<boolean>('open', { required: true })

const { t } = useI18n()

const smdp = ref('')
const activationCode = ref('')
const confirmationCode = ref('')

const resetForm = () => {
  smdp.value = ''
  activationCode.value = ''
  confirmationCode.value = ''
}

const closeDialog = () => {
  open.value = false
  resetForm()
}

const confirmInstall = () => {
  emit('confirm', {
    smdp: smdp.value,
    activationCode: activationCode.value,
    confirmationCode: confirmationCode.value,
  })
  closeDialog()
}

watch(open, (value) => {
  if (!value) {
    resetForm()
  }
})
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-30 flex items-center justify-center bg-black/30 px-4 backdrop-blur-sm"
  >
    <div class="w-full max-w-sm rounded-2xl bg-white p-6 shadow-2xl">
      <div class="flex items-center justify-between gap-4">
        <h3 class="text-sm font-semibold text-foreground">
          {{ t('modemDetail.esim.installTitle') }}
        </h3>
        <button
          type="button"
          class="flex size-9 items-center justify-center rounded-full border border-border text-muted-foreground transition hover:text-foreground"
          :aria-label="t('modemDetail.esim.scan')"
          :title="t('modemDetail.esim.scan')"
        >
          <ScanQrCode class="size-4" />
        </button>
      </div>
      <div class="mt-4 space-y-3">
        <label class="space-y-1 text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          <span>{{ t('modemDetail.esim.smdp') }}</span>
          <input
            v-model="smdp"
            type="text"
            class="w-full rounded-2xl border border-border/60 px-3 py-2 text-sm text-foreground outline-none transition focus:border-ring"
          />
        </label>
        <label class="space-y-1 text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          <span>{{ t('modemDetail.esim.activationCode') }}</span>
          <input
            v-model="activationCode"
            type="text"
            class="w-full rounded-2xl border border-border/60 px-3 py-2 text-sm text-foreground outline-none transition focus:border-ring"
          />
        </label>
        <label class="space-y-1 text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
          <span>{{ t('modemDetail.esim.confirmationCode') }}</span>
          <input
            v-model="confirmationCode"
            type="text"
            class="w-full rounded-2xl border border-border/60 px-3 py-2 text-sm text-foreground outline-none transition focus:border-ring"
          />
        </label>
      </div>
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
          @click="confirmInstall"
        >
          {{ t('modemDetail.esim.installConfirm') }}
        </button>
      </div>
    </div>
  </div>
</template>
