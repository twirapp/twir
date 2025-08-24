<script setup lang="ts">
import { useVModel } from '@vueuse/core'
import { Eye, EyeClosed } from 'lucide-vue-next'
import { ref } from 'vue'

import type { HTMLAttributes, InputHTMLAttributes } from 'vue'

import Button from '@/components/ui/button/Button.vue'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

interface Props {
	modelValue?: string | number
	placeholder?: string
	id?: string
	class?: HTMLAttributes['class']
	disabled?: InputHTMLAttributes['disabled']
	type?: InputHTMLAttributes['type']
}

const props = withDefaults(defineProps<Props>(), {
	type: 'text',
	placeholder: '',
	id: '',
})

const emits = defineEmits<{
	(e: 'update:modelValue', payload: string | number): void
}>()

const modelValue = useVModel(props, 'modelValue', emits, {
	passive: true,
})

// Создаем локальное состояние для типа ввода
const inputType = ref(props.type)

// Функция для переключения типа ввода
function toggleInputType() {
	inputType.value = inputType.value === 'password' ? 'text' : 'password'
}
</script>

<template>
	<div
		class="group relative w-full items-center overflow-hidden rounded-md border border-white/10 hover:border-white/15"
	>
		<span class="absolute inset-y-0 start-0 flex items-center justify-center pl-3">
			<slot />
		</span>
		<span class="absolute inset-y-0 start-0 flex items-center justify-center pl-3 blur-md">
			<slot />
		</span>
		<Input
			:id="props.id"
			v-model="modelValue"
			:type="inputType"
			:placeholder="props.placeholder"
			:disabled="props.disabled"
			:class="
				cn(
					'border-transparent bg-transparent pl-10 font-medium transition-colors hover:bg-white/5 focus-visible:bg-white/5',
					props.class,
				)
			"
		/>
		<div
			v-if="props.type === 'password'"
			class="show-password opacity-0 transition-all group-hover:opacity-100"
		>
			<Button
				class="absolute inset-y-0 end-0 flex h-full items-center justify-center rounded-none border-l border-white/15 backdrop-blur-sm"
				size="icon"
				variant="ghost"
				@click="toggleInputType"
			>
				<Eye v-if="inputType === 'password'" class="size-4" />
				<EyeClosed v-else class="size-4" />
			</Button>
		</div>
	</div>
</template>
