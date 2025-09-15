<script lang="ts" setup>
const props = defineProps<{
  selectedColor: string
}>()

const emit = defineEmits<{
  select: [color: string]
}>()

const presetColors: string[] = [
  '#f2f2f2',
  '#f25c54',
  '#fa7268',
  '#f1b79b',
  '#ff9e6d',
  '#ff9f00',
  '#f4a300',
  '#09bc8a',
  '#43b2a1',
  '#00ffe3',
  '#4c8c8b',
  '#004f8c',
  '#1a60a8',
  '#07090e',
  '#b4a0d0',
  '#a78ec1',
  '#e8c6d2',
  '#f2c1d1',
  '#bb2649',
  '#ff7986',
]

function selectColor(color: string) {
  emit('select', color)
}
</script>

<template>
  <div class="presets">
    <div
      v-for="color in presetColors"
      :key="color"
      :style="{ backgroundColor: color }"
      class="preset-color"
      :class="{ active: selectedColor === color }"
      @click="selectColor(color)"
    />
  </div>
</template>

<style scoped>
.presets {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 7px;
}

.preset-color {
  width: 22px;
  height: 22px;
  box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.1);
  transition: 0.2s box-shadow;
  border-radius: calc(var(--radius) - 4px);
  cursor: pointer;
}

.preset-color.active,
.preset-color:hover {
  box-shadow: 0 0 0 2px rgba(255, 255, 255, 1);
}
</style>
