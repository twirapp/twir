<script setup lang="ts">
import { useTimeout } from '@vueuse/core'
import { Check, Copy } from 'lucide-vue-next'

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
			<Check
				v-if="isPending"
				class="h-4 w-4"
			/>
			<Copy
				v-else
				class="h-4 w-4"
			/>
			<span class="sr-only">Copy</span>
		</Button>
	</div>
</template>
