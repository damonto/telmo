<script setup lang="ts">
import { ref, watch } from 'vue'

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
  open?: boolean
  title?: string
  message?: string
}

const props = withDefaults(defineProps<Props>(), {
  open: false,
  title: 'Error',
  message: '',
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  close: []
}>()

const isOpen = ref(props.open)

watch(
  () => props.open,
  (value) => {
    isOpen.value = value
  },
)

watch(isOpen, (value) => {
  if (!value) {
    emit('close')
    emit('update:open', false)
  }
})

const handleClose = () => {
  isOpen.value = false
}
</script>

<template>
  <AlertDialog :open="isOpen">
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
