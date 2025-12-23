<script setup lang="ts">
import type { Modem } from '@/types/modem'

defineProps<{
  modem: Modem
}>()
</script>

<template>
  <div class="space-y-6">
    <!-- Basic Info -->
    <div class="rounded-lg border bg-card p-6 shadow-sm">
      <h3 class="mb-4 text-lg font-semibold">Basic Information</h3>
      <div class="grid gap-3 text-sm">
        <div class="flex justify-between">
          <span class="text-muted-foreground">ID</span>
          <span class="font-mono">{{ modem.id }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Name</span>
          <span>{{ modem.name }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Manufacturer</span>
          <span>{{ modem.manufacturer }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Firmware Revision</span>
          <span class="font-mono">{{ modem.firmwareRevision }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Hardware Revision</span>
          <span class="font-mono">{{ modem.hardwareRevision }}</span>
        </div>
        <div v-if="modem.number" class="flex justify-between">
          <span class="text-muted-foreground">Number</span>
          <span class="font-mono">{{ modem.number }}</span>
        </div>
      </div>
    </div>

    <!-- SIM Info -->
    <div class="rounded-lg border bg-card p-6 shadow-sm">
      <h3 class="mb-4 text-lg font-semibold">SIM Information</h3>
      <div class="grid gap-3 text-sm">
        <div class="flex justify-between">
          <span class="text-muted-foreground">Operator</span>
          <span>{{ modem.sim.operatorName }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Operator ID</span>
          <span class="font-mono">{{ modem.sim.operatorIdentifier }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Region</span>
          <span class="font-mono">{{ modem.sim.regionCode }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">ICCID</span>
          <span class="font-mono">{{ modem.sim.identifier }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Active</span>
          <span>{{ modem.sim.active ? 'Yes' : 'No' }}</span>
        </div>
      </div>
    </div>

    <!-- Network Info -->
    <div class="rounded-lg border bg-card p-6 shadow-sm">
      <h3 class="mb-4 text-lg font-semibold">Network Information</h3>
      <div class="grid gap-3 text-sm">
        <div class="flex justify-between">
          <span class="text-muted-foreground">Access Technology</span>
          <span>{{ modem.accessTechnology || 'N/A' }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Registration State</span>
          <span>{{ modem.registrationState }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Registered Operator</span>
          <span>{{ modem.registeredOperator.name }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Operator Code</span>
          <span class="font-mono">{{ modem.registeredOperator.code }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-muted-foreground">Signal Quality</span>
          <span class="font-mono">{{ modem.signalQuality }}%</span>
        </div>
      </div>
    </div>

    <!-- Features -->
    <div class="rounded-lg border bg-card p-6 shadow-sm">
      <h3 class="mb-4 text-lg font-semibold">Features</h3>
      <div class="flex gap-2 flex-wrap">
        <span
          v-if="modem.supportsEsim"
          class="rounded-full bg-primary px-3 py-1 text-xs font-medium text-primary-foreground"
        >
          eSIM Support
        </span>
        <span
          v-else
          class="rounded-full bg-secondary px-3 py-1 text-xs font-medium text-secondary-foreground"
        >
          Physical SIM Only
        </span>
      </div>
    </div>
  </div>
</template>
