<script setup lang="ts">
import type { ComboboxViewportProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { ComboboxViewport as RekaComboboxViewport, useForwardProps } from "reka-ui"
import { cn } from "@/lib/utils"

const props = defineProps<ComboboxViewportProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")

const forwarded = useForwardProps(delegatedProps)
</script>

<template>
  <RekaComboboxViewport
    data-slot="combobox-viewport"
    v-bind="forwarded"
    :class="cn('max-h-[300px] scroll-py-1 overflow-x-hidden overflow-y-auto', props.class)"
  >
    <slot />
  </RekaComboboxViewport>
</template>
