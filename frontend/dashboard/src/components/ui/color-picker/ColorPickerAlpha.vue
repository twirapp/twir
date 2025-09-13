<script lang="ts" setup>
import { ref } from 'vue'

const props = defineProps<{
  alpha: number
  hex: string
}>()

const emit = defineEmits<{
  'update:alpha': [value: number]
}>()

const alphaRef = ref<HTMLDivElement | null>(null)

function startAlphaDrag(e: MouseEvent): void {
  e.preventDefault()

  const onMouseMove = (e: MouseEvent) => {
    if (!alphaRef.value) return
    const rect = alphaRef.value.getBoundingClientRect()
    const x = Math.min(Math.max(0, e.clientX - rect.left), rect.width)
    const newAlpha = Math.round((x / rect.width) * 100)
    emit('update:alpha', newAlpha)
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
      ref="alphaRef"
      class="alpha-checker relative w-full h-2.5 cursor-pointer rounded-full"
      @mousedown="startAlphaDrag"
    >
      <div class="absolute inset-0 rounded-full bg-gradient-to-r from-transparent to-black/50" />
      <div
        :style="{ left: `calc(${props.alpha}%  - 7px)` }"
        class="absolute w-[14px] h-[14px] bg-white rounded-full border-[3px] border-white cursor-pointer z-[999] top-1/2 -translate-y-1/2 pointer-events-none"
      />
    </div>
  </div>
</template>

<style scoped>
.alpha-checker {
  background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAwAAAAMCAIAAADZF8uwAAAAGUlEQVQYV2M4gwH+YwCGIasIUwhT25BVBADtzYNYrHvv4gAAAABJRU5ErkJggg==');
  background-repeat: repeat;
}
</style>
