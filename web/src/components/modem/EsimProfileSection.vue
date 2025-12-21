<script setup lang="ts">
import { EllipsisVertical } from 'lucide-vue-next'
import * as v from 'valibot'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { Badge } from '@/components/ui/badge'
import { useModemDisplay } from '@/composables/useModemDisplay'
import type { EsimProfile } from '@/types/modem'

const profiles = defineModel<EsimProfile[]>('profiles', { required: true })

const { t } = useI18n()
const { flagClass } = useModemDisplay()

const profileCount = computed(() => profiles.value.length)
const hasProfiles = computed(() => profiles.value.length > 0)

const openMenuId = ref<string | null>(null)

const toggleOpen = ref(false)
const toggleProfile = ref<EsimProfile | null>(null)
const toggleNextValue = ref(false)
const toggleLoading = ref(false)

const renameOpen = ref(false)
const renameProfile = ref<EsimProfile | null>(null)
const renameName = ref('')
const renameError = ref('')
const renameLoading = ref(false)

const deleteOpen = ref(false)
const deleteProfile = ref<EsimProfile | null>(null)
const deleteLoading = ref(false)

const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

const openToggleDialog = (profile: EsimProfile) => {
  toggleOpen.value = true
  toggleProfile.value = profile
  toggleNextValue.value = !profile.enabled
}

const closeToggleDialog = () => {
  toggleOpen.value = false
  toggleProfile.value = null
  toggleNextValue.value = false
  toggleLoading.value = false
}

const confirmToggle = async () => {
  if (!toggleProfile.value) return
  toggleLoading.value = true
  await sleep(650)
  toggleProfile.value.enabled = toggleNextValue.value
  closeToggleDialog()
}

const openRenameDialog = (profile: EsimProfile) => {
  renameOpen.value = true
  renameProfile.value = profile
  renameName.value = profile.name
  renameError.value = ''
}

const closeRenameDialog = () => {
  renameOpen.value = false
  renameProfile.value = null
  renameName.value = ''
  renameError.value = ''
  renameLoading.value = false
}

const validateRename = () => {
  const schema = v.pipe(
    v.string(),
    v.trim(),
    v.minLength(1, t('modemDetail.validation.required')),
    v.maxBytes(64, t('modemDetail.validation.maxBytes')),
  )

  const result = v.safeParse(schema, renameName.value)
  if (!result.success) {
    renameError.value = result.issues[0]?.message ?? t('modemDetail.validation.required')
    return null
  }

  renameError.value = ''
  return result.output
}

const confirmRename = async () => {
  const value = validateRename()
  if (!value || !renameProfile.value) return
  renameLoading.value = true
  await sleep(650)
  renameProfile.value.name = value
  closeRenameDialog()
}

const openDeleteDialog = (profile: EsimProfile) => {
  if (profile.enabled) return
  deleteOpen.value = true
  deleteProfile.value = profile
}

const closeDeleteDialog = () => {
  deleteOpen.value = false
  deleteProfile.value = null
  deleteLoading.value = false
}

const confirmDelete = async () => {
  if (!deleteProfile.value) return
  deleteLoading.value = true
  await sleep(650)
  profiles.value = profiles.value.filter((profile) => profile.id !== deleteProfile.value?.id)
  closeDeleteDialog()
}

const closeMenu = (event: Event) => {
  const target = event.currentTarget as HTMLElement
  const details = target.closest('details')
  if (details) {
    details.removeAttribute('open')
  }
  openMenuId.value = null
}

const handleMenuToggle = (profileId: string, event: Event) => {
  const details = event.currentTarget as HTMLDetailsElement
  openMenuId.value = details.open ? profileId : null
}

const handleRenameClick = (profile: EsimProfile, event: Event) => {
  openRenameDialog(profile)
  closeMenu(event)
}

const handleDeleteClick = (profile: EsimProfile, event: Event) => {
  openDeleteDialog(profile)
  closeMenu(event)
}
</script>

<template>
  <section class="space-y-3">
    <div class="flex items-center justify-between">
      <h2 class="text-sm font-semibold text-muted-foreground">
        {{ t('modemDetail.esim.listTitle') }}
      </h2>
      <Badge variant="outline" class="text-[10px] uppercase tracking-[0.2em]">
        {{ profileCount }}
      </Badge>
    </div>

    <div
      v-if="!hasProfiles"
      class="rounded-2xl border border-dashed border-border p-6 text-sm text-muted-foreground"
    >
      {{ t('modemDetail.esim.noProfiles') }}
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="profile in profiles"
        :key="profile.id"
        class="relative flex items-center justify-between rounded-2xl border border-white/40 bg-white/80 px-3 py-2.5 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/60"
        :class="openMenuId === profile.id ? 'z-20' : 'z-0'"
      >
        <div class="flex min-w-0 items-center gap-3">
          <div
            class="flex size-11 shrink-0 items-center justify-center rounded-2xl bg-white/80 ring-1 ring-border/60 dark:bg-slate-950/60"
          >
            <img
              v-if="profile.logoUrl"
              :src="profile.logoUrl"
              :alt="`${profile.name} logo`"
              class="size-6 object-contain"
            />
            <span v-else class="rounded-sm text-[18px]">
              <span v-if="flagClass(profile.regionCode)" :class="flagClass(profile.regionCode)" />
              <span v-else class="text-xs font-semibold text-muted-foreground">
                {{ profile.regionCode }}
              </span>
            </span>
          </div>
          <div class="min-w-0">
            <p class="truncate text-sm font-semibold text-foreground">
              {{ profile.name }}
            </p>
            <p class="truncate text-xs text-muted-foreground">
              {{ profile.iccid }}
            </p>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button
            type="button"
            role="switch"
            class="relative inline-flex h-6 w-11 items-center rounded-full transition"
            :class="profile.enabled ? 'bg-emerald-500' : 'bg-muted'"
            :aria-checked="profile.enabled"
            @click="openToggleDialog(profile)"
          >
            <span
              class="inline-block size-4 transform rounded-full bg-white transition"
              :class="profile.enabled ? 'translate-x-6' : 'translate-x-1'"
            />
          </button>

          <details class="relative z-20" @toggle="handleMenuToggle(profile.id, $event)">
            <summary
              class="flex cursor-pointer list-none items-center justify-center rounded-full p-2 text-muted-foreground transition hover:bg-muted"
            >
              <EllipsisVertical class="size-4" />
            </summary>
            <div
              class="absolute right-0 mt-2 w-40 rounded-xl border border-white/60 bg-white/90 p-1 text-sm shadow-lg backdrop-blur-lg"
            >
              <button
                type="button"
                class="w-full rounded-lg px-3 py-2 text-left text-foreground transition hover:bg-muted"
                @click="handleRenameClick(profile, $event)"
              >
                {{ t('modemDetail.actions.rename') }}
              </button>
              <button
                type="button"
                class="w-full rounded-lg px-3 py-2 text-left transition hover:bg-muted"
                :class="
                  profile.enabled
                    ? 'cursor-not-allowed text-muted-foreground/60'
                    : 'text-rose-500'
                "
                :disabled="profile.enabled"
                @click="handleDeleteClick(profile, $event)"
              >
                {{ t('modemDetail.actions.delete') }}
              </button>
            </div>
          </details>
        </div>
      </div>
    </div>
  </section>

  <div
    v-if="toggleOpen"
    class="fixed inset-0 z-30 flex items-center justify-center bg-black/30 px-4 backdrop-blur-sm"
  >
    <div class="w-full max-w-sm rounded-2xl bg-white p-6 shadow-2xl">
      <h3 class="text-sm font-semibold text-foreground">
        {{
          toggleNextValue
            ? t('modemDetail.confirm.enable', { name: toggleProfile?.name ?? '' })
            : t('modemDetail.confirm.disable', { name: toggleProfile?.name ?? '' })
        }}
      </h3>
      <div class="mt-6 flex gap-3">
        <button
          type="button"
          class="flex-1 rounded-full border border-border px-4 py-2 text-sm font-semibold text-muted-foreground"
          @click="closeToggleDialog"
          :disabled="toggleLoading"
        >
          {{ t('modemDetail.actions.cancel') }}
        </button>
        <button
          type="button"
          class="flex-1 rounded-full bg-foreground px-4 py-2 text-sm font-semibold text-background"
          @click="confirmToggle"
          :disabled="toggleLoading"
        >
          <span v-if="toggleLoading" class="inline-flex items-center gap-2">
            <span
              class="size-4 animate-spin rounded-full border-2 border-background/60 border-t-background"
            />
            {{ t('modemDetail.actions.confirm') }}
          </span>
          <span v-else>{{ t('modemDetail.actions.confirm') }}</span>
        </button>
      </div>
    </div>
  </div>

  <div
    v-if="renameOpen"
    class="fixed inset-0 z-30 flex items-center justify-center bg-black/30 px-4 backdrop-blur-sm"
  >
    <div class="w-full max-w-sm rounded-2xl bg-white p-6 shadow-2xl">
      <h3 class="text-sm font-semibold text-foreground">
        {{ t('modemDetail.actions.rename') }}
      </h3>
      <div class="mt-4 space-y-2">
        <input
          v-model="renameName"
          type="text"
          class="w-full rounded-2xl border border-border/60 px-3 py-2 text-sm text-foreground outline-none transition focus:border-ring"
          :placeholder="t('modemDetail.actions.rename')"
        />
        <p v-if="renameError" class="text-xs text-rose-500">
          {{ renameError }}
        </p>
      </div>
      <div class="mt-6 flex gap-3">
        <button
          type="button"
          class="flex-1 rounded-full border border-border px-4 py-2 text-sm font-semibold text-muted-foreground"
          @click="closeRenameDialog"
          :disabled="renameLoading"
        >
          {{ t('modemDetail.actions.cancel') }}
        </button>
        <button
          type="button"
          class="flex-1 rounded-full bg-foreground px-4 py-2 text-sm font-semibold text-background"
          @click="confirmRename"
          :disabled="renameLoading"
        >
          <span v-if="renameLoading" class="inline-flex items-center gap-2">
            <span
              class="size-4 animate-spin rounded-full border-2 border-background/60 border-t-background"
            />
            {{ t('modemDetail.actions.update') }}
          </span>
          <span v-else>{{ t('modemDetail.actions.update') }}</span>
        </button>
      </div>
    </div>
  </div>

  <div
    v-if="deleteOpen"
    class="fixed inset-0 z-30 flex items-center justify-center bg-black/30 px-4 backdrop-blur-sm"
  >
    <div class="w-full max-w-sm rounded-2xl bg-white p-6 shadow-2xl">
      <h3 class="text-sm font-semibold text-foreground">
        {{ t('modemDetail.confirm.delete', { name: deleteProfile?.name ?? '' }) }}
      </h3>
      <div class="mt-6 flex gap-3">
        <button
          type="button"
          class="flex-1 rounded-full border border-border px-4 py-2 text-sm font-semibold text-muted-foreground"
          @click="closeDeleteDialog"
          :disabled="deleteLoading"
        >
          {{ t('modemDetail.actions.cancel') }}
        </button>
        <button
          type="button"
          class="flex-1 rounded-full bg-rose-500 px-4 py-2 text-sm font-semibold text-white"
          @click="confirmDelete"
          :disabled="deleteLoading"
        >
          <span v-if="deleteLoading" class="inline-flex items-center gap-2">
            <span class="size-4 animate-spin rounded-full border-2 border-white/60 border-t-white" />
            {{ t('modemDetail.actions.confirm') }}
          </span>
          <span v-else>{{ t('modemDetail.actions.confirm') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>
