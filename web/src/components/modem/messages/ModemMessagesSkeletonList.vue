<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    count?: number
  }>(),
  {
    count: 4,
  },
)
</script>

<template>
  <div class="space-y-2">
    <div
      v-for="index in props.count"
      :key="`message-skeleton-${index}`"
      class="flex items-center justify-between rounded-lg border border-border bg-card px-4 py-3"
    >
      <div class="flex min-w-0 flex-1 flex-col gap-1.5">
        <div class="h-4 w-40 animate-pulse rounded bg-muted/60" />
        <div class="h-3 w-full animate-pulse rounded bg-muted/40" />
      </div>
      <div class="h-5 w-16 animate-pulse rounded-full bg-muted/50" />
    </div>
  </div>
</template>
