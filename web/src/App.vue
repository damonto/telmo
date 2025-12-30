<script setup lang="ts">
import { RouterView } from 'vue-router'
import 'vue-sonner/style.css'

import ErrorAlert from '@/components/ErrorAlert.vue'
import { Toaster } from '@/components/ui/sonner'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { isErrorOpen, errorTitle, errorMessage, clearError } = useErrorHandler()
</script>

<template>
  <RouterView />
  <ErrorAlert
    v-model:open="isErrorOpen"
    :title="errorTitle"
    :message="errorMessage"
    @close="clearError"
  />
  <Toaster position="top-center" />
</template>
