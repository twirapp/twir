<script lang="ts" setup>
import { ref } from 'vue'

const props = defineProps<{
  hue: number
}>()

const emit = defineEmits<{
  'update:hue': [value: number]
}>()

const hueRef = ref<HTMLDivElement | null>(null)

function startHueDrag(e: MouseEvent): void {
  e.preventDefault()

  const onMouseMove = (e: MouseEvent) => {
    if (!hueRef.value) return
    const rect = hueRef.value.getBoundingClientRect()
    const x = Math.min(Math.max(0, e.clientX - rect.left), rect.width)
    const newHue = (x / rect.width) * 360
    emit('update:hue', newHue)
  }

  const onMouseUp = () => {
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)

  onMouseMove(e)
}
</script>

<template>
  <div class="relative">
    <div
      ref="hueRef"
      class="relative z-10 w-full h-2 rgb-gradient cursor-pointer rounded-full"
      @mousedown="startHueDrag"
    >
      <div
        :style="{ left: `calc(${(hue / 360) * 100}% - 7px)` }"
        class="absolute w-[14px] h-[14px] bg-white rounded-full border-[3px] border-white cursor-pointer z-[999] top-1/2 -translate-y-1/2 pointer-events-none"
      />
    </div>
  </div>
</template>

<style scoped>
.rgb-gradient {
  background: linear-gradient(to right, red, yellow, lime, aqua, blue, magenta, red);
}
</style>
