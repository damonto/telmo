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

type DownloadProfilePreview = {
  iccid: string
  serviceProviderName: string
  profileName: string
  profileNickname?: string
  profileState: string
  icon?: string
  iconType?: string
  ownerMcc?: string
  ownerMnc?: string
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

const profileName = computed(() => {
  return (
    props.profile?.profileName ||
    props.profile?.serviceProviderName ||
    props.profile?.profileNickname ||
    ''
  )
})

const profileSubtitle = computed(() => props.profile?.iccid ?? '')

const ownerText = computed(() => {
  const mcc = props.profile?.ownerMcc ?? ''
  const mnc = props.profile?.ownerMnc ?? ''
  return `${mcc}${mnc}`.trim()
})

const logoUrl = computed(() => {
  if (!props.profile?.icon || !props.profile.iconType) return ''
  return `data:${props.profile.iconType};base64,${props.profile.icon}`
})

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
            <span v-else class="text-xs font-semibold text-muted-foreground">{{ ownerText }}</span>
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
