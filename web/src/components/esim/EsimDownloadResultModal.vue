<script setup lang="ts">
import { computed } from 'vue'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'

const props = defineProps<{
  open: boolean
  title: string
  message: string
  confirmLabel: string
  tone: 'success' | 'error'
}>()

const emit = defineEmits<{
  (event: 'confirm'): void
}>()

const messageClass = computed(() =>
  props.tone === 'error' ? 'text-destructive' : 'text-muted-foreground',
)
</script>

<template>
  <AlertDialog :open="props.open">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ title }}</AlertDialogTitle>
        <AlertDialogDescription :class="messageClass">
          {{ message }}
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogAction @click="emit('confirm')">
          {{ confirmLabel }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
