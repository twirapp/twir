<script setup lang="ts">
import { PanelLeft } from 'lucide-vue-next'

import { useSidebar } from './utils'

import type { HTMLAttributes } from 'vue'

import Button from '@/components/ui/button/Button.vue'
import { cn } from '@/lib/utils'

const props = defineProps<{
	class?: HTMLAttributes['class']
}>()

const { toggleSidebar } = useSidebar()
</script>

<template>
	<Button
		data-sidebar="trigger"
		variant="ghost"
		size="icon"
		:class="cn('size-[32px]', props.class)"
		@click="toggleSidebar"
	>
		<PanelLeft class="size-5" />
		<span class="sr-only">Toggle Sidebar</span>
	</Button>
</template>
