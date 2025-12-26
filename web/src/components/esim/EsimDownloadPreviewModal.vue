<script setup lang="ts">
import { computed } from 'vue'

import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useModemDisplay } from '@/composables/useModemDisplay'

type DownloadProfilePreview = {
  iccid: string
  serviceProviderName: string
  profileName: string
  profileNickname?: string
  profileState: string
  icon?: string
  regionCode?: string
}

const props = defineProps<{
  open: boolean
  title: string
  hint: string
  profile: DownloadProfilePreview | null
  confirmLabel: string
  cancelLabel: string
}>()

const emit = defineEmits<{
  (event: 'confirm'): void
  (event: 'cancel'): void
}>()

const { flagClass } = useModemDisplay()

const profileName = computed(() => {
  return props.profile?.profileName || props.profile?.serviceProviderName || ''
})

const profileSubtitle = computed(() => props.profile?.iccid ?? '')

const logoUrl = computed(() => props.profile?.icon ?? '')
const regionCode = computed(() => props.profile?.regionCode ?? '')
const regionFlagClass = computed(() => flagClass(regionCode.value))

const handleOpenChange = (nextOpen: boolean) => {
  if (!nextOpen) emit('cancel')
}
</script>

<template>
  <Dialog :open="props.open" @update:open="handleOpenChange">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>{{ title }}</DialogTitle>
        <DialogDescription>{{ hint }}</DialogDescription>
      </DialogHeader>
      <Card class="border-dashed">
        <CardContent class="flex items-center gap-3 p-3">
          <div
            class="flex size-12 shrink-0 items-center justify-center rounded-md border border-border bg-muted/30"
          >
            <img v-if="logoUrl" :src="logoUrl" class="size-7 object-contain" />
            <span v-else class="rounded-sm text-[18px]">
              <span v-if="regionFlagClass" :class="regionFlagClass" />
              <span v-else class="text-xs font-semibold text-muted-foreground">
                {{ regionCode }}
              </span>
            </span>
          </div>
          <div class="min-w-0">
            <p class="truncate text-sm font-semibold text-foreground">{{ profileName }}</p>
            <p class="truncate text-xs text-muted-foreground">{{ profileSubtitle }}</p>
          </div>
        </CardContent>
      </Card>
      <DialogFooter class="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <Button variant="outline" type="button" class="w-full" @click="emit('cancel')">
          {{ cancelLabel }}
        </Button>
        <Button type="button" class="w-full" @click="emit('confirm')">
          {{ confirmLabel }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
