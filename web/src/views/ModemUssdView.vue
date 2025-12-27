<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import ModemUssdComposer from '@/components/modem/ussd/ModemUssdComposer.vue'
import ModemUssdConversation from '@/components/modem/ussd/ModemUssdConversation.vue'
import ModemUssdHeader from '@/components/modem/ussd/ModemUssdHeader.vue'
import { useUssdSession } from '@/composables/useUssdSession'

const route = useRoute()

const modemId = computed(() => (route.params.id ?? 'unknown') as string)

const { items, draft, isSending, hasEntries, sendMessage } = useUssdSession(modemId)
</script>

<template>
  <div class="flex h-[calc(100dvh-6.5rem)] flex-col overflow-hidden">
    <ModemUssdHeader />

    <ModemUssdConversation :items="items" :has-entries="hasEntries" />

    <ModemUssdComposer v-model="draft" :is-sending="isSending" @submit="sendMessage" />
  </div>
</template>
