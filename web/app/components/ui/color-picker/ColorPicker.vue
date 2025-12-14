<script lang="ts" setup>
import tinycolor from 'tinycolor2'
import { type HTMLAttributes, onMounted, reactive, ref, watch } from 'vue'

import {
	ColorPickerAlpha,
	ColorPickerEyeDropper,
	ColorPickerHue,
	ColorPickerInputs,
	ColorPickerPresets,
	ColorPickerSaturation,
} from '.'

import { cn } from '@/lib/utils'
import { Popover, PopoverContent, PopoverTrigger } from '~/components/ui/popover'

interface RGB {
	r: number
	g: number
	b: number
}

interface RGBA {
	r: number
	g: number
	b: number
	a: number
}

interface Props {
	modelValue?: string
	class?: HTMLAttributes['class']
	showPresets?: boolean
	showPipette?: boolean
	showCopy?: boolean
	outputFormat?: 'hex' | 'rgb' | 'hsl' | 'auto'
}

const props = withDefaults(defineProps<Props>(), {
	modelValue: 'rgba(255, 255, 255, 1)',
	showPresets: true,
	showPipette: true,
	showCopy: true,
	outputFormat: 'auto',
})

const emit = defineEmits<{
	'update:modelValue': [value: string]
}>()

const colorPickerState = ref<boolean>(false)

function openColorPickerPopover() {
	colorPickerState.value = true
}

function closeColorPickerPopover() {
	colorPickerState.value = false
}

const hue = ref<number>(0)
const saturation = ref<number>(1)
const brightness = ref<number>(1)
const alpha = ref<number>(100)
const rgb = reactive<RGB>({ r: 255, g: 255, b: 255 })
const rgba = reactive<RGBA>({ r: 255, g: 255, b: 255, a: 1 })
const hex = ref<string>('#ffffff')

const selectedColorModel = ref<'hex' | 'rgb' | 'hsl'>('hex')

function updateColor(): void {
	const color = tinycolor({ h: hue.value, s: saturation.value, v: brightness.value })
	const { r, g, b } = color.toRgb()

	rgb.r = r
	rgb.g = g
	rgb.b = b
	rgba.r = r
	rgba.g = g
	rgba.b = b
	rgba.a = alpha.value / 100
	hex.value = color.toHexString()

	const outputColor = getOutputColor(color)
	emit('update:modelValue', outputColor)
}

function getOutputColor(color: tinycolor.Instance): string {
	const colorWithAlpha = color.setAlpha(alpha.value / 100)

	if (props.outputFormat !== 'auto') {
		switch (props.outputFormat) {
			case 'hex':
				return alpha.value < 100 ? colorWithAlpha.toHex8String() : colorWithAlpha.toHexString()
			case 'rgb':
				return colorWithAlpha.toRgbString()
			case 'hsl':
				return colorWithAlpha.toHslString()
		}
	}

	switch (selectedColorModel.value) {
		case 'hex':
			return alpha.value < 100 ? colorWithAlpha.toHex8String() : colorWithAlpha.toHexString()
		case 'rgb':
			return colorWithAlpha.toRgbString()
		case 'hsl':
			return colorWithAlpha.toHslString()
		default:
			return colorWithAlpha.toRgbString()
	}
}

function updateFromHSV(h: number, s: number, v: number): void {
	hue.value = h
	saturation.value = s
	brightness.value = v
	updateColor()
}

function updateAlpha(newAlpha: number): void {
	alpha.value = newAlpha
	updateColor()
}

function updateColorFromHex(): void {
	const color = tinycolor(hex.value)
	if (color.isValid()) {
		const hsv = color.toHsv()
		updateFromHSV(hsv.h, hsv.s, hsv.v)
	}
}

function updateColorFromRGB(): void {
	const color = tinycolor({ r: rgb.r, g: rgb.g, b: rgb.b })
	if (color.isValid()) {
		const hsv = color.toHsv()
		updateFromHSV(hsv.h, hsv.s, hsv.v)
	}
}

function setPresetColor(color: string): void {
	const selectedColor = tinycolor(color)
	const hsv = selectedColor.toHsv()
	alpha.value = Math.round(selectedColor.getAlpha() * 100)
	updateFromHSV(hsv.h, hsv.s, hsv.v)
}

function handleColorPick(pickedColor: string): void {
	const color = tinycolor(pickedColor)
	if (color.isValid()) {
		const hsv = color.toHsv()
		const rgbValues = color.toRgb()
		hue.value = hsv.h || 0
		saturation.value = hsv.s
		brightness.value = hsv.v
		alpha.value = Math.round((rgbValues.a || 1) * 100)
		updateColor()
	}
}

function handleColorModelChange(newModel: 'hex' | 'rgb' | 'hsl'): void {
	selectedColorModel.value = newModel
	updateColor()
}

watch(
	() => props.modelValue,
	(newValue) => {
		if (!newValue) return
		const color = tinycolor(newValue)
		if (color.isValid()) {
			const hsv = color.toHsv()
			const rgbValues = color.toRgb()
			rgb.r = rgbValues.r
			rgb.g = rgbValues.g
			rgb.b = rgbValues.b
			rgba.r = rgbValues.r
			rgba.g = rgbValues.g
			rgba.b = rgbValues.b
			rgba.a = rgbValues.a || 1
			hue.value = hsv.h || 0
			saturation.value = hsv.s
			brightness.value = hsv.v
			alpha.value = Math.round((rgbValues.a || 1) * 100)
			hex.value = color.toHexString()
		}
	},
	{ immediate: true }
)

onMounted(() => {
	if (props.modelValue) {
		const color = tinycolor(props.modelValue)
		if (color.isValid()) {
			const hsv = color.toHsv()
			const rgbValues = color.toRgb()
			hue.value = hsv.h || 0
			saturation.value = hsv.s
			brightness.value = hsv.v
			alpha.value = Math.round((rgbValues.a || 1) * 100)
			updateColor()
		}
	}
})
</script>

<template>
	<Popover
		:modal="colorPickerState"
		:open="colorPickerState"
		@update:open="(v) => (colorPickerState = v)"
	>
		<PopoverTrigger
			@click="colorPickerState = true"
			:style="{ backgroundColor: `rgba(${rgba.r}, ${rgba.g}, ${rgba.b}, ${rgba.a})` }"
			:class="cn('w-5 h-5 cursor-pointer rounded-sm border ', props.class)"
		/>
		<PopoverContent align="start" class="w-80 overflow-hidden">
			<div class="grid gap-3">
				<ColorPickerSaturation
					:hex="hex"
					:saturation="saturation"
					:brightness="brightness"
					:hue="hue"
					@update="updateFromHSV"
				/>
				<div class="inline-flex items-center gap-3">
					<div class="flex-1 grid w-full gap-1">
						<ColorPickerHue
							:hue="hue"
							@update:hue="(h) => updateFromHSV(h, saturation, brightness)"
						/>
						<ColorPickerAlpha :alpha="alpha" :hex="hex" @update:alpha="updateAlpha" />
					</div>
					<ColorPickerEyeDropper
						v-if="props.showPipette"
						class="w-8 h-8 bg-transparent"
						@pick-color="handleColorPick"
					/>
				</div>
				<ColorPickerInputs
					:hex="hex"
					:rgb="rgb"
					:alpha="alpha"
					:show-copy="props.showCopy"
					@update:hex="
						(newHex) => {
							hex = newHex
							updateColorFromHex()
						}
					"
					@update:rgb="updateColorFromRGB"
					@update:alpha="updateAlpha"
					@update:color-model="handleColorModelChange"
				/>
				<ColorPickerPresets
					v-if="props.showPresets"
					:selectedColor="hex"
					@select="setPresetColor"
				/>
			</div>
		</PopoverContent>
	</Popover>
</template>
