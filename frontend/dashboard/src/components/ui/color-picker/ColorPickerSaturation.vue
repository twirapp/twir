<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from 'vue'

const props = defineProps<{
  hex: string
  saturation: number
  brightness: number
  hue: number
}>()

const emit = defineEmits<{
  update: [h: number, s: number, v: number]
}>()

const saturationRef = ref<HTMLDivElement | null>(null)
let isDragging = false
let containerRect: DOMRect | null = null

// Независимые координаты курсора в пикселях
const cursorX = ref(0)
const cursorY = ref(0)

// Инициализируем позицию курсора при монтировании
onMounted(() => {
  if (saturationRef.value) {
    const rect = saturationRef.value.getBoundingClientRect()
    cursorX.value = props.saturation * rect.width
    cursorY.value = (1 - props.brightness) * rect.height
  }
})

function startSaturationDrag(e: MouseEvent): void {
  e.preventDefault()
  isDragging = true
  containerRect = saturationRef.value?.getBoundingClientRect() || null
  if (!containerRect) return

  document.addEventListener('mousemove', onMouseMove, { passive: false })
  document.addEventListener('mouseup', onMouseUp, { once: true })
  onMouseMove(e)
}

function onMouseMove(e: MouseEvent): void {
  if (!isDragging || !containerRect) return
  e.preventDefault()

  // Вычисляем координаты мыши в пикселях
  const x = Math.min(Math.max(0, e.clientX - containerRect.left), containerRect.width)
  const y = Math.min(Math.max(0, e.clientY - containerRect.top), containerRect.height)

  // Сохраняем ТОЧНЫЕ пиксельные координаты для курсора
  cursorX.value = x
  cursorY.value = y

  // Вычисляем saturation и brightness для эмита, но НЕ используем их для позиционирования
  const newSaturation = x / containerRect.width
  const newBrightness = 1 - y / containerRect.height

  emit('update', props.hue, newSaturation, newBrightness)
}

function onMouseUp(): void {
  isDragging = false
  containerRect = null
  document.removeEventListener('mousemove', onMouseMove)
}

function updateRect(): void {
  if (isDragging && saturationRef.value) {
    containerRect = saturationRef.value.getBoundingClientRect()
  }
}

onMounted(() => {
  window.addEventListener('resize', updateRect)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateRect)
  document.removeEventListener('mousemove', onMouseMove)
})
</script>

<template>
  <div
    ref="saturationRef"
    :style="{ backgroundColor: `hsl(${hue}, 100%, 50%)` }"
    class="relative w-full h-56 cursor-crosshair rounded-sm select-none touch-none"
    @mousedown="startSaturationDrag"
  >
    <!-- Курсор теперь позиционируется по точным пиксельным координатам -->
    <div
      :style="{
        left: cursorX + 'px',
        top: cursorY + 'px',
      }"
      class="absolute size-3.5 border-2 border-white rounded-full z-10 shadow-[0_0_6px_rgba(0,0,0,0.2)] -translate-x-1/2 -translate-y-1/2 pointer-events-none"
    />
    <div
      class="absolute inset-0 bg-gradient-to-r from-white to-transparent rounded-[inherit] pointer-events-none"
    />
    <div
      class="absolute inset-0 bg-gradient-to-t from-black to-transparent rounded-[inherit] pointer-events-none"
    />
  </div>
</template>
