<script setup lang="ts">
import type { ComboboxSeparatorProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { ComboboxSeparator as RekaComboboxSeparator } from "reka-ui"
import { cn } from "@/lib/utils"

const props = defineProps<ComboboxSeparatorProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")
</script>

<template>
  <RekaComboboxSeparator
    data-slot="combobox-separator"
    v-bind="delegatedProps"
    :class="cn('bg-border -mx-1 h-px', props.class)"
  >
    <slot />
  </RekaComboboxSeparator>
</template>
