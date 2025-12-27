<script setup lang="ts">
import { computed } from 'vue'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Progress } from '@/components/ui/progress'

const props = defineProps<{
  open: boolean
  title: string
  stageLabel: string
  progress: number
  cancelLabel: string
}>()

const emit = defineEmits<{
  (event: 'cancel'): void
}>()

const progressValue = computed(() => Math.min(Math.max(props.progress, 0), 100))
</script>

<template>
  <Dialog :open="props.open">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>{{ title }}</DialogTitle>
        <DialogDescription>{{ stageLabel }}</DialogDescription>
      </DialogHeader>
      <Progress :model-value="progressValue" />
      <DialogFooter>
        <Button variant="ghost" type="button" @click="emit('cancel')">
          {{ cancelLabel }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
