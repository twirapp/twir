<script lang="ts" setup>
import MyBtn from '@elements/MyBtn.vue';

interface Props {
  button?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  button: true,
});
</script>

<template>
  <div class="flex h-full select-none">
    <div class="m-auto space-y-2 text-center">
      <h1 class="font-bold text-6xl">
        Oops!
      </h1>
      <p
        class="text-lg"
      >
        This feature is coming soon!
      </p>

      <MyBtn
        v-if="props.button"
        size="default"
        color="purple"
        class="px-10"
        @click="$router.go(-1)"
      >
        <svg
          class="h-6 w-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1"
          d="M12.066 11.2a1 1 0 000 1.6l5.334 4A1 1 0 0019 16V8a1 1 0 00-1.6-.8l-5.333 4zM4.066 11.2a1 1 0 000 1.6l5.334 4A1 1 0 0011 16V8a1 1 0 00-1.6-.8l-5.334 4z"
        /></svg>
      </MyBtn>
    </div>
  </div>
</template>