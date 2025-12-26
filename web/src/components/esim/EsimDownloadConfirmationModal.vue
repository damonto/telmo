<script setup lang="ts">
import { watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'

const props = defineProps<{
  open: boolean
  title: string
  hint: string
  placeholder: string
  confirmLabel: string
  cancelLabel: string
}>()

const emit = defineEmits<{
  (event: 'submit'): void
  (event: 'cancel'): void
}>()

const { t } = useI18n()
const code = defineModel<string>('code', { required: true })

const confirmationSchemaDefinition = z.object({
  code: z.string().trim().min(1, t('modemDetail.validation.required')),
})

type ConfirmationValues = z.infer<typeof confirmationSchemaDefinition>

const confirmationSchema = toTypedSchema(confirmationSchemaDefinition)

const { handleSubmit, resetForm, isSubmitting } = useForm<ConfirmationValues>({
  validationSchema: confirmationSchema,
  initialValues: {
    code: '',
  },
})

const resetValues = () => {
  resetForm({
    values: {
      code: code.value,
    },
  })
}

const handleOpenChange = (nextOpen: boolean) => {
  if (!nextOpen) {
    code.value = ''
    emit('cancel')
  }
}

const onSubmit = handleSubmit((values) => {
  code.value = values.code.trim()
  emit('submit')
})

watch(
  () => props.open,
  (value) => {
    if (!value) {
      resetForm({ values: { code: '' } })
      return
    }
    resetValues()
  },
)
</script>

<template>
  <Dialog :open="props.open" @update:open="handleOpenChange">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>{{ title }}</DialogTitle>
        <DialogDescription>{{ hint }}</DialogDescription>
      </DialogHeader>
      <form class="space-y-4" @submit="onSubmit">
        <FormField v-slot="{ componentField }" name="code">
          <FormItem>
            <FormLabel class="sr-only">{{ placeholder }}</FormLabel>
            <FormControl>
              <Input type="text" :placeholder="placeholder" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <DialogFooter class="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <Button variant="outline" type="button" class="w-full" @click="emit('cancel')">
            {{ cancelLabel }}
          </Button>
          <Button type="submit" class="w-full" :disabled="isSubmitting">
            {{ confirmLabel }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
