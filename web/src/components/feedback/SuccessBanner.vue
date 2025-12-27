<script setup lang="ts">
import { onBeforeUnmount, watch } from 'vue'
import { CheckCircle2 } from 'lucide-vue-next'

import { Alert, AlertTitle } from '@/components/ui/alert'

const open = defineModel<boolean>('open', { required: true })

const props = defineProps<{
  message: string
}>()

const AUTO_CLOSE_MS = 2600
let closeTimer: ReturnType<typeof setTimeout> | null = null

const clearTimer = () => {
  if (!closeTimer) return
  clearTimeout(closeTimer)
  closeTimer = null
}

const scheduleClose = () => {
  clearTimer()
  closeTimer = setTimeout(() => {
    open.value = false
  }, AUTO_CLOSE_MS)
}

watch(
  () => [open.value, props.message],
  ([isOpen]) => {
    if (!isOpen) {
      clearTimer()
      return
    }
    scheduleClose()
  },
)

onBeforeUnmount(() => {
  clearTimer()
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="-translate-y-3 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="-translate-y-3 opacity-0"
    >
      <div
        v-if="open"
        class="pointer-events-none fixed inset-x-0 top-3 z-50 flex justify-center px-4 sm:top-6"
      >
        <Alert
          class="pointer-events-auto w-full max-w-lg border-emerald-200 bg-emerald-50/95 text-emerald-950 shadow-lg backdrop-blur"
        >
          <CheckCircle2 class="text-emerald-600" />
          <AlertTitle class="text-sm font-semibold leading-snug">
            {{ props.message }}
          </AlertTitle>
        </Alert>
      </div>
    </Transition>
  </Teleport>
</template>
