<template>
  <button type="button">
    {{ text }}
  </button>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'MyButton',
  props: {
    text: {
      type: String,
      default: 'Hello World',
    },
  },
});
</script>
