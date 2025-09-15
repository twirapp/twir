<script lang="ts" setup>
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { CheckIcon, CopyIcon } from 'lucide-vue-next'
import tinycolor from 'tinycolor2'
import { computed, ref, watch } from 'vue'

const props = defineProps<{
	hex: string
	rgb: { r: number; g: number; b: number }
	alpha: number
	showCopy: boolean
}>()

const emit = defineEmits<{
	'update:hex': [value: string]
	'update:rgb': []
	'update:alpha': [value: number]
	'update:color-model': [model: 'hex' | 'rgb' | 'hsl']
}>()

const localHex = ref(props.hex)
const localRgb = ref({ ...props.rgb })
const localAlpha = ref(props.alpha)
const colorModel = ref<'hex' | 'rgb' | 'hsl'>('hex')
const isCopied = ref(false)

const localHsl = computed(() => {
	const color = tinycolor(localRgb.value)
	const hsl = color.toHsl()
	return {
		h: Math.round(hsl.h || 0),
		s: Math.round(hsl.s * 100),
		l: Math.round(hsl.l * 100),
	}
})

const editableHsl = ref({ h: 0, s: 0, l: 0 })

const colorString = computed(() => {
	const color = tinycolor(localRgb.value).setAlpha(localAlpha.value / 100)
	switch (colorModel.value) {
		case 'hex':
			return localAlpha.value < 100 ? color.toHex8String() : color.toHexString()
		case 'rgb':
			return color.toRgbString()
		case 'hsl':
			return color.toHslString()
		default:
			return color.toHexString()
	}
})

watch(
	() => props.hex,
	(v) => (localHex.value = v)
)
watch(
	() => props.rgb,
	(v) => {
		localRgb.value = { ...v }
		editableHsl.value = { ...localHsl.value }
	},
	{ deep: true }
)
watch(
	() => props.alpha,
	(v) => (localAlpha.value = v)
)
watch(colorModel, (newModel) => {
	if (newModel === 'hsl') {
		editableHsl.value = { ...localHsl.value }
	}
	emit('update:color-model', newModel)
})

const updateHex = () => emit('update:hex', localHex.value)
const updateRgb = () => emit('update:rgb')
const updateAlpha = () => emit('update:alpha', localAlpha.value)
const updateFromHsl = () => {
	const color = tinycolor({
		h: editableHsl.value.h,
		s: editableHsl.value.s / 100,
		l: editableHsl.value.l / 100,
	})
	const rgb = color.toRgb()
	localRgb.value = { r: rgb.r, g: rgb.g, b: rgb.b }
	localHex.value = color.toHexString()
	emit('update:hex', localHex.value)
	emit('update:rgb')
}

async function copyColor() {
	try {
		await navigator.clipboard.writeText(colorString.value)
		isCopied.value = true
		setTimeout(() => {
			isCopied.value = false
		}, 2000)
	} catch (error) {
		console.error('Ошибка копирования:', error)
		const textArea = document.createElement('textarea')
		textArea.value = colorString.value
		document.body.appendChild(textArea)
		textArea.select()
		document.execCommand('copy')
		document.body.removeChild(textArea)
		isCopied.value = true
		setTimeout(() => {
			isCopied.value = false
		}, 2000)
	}
}

const inputClass =
	'h-8 bg-transparent text-xs font-medium px-2.5 transition-colors hover:bg-white/5 [-moz-appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none'
</script>

<template>
	<div class="flex items-center gap-1 w-full overflow-hidden">
		<Select v-model="colorModel">
			<SelectTrigger size="sm" class="px-2 text-xs shrink-0">
				<SelectValue />
			</SelectTrigger>
			<SelectContent>
				<SelectItem value="hex" class="text-xs">HEX</SelectItem>
				<SelectItem value="rgb" class="text-xs">RGB</SelectItem>
				<SelectItem value="hsl" class="text-xs">HSL</SelectItem>
			</SelectContent>
		</Select>

		<!-- HEX -->
		<div v-if="colorModel === 'hex'" class="inline-flex -space-x-px min-w-0">
			<Input
				v-model="localHex"
				:class="[inputClass, 'uppercase rounded-r-none border-r-0 min-w-0']"
				placeholder="HEX"
				@blur="updateHex"
				@keydown.enter="updateHex"
			/>
			<Input
				v-model.number="localAlpha"
				type="number"
				min="0"
				max="100"
				placeholder="A"
				:class="[inputClass, 'rounded-l-none  min-w-0']"
				@blur="updateAlpha"
				@keydown.enter="updateAlpha"
			/>
		</div>

		<!-- RGB -->
		<div v-else-if="colorModel === 'rgb'" class="inline-flex -space-x-px min-w-0">
			<Input
				v-model.number="localRgb.r"
				type="number"
				min="0"
				max="255"
				placeholder="R"
				:class="[inputClass, 'rounded-r-none border-r-0 min-w-0']"
				@blur="updateRgb"
				@keydown.enter="updateRgb"
			/>
			<Input
				v-model.number="localRgb.g"
				type="number"
				min="0"
				max="255"
				placeholder="G"
				:class="[inputClass, 'rounded-none border-r-0 min-w-0']"
				@blur="updateRgb"
				@keydown.enter="updateRgb"
			/>
			<Input
				v-model.number="localRgb.b"
				type="number"
				min="0"
				max="255"
				placeholder="B"
				:class="[inputClass, 'rounded-none border-r-0 min-w-0']"
				@blur="updateRgb"
				@keydown.enter="updateRgb"
			/>
			<Input
				v-model.number="localAlpha"
				type="number"
				min="0"
				max="100"
				placeholder="A"
				:class="[inputClass, 'rounded-l-none min-w-0']"
				@blur="updateAlpha"
				@keydown.enter="updateAlpha"
			/>
		</div>

		<!-- HSL -->
		<div v-else-if="colorModel === 'hsl'" class="inline-flex -space-x-px min-w-0">
			<Input
				v-model.number="editableHsl.h"
				type="number"
				min="0"
				max="360"
				placeholder="H"
				:class="[inputClass, 'border-r-0 rounded-r-none min-w-0']"
				@blur="updateFromHsl"
				@keydown.enter="updateFromHsl"
			/>
			<Input
				v-model.number="editableHsl.s"
				type="number"
				min="0"
				max="100"
				placeholder="S"
				:class="[inputClass, 'rounded-l-none rounded-r-none border-r-0 min-w-0']"
				@blur="updateFromHsl"
				@keydown.enter="updateFromHsl"
			/>
			<Input
				v-model.number="editableHsl.l"
				type="number"
				min="0"
				max="100"
				placeholder="L"
				:class="[inputClass, 'rounded-l-none rounded-r-none border-r-0 min-w-0']"
				@blur="updateFromHsl"
				@keydown.enter="updateFromHsl"
			/>
			<Input
				v-model.number="localAlpha"
				type="number"
				min="0"
				max="100"
				placeholder="A"
				:class="[inputClass, 'rounded-l-none min-w-0']"
				@blur="updateAlpha"
				@keydown.enter="updateAlpha"
			/>
		</div>

		<Button
			v-if="props.showCopy"
			@click="copyColor"
			variant="outline"
			size="none"
			class="w-8 h-8 bg-transparent"
			:class="{ 'bg-green-50 border-green-200': isCopied }"
		>
			<CheckIcon v-if="isCopied" class="h-3 w-3 text-green-600" />
			<CopyIcon v-else class="h-3 w-3" />
		</Button>
	</div>
</template>
