<script setup lang="ts">
import { useTitle } from '@vueuse/core';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const error = ref('');

onMounted(async () => {
  const title = useTitle();
  title.value = 'Tsuwari - Login';

  const params = new URLSearchParams(window.location.search.substring(1));

  const code = params.get('code');
  const isError = params.get('error');
  if (code) {
    const request = await fetch(
      '/api/auth/token?' +
        new URLSearchParams({
          code,
          state: window.btoa(window.location.origin + '/login'),
        }),
    );

    if (!request.ok) {
      return router.push('/');
    }

    const response = await request.json();
    console.log(response);
    localStorage.setItem('accessToken', response.accessToken);
    localStorage.setItem('refreshToken', response.refreshToken);

    router.push('/dashboard');
  }

  if (isError) {
    error.value = isError;
  }
});
</script>

<template>
  <div v-if="error" class="flex h-screen select-none">
    <div class="m-auto space-y-2 text-center text-white">
      <h1 class="font-bold text-6xl">Oops!</h1>
      <p class="text-lg">
        That seems like you denied us to use your account data.<br />Sorry, but we can't provide our
        service then.
      </p>

      <button
        class="bg-[#9146FF] duration-150 ease-in-out focus:outline-none focus:ring-0 hover:bg-[#772CE8] inline-block leading-tight px-10 py-2.5 rounded shadow text-lg text-white transition uppercase"
        @click="router.push('/')">
        Home
      </button>
    </div>
  </div>

  <div v-else class="flex h-screen select-none">
    <div class="m-auto">
      <div class="flex items-center text-white">
        <div
          class="animate-spin border-4 h-8 inline-block rounded-full spinner-border w-8"
          role="status" />
        <span class="ml-2 text-lg">Loading...</span>
      </div>
    </div>
  </div>
</template>
