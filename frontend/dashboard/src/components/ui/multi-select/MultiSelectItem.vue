<script setup lang="ts">
import { computed } from 'vue'

import { cn } from '@/lib/utils'

interface Props {
  value: string
  label: string
  selected?: boolean
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  selected: false,
})

const emit = defineEmits<{
  select: [value: string]
}>()

const handleSelect = () => {
  emit('select', props.value)
}
</script>

<template>
  <div
    :class="cn(
      'relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none aria-selected:bg-accent aria-selected:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
      props.class
    )"
    :data-selected="selected"
    @click="handleSelect"
  >
    <slot />
  </div>
</template>
