<script setup lang="ts">
import type { MenubarGroupProps } from "reka-ui"
import { MenubarGroup as RekaMenubarGroup } from "reka-ui"

const props = defineProps<MenubarGroupProps>()
</script>

<template>
  <RekaMenubarGroup
    data-slot="menubar-group"
    v-bind="props"
  >
    <slot />
  </RekaMenubarGroup>
</template>
