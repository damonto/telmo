<script setup lang="ts">
import { RouterView } from 'vue-router'

import ErrorAlert from '@/components/ErrorAlert.vue'
import { useErrorHandler } from '@/composables/useErrorHandler'

const { isErrorOpen, errorTitle, errorMessage, clearError } = useErrorHandler()
</script>

<template>
  <RouterView />
  <ErrorAlert
    :open="isErrorOpen"
    :title="errorTitle"
    :message="errorMessage"
    @close="clearError"
    @update:open="clearError"
  />
</template>
