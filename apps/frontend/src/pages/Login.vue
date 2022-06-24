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
  <div
    v-if="error"
    class="flex h-screen select-none"
  >
    <div class="m-auto text-center space-y-2 text-white">
      <h1 class="font-bold text-6xl">
        Oops!
      </h1>
      <p
        class="text-lg"
      >
        That seems like you denied us to use your account data.<br>Sorry, but we can't provide our service then.
      </p>

      <button
        class="inline-block px-10 py-2.5 bg-[#9146FF] text-white text-lg leading-tight uppercase rounded shadow hover:bg-[#772CE8]  focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
        @click="router.push('/')"
      >
        Home
      </button>
    </div>
  </div>
</template>