<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { useModemDisplay } from '@/composables/useModemDisplay'

const props = defineProps<{
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
const { flagClass, formatSignal, signalIcon, signalTone } = useModemDisplay()

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
const displayNumber = computed(() => (props.number.trim() ? props.number : t('home.noNumber')))
const signalValue = computed(() => formatSignal(props.signalQuality))
const signalIconComponent = computed(() => signalIcon(props.signalQuality))
const signalToneClass = computed(() => signalTone(props.signalQuality))
const signalTitle = computed(() => `${t('labels.signal')}: ${signalValue.value}`)
</script>

<template>
  <Card class="h-full border-0 py-4 shadow-sm transition duration-300 group-hover:-translate-y-0.5">
    <CardContent class="flex items-center gap-3 px-4">
      <div class="flex size-9 shrink-0 items-center justify-center rounded-lg border bg-background">
        <span
          v-if="regionFlagClass"
          :class="regionFlagClass"
          class="rounded-sm text-base"
          :aria-label="props.regionCode"
          :title="props.regionCode"
        />
        <span
          v-else
          class="rounded-sm text-xs font-semibold text-muted-foreground"
          :aria-label="props.regionCode"
          :title="props.regionCode"
        >
          {{ props.regionCode }}
        </span>
      </div>
      <div class="flex min-w-0 flex-1 flex-col gap-1">
        <div class="flex items-center justify-between gap-2">
          <p class="truncate text-sm font-semibold text-foreground">
            <span class="text-foreground">{{ props.operatorName }}</span>
            <span v-if="showRoamingLabel" class="text-xs text-muted-foreground">
              ({{ roamingLabel }})
            </span>
          </p>
          <div class="flex items-center gap-1.5">
            <Badge :variant="techVariant">
              {{ tech }}
            </Badge>
            <Badge variant="outline">
              {{ esimLabel }}
            </Badge>
          </div>
        </div>
        <div class="mt-auto flex items-center justify-between gap-3">
          <p class="truncate text-xs text-muted-foreground">
            {{ displayNumber }}
          </p>
          <div class="flex items-center gap-1">
            <component
              :is="signalIconComponent"
              class="size-4 shrink-0"
              :class="signalToneClass"
              :title="signalTitle"
            />
            <Badge
              v-if="isRoaming"
              variant="secondary"
              class="h-4 px-1.5 text-xs font-semibold"
              :aria-label="t('labels.roaming')"
              :title="t('labels.roaming')"
            >
              R
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
