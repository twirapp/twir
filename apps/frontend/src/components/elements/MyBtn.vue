<script lang="ts" setup>

type Props = {
  color: 'green' | 'purple' | 'red'
  disabled?: boolean
  upperCase?: boolean
  size?: 'default' | 'small',
  show?: boolean,
  type?: 'button' | 'submit'
}

withDefaults(defineProps<Props>(), {
  size: 'default',
  upperCase: true,
  show: true,
  type: 'button',
});

defineEmits(['click']);
</script>

<template>
  <button
    v-if="show"
    class="duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium inline-block leading-tight rounded shadow text-white text-xs transition"
    :type="type"
    :class="{
      'bg-green-600 hover:bg-green-700': color === 'green',
      'bg-red-600 hover:bg-red-700': color === 'red',
      'bg-purple-600 hover:bg-purple-700': color === 'purple',
      'uppercase': upperCase,
      'px-6 py-2.5': size === 'default',
      'px-1 py-1': size === 'small'
    }"
    @click="$emit('click', $event)"
  >
    <slot />
  </button>
</template>