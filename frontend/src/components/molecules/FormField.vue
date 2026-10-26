<script setup lang="ts">
defineProps<{
  label: string;
  /** id of the control inside the slot, so the label is clickable and announced. */
  for?: string;
  error?: string;
  hint?: string;
  required?: boolean;
}>();
</script>

<template>
  <div class="space-y-1.5">
    <label :for="$props.for" class="block text-sm font-medium text-jati">
      {{ label }}<span v-if="required" class="ml-0.5 text-bata" aria-hidden="true">*</span>
    </label>
    <slot />
    <p v-if="error" class="text-xs font-medium text-bata" role="alert">{{ error }}</p>
    <p v-else-if="hint" class="text-xs text-jati-muda">{{ hint }}</p>
  </div>
</template>
