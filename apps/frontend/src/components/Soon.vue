<script lang="ts" setup>
import MyBtn from '@elements/MyBtn.vue';

import Back from '../assets/buttons/back.svg?component';

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
        <Back />
      </MyBtn>
    </div>
  </div>
</template>