<script lang="ts" setup>
import type { ToasterProps } from "vue-sonner"
import { Toaster as Sonner } from "vue-sonner"
import { cn } from "@/lib/utils"

const props = defineProps<ToasterProps>()
</script>

<template>
  <Sonner
    :class="cn('toaster group', props.class)"
    :style="{
      '--normal-bg': 'var(--popover)',
      '--normal-text': 'var(--popover-foreground)',
      '--normal-border': 'var(--border)',
      '--border-radius': 'var(--radius)',
    }"
    v-bind="props"
  >
	<template #success-icon>
	  <Icon name="lucide:circle-check" class="size-4" />
	</template>
	<template #info-icon>
	  <Icon name="lucide:info" class="size-4" />
	</template>
	<template #warning-icon>
	  <Icon name="lucide:triangle-alert" class="size-4" />
	</template>
	<template #error-icon>
	  <Icon name="lucide:octagon-x" class="size-4" />
	</template>
	<template #loading-icon>
	  <div>
		<Icon name="lucide:loader-2" class="size-4 animate-spin" />
	  </div>
	</template>
	<template #close-icon>
	  <Icon name="lucide:x" class="size-4" />
	</template>
  </Sonner>
</template>
