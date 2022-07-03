<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { defineProps } from 'vue';
import { RouterLink } from 'vue-router';


import { redirectToLogin } from '@/functions/redirectToLogin';
import { api } from '@/plugins/api';
import { setUser, userStore } from '@/stores/userStore';

const user = useStore(userStore); 

interface Props {
  size?: 'large' | 'small'
  type?: 'col' | 'normal'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'small',
  type: 'normal',
});

async function logOut() {
  await api.post('/auth/logout');
  localStorage.clear();
  setUser(null);
}
</script>

<template>
  <div
    class="select-none space-x-2"
    :class="[
      props.type === 'col' ? 'space-x-0 space-y-2 flex flex-col w-full' : ''
    ]"
  >
    <RouterLink
      v-if="user"
      type="button"
      to="/dashboard"
      class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight rounded shadow text-white text-xs transition uppercase"
      :class="[
        props.size === 'small' ? 'px-6 py-3' : 'px-7 py-3'
      ]"
    >
      Dashboard
    </RouterLink>

    <button
      v-else
      type="button"
      class="border-2 border-[#9146FF] duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-[#9146FF] hover:bg-opacity-5 inline-block leading-tight rounded text-white text-xs transition uppercase"
      :class="[
        props.size === 'small' ? 'px-6 py-2.5' : 'px-7 py-2.5'
      ]"
      @click="redirectToLogin"
    >
      Login
    </button>
    
    <button
      v-if="user"
      type="button"
      class="border-2 border-red-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-opacity-5 hover:bg-red-200 inline-block leading-tight rounded text-white text-xs transition uppercase"
      :class="[
        props.size === 'small' ? 'px-4 py-2.5' : 'px-7 py-2.5',
        props.type === 'col' ? 'w-full' : ''
      ]"
      @click="logOut"
    >
      Logout
    </button>
  </div>
</template>