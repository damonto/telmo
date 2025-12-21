<script setup lang="ts">
import type { Component } from 'vue'
import { computed } from 'vue'
import { RouterLink, useRoute, type RouteLocationRaw } from 'vue-router'

export type TabBarItem = {
  key: string
  to: RouteLocationRaw
  routeName: string
  label: string
  icon: Component
}

const props = withDefaults(
  defineProps<{
    items: TabBarItem[]
    containerClass?: string
  }>(),
  {
    containerClass: 'max-w-6xl',
  },
)

const route = useRoute()

const tabs = computed(() =>
  props.items.map((item) => ({
    ...item,
    isActive: route.name === item.routeName,
  })),
)
</script>

<template>
  <nav
    aria-label="Primary navigation"
    class="fixed bottom-0 left-0 right-0 z-20 h-16 w-full border-t border-white/40 bg-white/70 px-6 py-2 shadow-[0_-12px_30px_rgba(15,23,42,0.08)] backdrop-blur-2xl dark:border-white/10 dark:bg-slate-950/60"
  >
    <div
      class="mx-auto flex h-full w-full items-center justify-around"
      :class="props.containerClass"
    >
      <RouterLink
        v-for="item in tabs"
        :key="item.key"
        :to="item.to"
        class="flex flex-1 items-center justify-center"
        :aria-current="item.isActive ? 'page' : undefined"
      >
        <span
          class="flex size-10 items-center justify-center rounded-2xl transition"
          :class="
            item.isActive
              ? 'bg-black/5 text-foreground dark:bg-white/10'
              : 'text-muted-foreground'
          "
        >
          <component :is="item.icon" class="size-5" />
        </span>
        <span class="sr-only">{{ item.label }}</span>
      </RouterLink>
    </div>
  </nav>
</template>
