<script setup lang="ts">
import { useTitle } from '@vueuse/core';
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';

import { fetchAndSetUser } from '@/functions/fetchAndSetUser';


const router = useRouter();

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
          state: window.btoa(window.location.origin),
        }),
    );

    if (!request.ok) {
      return router.push('/');
    }

    const response = await request.json();
    localStorage.setItem('accessToken', response.accessToken);
    localStorage.setItem('refreshToken', response.refreshToken);
    fetchAndSetUser();
    router.push('/dashboard');
  }
});
</script>
