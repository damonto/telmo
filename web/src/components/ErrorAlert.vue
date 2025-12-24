<script setup lang="ts">
import { watch } from 'vue'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'

interface Props {
  title?: string
  message?: string
}

withDefaults(defineProps<Props>(), {
  title: 'Error',
  message: '',
})

const emit = defineEmits<{
  close: []
}>()

const open = defineModel<boolean>('open', { default: false })

watch(open, (value) => {
  if (!value) {
    emit('close')
  }
})

const handleClose = () => {
  open.value = false
}
</script>

<template>
  <AlertDialog v-model:open="open">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ title }}</AlertDialogTitle>
        <AlertDialogDescription>{{ message }}</AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogAction @click="handleClose">OK</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
