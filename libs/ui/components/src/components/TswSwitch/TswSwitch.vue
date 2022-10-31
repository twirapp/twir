<template>
  <button
    role="switch"
    type="button"
    :aria-checked="isCheckedValue"
    class="tsw-switch"
    @click="toggle"
  ></button>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

const props =
  defineProps<{
    isChecked: boolean;
  }>();

const emit =
  defineEmits<{
    (event: 'update:isChecked', isChecked: boolean): void;
  }>();

const isCheckedValue = computed({
  get() {
    return props.isChecked;
  },
  set(isChecked: boolean) {
    emit('update:isChecked', isChecked);
  },
});

const toggle = () => {
  emit('update:isChecked', !isCheckedValue.value);
};
</script>

<style lang="postcss">
.tsw-switch {
  @apply w-11 h-6 border-2 bg-gray-30 rounded-full border-gray-30 relative;

  transition-property: border-color, background-color;
  transition-timing-function: theme('transitionTimingFunction.DEFAULT');
  transition-duration: 0.15s;

  &::before {
    content: '';
    @apply inline-block w-5 h-5 rounded-full bg-white-100 -translate-x-[10px] transition-transform;

    box-shadow: 0px 1px 3px rgba(16, 24, 40, 0.1), 0px 1px 2px rgba(16, 24, 40, 0.06);
  }

  &[aria-checked='true'] {
    @apply bg-purple-60 border-purple-60;

    &::before {
      @apply translate-x-[10px];
    }
  }
}
</style>
