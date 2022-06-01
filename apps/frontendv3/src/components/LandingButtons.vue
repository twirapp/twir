<script lang="ts" setup>
  import { useStore } from '@nanostores/vue';
  import { defineProps } from 'vue';

  import { userStore } from '@/stores/userStore';

  const user = useStore(userStore); 

  interface Props {
    size?: 'large' | 'small'
  }

  const props = withDefaults(defineProps<Props>(), {
    size: 'small',
  });

  function redirectToLogin() {
    window.location.replace(`/api/auth?state=${window.btoa(window.location.origin)}/login`);
  }
</script>

<template>
  <div>
    <router-link
      v-if="user"
      type="button"
      to="/dashboard"
      class="inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
      :class="[
        props.size === 'small' ? 'px-6 py-2.5' : 'px-7 py-3'
      ]"
    >
      Dashboard
    </router-link>
    <button
      v-else
      type="button"
      class="inline-block border-2 border-purple-600 text-white font-medium text-xs leading-tight uppercase rounded hover:bg-black hover:bg-opacity-5 focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
      :class="[
        props.size === 'small' ? 'px-6 py-2.5' : 'px-7 py-3'
      ]"
      @click="redirectToLogin"
    >
      Login
    </button>
  </div>
</template>