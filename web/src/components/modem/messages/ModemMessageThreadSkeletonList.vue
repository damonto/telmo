<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    count?: number
  }>(),
  {
    count: 6,
  },
)
</script>

<template>
  <div class="space-y-3">
    <div v-for="index in props.count" :key="`thread-skeleton-${index}`" class="flex">
      <div class="h-8 w-2/3 animate-pulse rounded-2xl bg-muted/60" />
    </div>
  </div>
</template>
