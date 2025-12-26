<script setup lang="ts">
import { computed, watch } from 'vue'
import { ScanQrCode } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
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

const domainPattern = /^(?=.{1,253}$)(?!-)([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,63}$/

const smdpPlaceholder = computed(() => t('modemDetail.esim.smdp'))
const activationPlaceholder = computed(() => t('modemDetail.esim.activationCode'))
const confirmationPlaceholder = computed(() => t('modemDetail.esim.confirmationCode'))

const normalizeSmdp = (value: string) => {
  const trimmed = value.trim()
  if (!trimmed) return ''
  const normalized = trimmed.includes('://') ? trimmed : `https://${trimmed}`
  try {
    const url = new URL(normalized)
    const host = url.hostname.toLowerCase()
    if (!domainPattern.test(host)) {
      return ''
    }
    return host
  } catch {
    return ''
  }
}

const installSchemaDefinition = z.object({
  smdp: z
    .string()
    .trim()
    .min(1, t('modemDetail.esim.validation.smdpRequired'))
    .refine((value) => normalizeSmdp(value).length > 0, t('modemDetail.esim.validation.smdpInvalid'))
    .transform((value) => normalizeSmdp(value)),
  activationCode: z
    .string()
    .trim()
    .min(1, t('modemDetail.esim.validation.activationCodeRequired')),
  confirmationCode: z.string().optional().transform((value) => value?.trim() ?? ''),
})

type InstallFormValues = z.infer<typeof installSchemaDefinition>

const installSchema = toTypedSchema(installSchemaDefinition)

const { handleSubmit, resetForm, isSubmitting } = useForm<InstallFormValues>({
  validationSchema: installSchema,
  initialValues: {
    smdp: '',
    activationCode: '',
    confirmationCode: '',
  },
})

const resetValues = () => {
  resetForm({
    values: {
      smdp: '',
      activationCode: '',
      confirmationCode: '',
    },
  })
}

const closeDialog = () => {
  open.value = false
  resetValues()
}

const onSubmit = handleSubmit((values) => {
  emit('confirm', {
    smdp: values.smdp,
    activationCode: values.activationCode,
    confirmationCode: values.confirmationCode ?? '',
  })
  closeDialog()
})

watch(open, (value) => {
  if (!value) {
    resetValues()
  }
})
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <div class="flex items-start justify-between gap-3">
          <DialogTitle>{{ t('modemDetail.esim.installTitle') }}</DialogTitle>
          <Button
            variant="outline"
            size="icon"
            type="button"
            :aria-label="t('modemDetail.esim.scan')"
            :title="t('modemDetail.esim.scan')"
          >
            <ScanQrCode class="size-4" />
          </Button>
        </div>
      </DialogHeader>

      <form class="space-y-4" @submit="onSubmit">
        <FormField v-slot="{ componentField }" name="smdp">
          <FormItem>
            <FormLabel>{{ t('modemDetail.esim.smdp') }}</FormLabel>
            <FormControl>
              <Input type="text" :placeholder="smdpPlaceholder" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="activationCode">
          <FormItem>
            <FormLabel>{{ t('modemDetail.esim.activationCode') }}</FormLabel>
            <FormControl>
              <Input type="text" :placeholder="activationPlaceholder" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="confirmationCode">
          <FormItem>
            <FormLabel>{{ t('modemDetail.esim.confirmationCode') }}</FormLabel>
            <FormControl>
              <Input type="text" :placeholder="confirmationPlaceholder" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <DialogFooter class="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <Button variant="outline" type="button" class="w-full" @click="closeDialog">
            {{ t('modemDetail.actions.cancel') }}
          </Button>
          <Button type="submit" class="w-full" :disabled="isSubmitting">
            {{ t('modemDetail.esim.installConfirm') }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
