<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'

const open = defineModel<boolean>('open', { required: true })

const props = defineProps<{
  message: string
}>()

const { t } = useI18n()
</script>

<template>
  <AlertDialog v-model:open="open">
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{{ t('modemDetail.settings.successTitle') }}</AlertDialogTitle>
        <AlertDialogDescription>{{ props.message }}</AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogAction>{{ t('modemDetail.actions.confirm') }}</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
