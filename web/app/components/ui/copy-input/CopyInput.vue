<script setup lang="ts">
import { useTimeout } from '@vueuse/core'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

const props = defineProps<{
	text: string
	type?: 'password' | 'text'
	class?: string
}>()

const { start: copyStart, isPending } = useTimeout(2000, { controls: true, immediate: false })

async function copyToClipboard() {
	await navigator.clipboard.writeText(props.text)
	copyStart()
}
</script>

<template>
	<div :class="cn('relative w-full', props.class)">
		<Input
			:type="type"
			:default-value="text"
			readonly
			class="font-mono pr-[42px]"
		/>
		<Button
			variant="ghost"
			size="icon"
			:disabled="!text"
			class="absolute right-[4px] top-[3px] h-[32px] w-[32px]"
			@click="copyToClipboard"
		>
			<Icon
				v-if="isPending"
				name="lucide:check"
				class="h-4 w-4"
			/>
			<Icon
				v-else
				name="lucide:copy"
				class="h-4 w-4"
			/>
			<span class="sr-only">Copy</span>
		</Button>
	</div>
</template>
