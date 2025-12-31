<script setup lang="ts">
import { computed } from 'vue'
import { SignalHigh } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { useModemDisplay } from '@/composables/useModemDisplay'

const props = defineProps<{
  name: string
  regionCode: string
  operatorName: string
  registeredOperatorName: string
  registeredOperatorCode: string
  registrationState: string
  accessTechnology: string | null
  supportsEsim: boolean
  number: string
  signalQuality: number
}>()

const { t } = useI18n()
const {
  flagClass,
  formatSignal,
  signalIcon,
  signalTone,
  registrationStateIcon,
  registrationStateLabel,
  registrationStateTone,
  shouldShowRegistrationIcon,
  getSignalColorOverride,
} = useModemDisplay()

const getTech = (accessTechnology: string | null): '4G' | '5G' => {
  if (!accessTechnology) return '4G'
  const upperTech = accessTechnology.toUpperCase()
  if (upperTech.includes('5G') || upperTech.includes('NR') || upperTech.includes('SA')) {
    return '5G'
  }
  return '4G'
}

const regionFlagClass = computed(() => flagClass(props.regionCode))
const tech = computed(() => getTech(props.accessTechnology))
const techVariant = computed(() => (tech.value === '5G' ? 'default' : 'secondary'))
const isRoaming = computed(() => props.registrationState === 'Roaming')
const showRoamingLabel = computed(() => isRoaming.value && Boolean(props.registeredOperatorCode))
const roamingLabel = computed(() => props.registeredOperatorName || props.registeredOperatorCode)
const esimLabel = computed(() => (props.supportsEsim ? t('labels.esim') : t('labels.psim')))
const displayName = computed(() =>
  props.name.trim().length > 0 ? props.name : props.operatorName,
)
const displayNumber = computed(() => (props.number.trim() ? props.number : t('home.noNumber')))
const isSearching = computed(() => props.registrationState.trim() === 'Searching')
const signalValue = computed(() => formatSignal(props.signalQuality))
const signalIconComponent = computed(() =>
  isSearching.value ? SignalHigh : signalIcon(props.signalQuality),
)
const signalToneClass = computed(() => {
  if (isSearching.value) return 'text-amber-500'
  const override = getSignalColorOverride(props.registrationState)
  return override ?? signalTone(props.signalQuality)
})
const signalTitle = computed(() => `${t('labels.signal')}: ${signalValue.value}`)
const showRegistrationIcon = computed(() =>
  shouldShowRegistrationIcon(props.registrationState),
)
const registrationIcon = computed(() => registrationStateIcon(props.registrationState))
const registrationLabel = computed(() => registrationStateLabel(props.registrationState))
const registrationToneClass = computed(() => registrationStateTone(props.registrationState))
</script>

<template>
  <Card class="h-full border-0 py-4 shadow-sm transition duration-300 group-hover:-translate-y-0.5">
    <CardContent class="flex items-center gap-3 px-4">
      <div class="flex size-12 shrink-0 items-center justify-center rounded-xl border bg-background">
        <span
          v-if="regionFlagClass"
          :class="regionFlagClass"
          class="rounded-sm text-xl"
          :aria-label="props.regionCode"
          :title="props.regionCode"
        />
        <span
          v-else
          class="rounded-sm text-base font-semibold text-muted-foreground"
          :aria-label="props.regionCode"
          :title="props.regionCode"
        >
          {{ props.regionCode }}
        </span>
      </div>
      <div class="flex min-w-0 flex-1 flex-col gap-0.5">
        <div class="flex items-center justify-between gap-2">
          <p
            class="min-w-0 flex-1 truncate text-sm font-semibold text-foreground"
            :title="displayName"
          >
            {{ displayName }}
          </p>
          <div class="flex shrink-0 items-center gap-1.5">
            <Badge :variant="techVariant">
              {{ tech }}
            </Badge>
            <Badge variant="outline">
              {{ esimLabel }}
            </Badge>
          </div>
        </div>
        <p class="truncate text-xs font-normal text-foreground/70">
          <span>{{ props.operatorName }}</span>
          <span v-if="showRoamingLabel" class="ml-1 text-muted-foreground">
            ({{ roamingLabel }})
          </span>
        </p>
        <div class="mt-auto flex items-center justify-between gap-3">
          <p class="truncate text-xs text-muted-foreground">
            {{ displayNumber }}
          </p>
          <div class="flex items-center gap-1">
            <component
              :is="signalIconComponent"
              class="size-4 shrink-0"
              :class="[signalToneClass, isSearching && 'animate-pulse']"
              :title="signalTitle"
            />
            <component
              v-if="showRegistrationIcon && registrationIcon"
              :is="registrationIcon"
              class="size-4 shrink-0"
              :class="registrationToneClass"
              :aria-label="props.registrationState"
              :title="props.registrationState"
            />
            <Badge
              v-else-if="showRegistrationIcon && registrationLabel"
              variant="secondary"
              class="h-4 px-1.5 text-xs font-semibold"
              :class="registrationToneClass"
              :aria-label="props.registrationState"
              :title="props.registrationState"
            >
              {{ registrationLabel }}
            </Badge>
            <span
              class="inline-flex h-4 min-w-6 items-center justify-end text-right font-mono text-xs text-muted-foreground tabular-nums"
            >
              {{ signalValue }}
            </span>
          </div>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
