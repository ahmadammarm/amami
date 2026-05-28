<script lang="ts" setup>
import type { RangeCalendarGridHeadProps } from "reka-ui"
import { RangeCalendarGridHead as RekaRangeCalendarGridHead } from "reka-ui"

const props = defineProps<RangeCalendarGridHeadProps>()
</script>

<template>
  <RekaRangeCalendarGridHead
    data-slot="range-calendar-grid-head"
    v-bind="props"
  >
    <slot />
  </RekaRangeCalendarGridHead>
</template>
