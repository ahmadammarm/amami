<script setup lang="ts">
import type { AlertDialogActionProps } from "reka-ui"
import type { HTMLAttributes } from "vue"
import { reactiveOmit } from "@vueuse/core"
import { AlertDialogAction as RekaAlertDialogAction } from "reka-ui"
import { cn } from "@/lib/utils"
import { buttonVariants } from '@/components/atoms/button'

const props = defineProps<AlertDialogActionProps & { class?: HTMLAttributes["class"] }>()

const delegatedProps = reactiveOmit(props, "class")
</script>

<template>
  <RekaAlertDialogAction v-bind="delegatedProps" :class="cn(buttonVariants(), props.class)">
    <slot />
  </RekaAlertDialogAction>
</template>
