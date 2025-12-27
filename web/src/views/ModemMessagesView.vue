<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ModemMessagesDeleteDialog from '@/components/modem/messages/ModemMessagesDeleteDialog.vue'
import ModemMessagesFab from '@/components/modem/messages/ModemMessagesFab.vue'
import ModemMessagesHeader from '@/components/modem/messages/ModemMessagesHeader.vue'
import ModemMessagesList from '@/components/modem/messages/ModemMessagesList.vue'
import { useModemMessages, type ConversationItem } from '@/composables/useModemMessages'

const route = useRoute()
const router = useRouter()

const modemId = computed(() => (route.params.id ?? 'unknown') as string)

const { items, count, isLoading, deleteConversation } = useModemMessages(modemId)

const deleteOpen = ref(false)
const deleteLoading = ref(false)
const deleteTarget = ref<ConversationItem | null>(null)

const deleteTargetLabel = computed(() => deleteTarget.value?.participantLabel ?? '')
const isFabDisabled = computed(() => modemId.value === 'unknown')

const openDeleteDialog = (item: ConversationItem) => {
  deleteTarget.value = item
  deleteOpen.value = true
}

const closeDeleteDialog = () => {
  deleteOpen.value = false
  deleteTarget.value = null
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await deleteConversation(deleteTarget.value.participantValue)
  } catch (err) {
    console.error('[ModemMessagesView] Failed to delete messages:', err)
  } finally {
    deleteLoading.value = false
    closeDeleteDialog()
  }
}

const startConversation = async () => {
  if (!modemId.value || modemId.value === 'unknown') return
  await router.push({
    name: 'modem-message-thread',
    params: { id: modemId.value, participant: 'new' },
    query: { new: '1' },
  })
}
</script>

<template>
  <div class="space-y-6">
    <ModemMessagesHeader :count="count" :is-loading="isLoading" />

    <ModemMessagesList
      :items="items"
      :modem-id="modemId"
      :is-loading="isLoading"
      @delete="openDeleteDialog"
    />
  </div>

  <ModemMessagesFab :disabled="isFabDisabled" @click="startConversation" />

  <ModemMessagesDeleteDialog
    v-model:open="deleteOpen"
    :target-label="deleteTargetLabel"
    :is-deleting="deleteLoading"
    @confirm="confirmDelete"
  />
</template>
