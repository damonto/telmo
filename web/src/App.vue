<script setup lang="ts">
import { RouterView } from 'vue-router'

import ErrorAlert from '@/components/ErrorAlert.vue'
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
</template>
