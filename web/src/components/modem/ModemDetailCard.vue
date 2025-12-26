<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import type { Modem } from '@/types/modem'

defineProps<{
  modem: Modem
}>()
</script>

<template>
  <div class="space-y-6">
    <!-- Basic Info -->
    <Card>
      <CardHeader class="pb-0">
        <CardTitle class="text-lg">Basic Information</CardTitle>
      </CardHeader>
      <CardContent class="grid gap-3 pt-4 text-sm">
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
      </CardContent>
    </Card>

    <!-- SIM Info -->
    <Card>
      <CardHeader class="pb-0">
        <CardTitle class="text-lg">SIM Information</CardTitle>
      </CardHeader>
      <CardContent class="grid gap-3 pt-4 text-sm">
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
      </CardContent>
    </Card>

    <!-- Network Info -->
    <Card>
      <CardHeader class="pb-0">
        <CardTitle class="text-lg">Network Information</CardTitle>
      </CardHeader>
      <CardContent class="grid gap-3 pt-4 text-sm">
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
      </CardContent>
    </Card>

    <!-- Features -->
    <Card>
      <CardHeader class="pb-0">
        <CardTitle class="text-lg">Features</CardTitle>
      </CardHeader>
      <CardContent class="flex flex-wrap gap-2 pt-4">
        <Badge v-if="modem.supportsEsim">eSIM Support</Badge>
        <Badge v-else variant="secondary">Physical SIM Only</Badge>
      </CardContent>
    </Card>
  </div>
</template>
