<template>
  <button
    :class="{
      'burger-menu-button': true,
      active: buttonState,
    }"
    @click.prevent="toggle"
  >
    <div></div>
  </button>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

const props =
  defineProps<{
    state: boolean;
  }>();

const emit =
  defineEmits<{
    (event: 'update:state', state: boolean): void;
  }>();

const buttonState = computed<boolean>({
  get() {
    return props.state;
  },
  set(value) {
    emit('update:state', value);
  },
});

const toggle = () => {
  buttonState.value = !buttonState.value;
};
</script>

<style lang="postcss">
.burger-menu-button {
  width: 38px;
  height: 38px;
  display: flex;
  justify-content: center;
  align-items: center;
  border: none;
  user-select: none;

  @apply min-lg:hidden;

  div {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    pointer-events: none;

    &::before,
    &::after {
      content: '';
      display: block;
      height: 1.4px;
      width: 22px;
      border-radius: 2px;
      background-color: theme('colors.gray.70');
      transition: transform 0.15s ease;
    }

    &::before {
      transform: translateY(-4px) rotate(0deg);
    }

    &::after {
      transform: translateY(4px) rotate(0deg);
    }
  }

  &.active div {
    &::before {
      transform: translateY(1px) rotate(45deg);
    }

    &::after {
      transform: translateY(0) rotate(-45deg);
    }
  }
}
</style>
