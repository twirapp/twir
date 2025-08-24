<script lang="ts" setup>
import tinycolor from 'tinycolor2'
import { onMounted, reactive, ref, watch } from 'vue'

import { Input } from '@/components/ui/input'
import InputWithIcon from '@/components/ui/InputWithIcon.vue'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

interface RGB {
	r: number
	g: number
	b: number
}

const props = defineProps({
	modelValue: {
		type: String,
		default: '#ffffff',
		required: false,
	},
})

const emit = defineEmits(['update:modelValue'])

const hue = ref<number>(0)
const saturation = ref<number>(1)
const brightness = ref<number>(1)
const rgb = reactive<RGB>({ r: 255, g: 255, b: 255 })
const hex = ref<string>(props.modelValue)

const presetColors: string[] = [
	'#f2f2f2',
	'#e8c6d2',
	'#f2c1d1',
	'#f1b79b',
	'#ff9e6d',
	'#ff7986',
	'#f25c54',
	'#ff9f00',
	'#f4a300',
	'#09bc8a',
	'#43b2a1',
	'#4c8c8b',
	'#00ffe3',
	'#004f8c',
	'#a78ec1',
	'#b4a0d0',
	'#bb2649',
	'#07090e',
]

const saturationRef = ref<HTMLDivElement | null>(null)
const hueRef = ref<HTMLDivElement | null>(null)

function updateColor(): void {
	const color = tinycolor({ h: hue.value, s: saturation.value, v: brightness.value })
	const { r, g, b } = color.toRgb()
	rgb.r = r
	rgb.g = g
	rgb.b = b
	hex.value = color.toHexString()
	emit('update:modelValue', hex.value)
}

watch(
	() => props.modelValue,
	(newValue) => {
		const color = tinycolor(newValue)
		if (color.isValid()) {
			const hsv = color.toHsv()
			hue.value = hsv.h
			saturation.value = hsv.s
			brightness.value = hsv.v
			updateColor()
		}
	},
	{ immediate: true },
)

function updateColorFromHex(): void {
	const color = tinycolor(hex.value)
	if (color.isValid()) {
		const hsv = color.toHsv()
		hue.value = hsv.h
		saturation.value = hsv.s
		brightness.value = hsv.v
		updateColor()
	}
}

function updateColorFromRGB(): void {
	const color = tinycolor({ r: rgb.r, g: rgb.g, b: rgb.b })
	if (color.isValid()) {
		const hsv = color.toHsv()
		hue.value = hsv.h
		saturation.value = hsv.s
		brightness.value = hsv.v
		updateColor()
	}
}

function setPresetColor(color: string): void {
	const selectedColor = tinycolor(color)
	const hsv = selectedColor.toHsv()
	hue.value = hsv.h
	saturation.value = hsv.s
	brightness.value = hsv.v
	updateColor()
}

function startSaturationDrag(): void {
	const onMouseMove = (e: MouseEvent) => {
		if (!saturationRef.value) return
		const rect = saturationRef.value.getBoundingClientRect()
		const x = Math.min(Math.max(0, e.clientX - rect.left), rect.width)
		const y = Math.min(Math.max(0, e.clientY - rect.top), rect.height)
		saturation.value = x / rect.width
		brightness.value = 1 - y / rect.height
		updateColor()
	}

	const onMouseUp = () => {
		window.removeEventListener('mousemove', onMouseMove)
		window.removeEventListener('mouseup', onMouseUp)
	}

	window.addEventListener('mousemove', onMouseMove)
	window.addEventListener('mouseup', onMouseUp)
}

function startHueDrag(): void {
	const onMouseMove = (e: MouseEvent) => {
		if (!hueRef.value) return
		const rect = hueRef.value.getBoundingClientRect()
		const x = Math.min(Math.max(0, e.clientX - rect.left), rect.width)
		hue.value = (x / rect.width) * 360
		updateColor()
	}

	const onMouseUp = () => {
		window.removeEventListener('mousemove', onMouseMove)
		window.removeEventListener('mouseup', onMouseUp)
	}

	window.addEventListener('mousemove', onMouseMove)
	window.addEventListener('mouseup', onMouseUp)
}

onMounted(() => {
	const color = tinycolor(hex.value)
	if (color.isValid()) {
		const hsv = color.toHsv()
		hue.value = hsv.h
		saturation.value = hsv.s
		brightness.value = hsv.v
		updateColor()
	}
})
</script>

<template>
	<Popover>
		<InputWithIcon v-model="hex" class="relative uppercase" @blur="updateColorFromHex">
			<PopoverTrigger
				:style="{
					backgroundColor: hex,
				}"
				class="z-10 size-5 cursor-pointer rounded-sm"
			/>
		</InputWithIcon>
		<PopoverContent align="start">
			<div class="color-picker grid gap-3">
				<div
					ref="saturationRef"
					:style="{ backgroundColor: `hsl(${hue}, 100%, 50%)` }"
					class="saturation"
					@mousedown="startSaturationDrag"
				>
					<div
						:style="{
							left: `${saturation * 100}%`,
							top: `${(1 - brightness) * 100}%`,
							backgroundColor: hex,
						}"
						class="saturation-pointer"
					></div>
					<div class="saturation-white"></div>
					<div class="saturation-black"></div>
				</div>
				<div class="relative">
					<div ref="hueRef" class="hue relative z-10" @mousedown="startHueDrag">
						<div :style="{ left: `calc(${(hue / 360) * 100}% - 7px)` }" class="hue-pointer"></div>
					</div>
					<div class="hue absolute top-0 blur-[16px]"></div>
				</div>
				<div class="inputs">
					<Input
						id="hex"
						v-model="hex"
						class="h-8 bg-transparent text-xs font-medium uppercase transition-colors hover:bg-white/5"
						@focusout="updateColorFromHex"
					/>
					<Input
						id="r"
						v-model.number="rgb.r"
						class="h-8 bg-transparent text-xs font-medium transition-colors hover:bg-white/5"
						@focusout="updateColorFromRGB"
					/>
					<Input
						id="g"
						v-model.number="rgb.g"
						class="h-8 bg-transparent text-xs font-medium transition-colors hover:bg-white/5"
						@focusout="updateColorFromRGB"
					/>
					<Input
						id="b"
						v-model.number="rgb.b"
						class="h-8 bg-transparent text-xs font-medium transition-colors hover:bg-white/5"
						@focusout="updateColorFromRGB"
					/>
				</div>

				<div class="presets">
					<div
						v-for="color in presetColors"
						:key="color"
						:style="{ backgroundColor: color }"
						class="preset-color"
						:class="[hex === color ? 'active' : '']"
						@click="setPresetColor(color)"
					></div>
				</div>
			</div>
		</PopoverContent>
	</Popover>
</template>

<style lang="scss" scoped>
.color-picker {
  .saturation,
  .saturation-white,
  .saturation-black {
    position: relative;
    width: 100%;
    height: 150px;
    background: red;
    cursor: crosshair;
    border-radius: calc(var(--radius) - 4px);

    .saturation-white {
      position: absolute;
      background: linear-gradient(90deg, #fff, hsla(0, 0%, 100%, 0));
    }

    .saturation-black {
      position: absolute;
      background: linear-gradient(0deg, #000, transparent);
    }

    .saturation-pointer {
      position: absolute;
      width: 14px;
      height: 14px;
      border: 2px solid #fff;
      z-index: 999;
      box-shadow: 0 0 6px rgba(0, 0, 0, 0.2);
      border-radius: 50%;
      transform: translate(-50%, -50%);
    }
  }

  .hue {
    width: 100%;
    height: 10px;
    background: linear-gradient(to right, red, yellow, lime, aqua, blue, magenta, red);
    cursor: pointer;
    border-radius: 999px;

    .hue-pointer {
      position: absolute;
      width: 14px;
      height: 14px;
      background-color: transparent;
      border-radius: 999px;
      border: 3px solid #fff;
      cursor: pointer;
      z-index: 999;
      transform: translateY(-50%);
      top: 50%;
    }
  }

  .presets {
    display: inline-flex;
    flex-wrap: wrap;
    gap: 7px;

    .preset-color {
      width: 22px;
      height: 22px;
      box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.1);
      transition: 0.2s box-shadow;
      border-radius: calc(var(--radius) - 4px);
      cursor: pointer;

      &.active,
      &:hover {
        box-shadow: 0 0 0 2px rgba(255, 255, 255, 1);
      }
    }
  }

  .inputs {
    display: grid;
    gap: 4px;
    align-items: center;
    align-content: center;
    grid-template-columns: 2fr repeat(3, 1fr);
    grid-template-rows: 1fr;
  }
}
</style>
