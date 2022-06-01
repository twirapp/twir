<script setup lang="ts">
import { useTitle } from '@vueuse/core';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { fetchAndSetUser } from '@/functions/fetchAndSetUser';


const router = useRouter();
const error = ref('');

onMounted(async () => {
  const title = useTitle();
  title.value = 'Tsuwari - Login';

  const params = new URLSearchParams(window.location.search.substring(1));

  const code = params.get('code');
  if (code) {
    const request = await fetch(
      '/api/auth/token?' +
        new URLSearchParams({
          code,
          state: window.btoa(window.location.origin + '/login'),
        }),
    );

    const response = await request.json();

    if (!request.ok) {
      error.value = response.message ?? 'Unknown error';
    } else {
      localStorage.setItem('accessToken', response.accessToken);
      localStorage.setItem('refreshToken', response.refreshToken);
      fetchAndSetUser();
      router.push('/dashboard');
    }
  }
});
</script>

<template>
  <div class="q-pa-md q-gutter-sm">
    <q-banner
      v-if="error"
      inline-actions
      rounded
      class="bg-orange text-white"
    >
      {{ error }}

      <template #action>
        <q-btn
          flat
          label="Go Back"
          to="/"
        />
      </template>
    </q-banner>
  </div>
</template>