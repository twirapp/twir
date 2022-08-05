<template>
  <button
    type="button"
    class="text-blue-600"
  >
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
