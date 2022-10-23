<template>
  <TswLoader v-if="isFetching" size="lg" />
  <button v-else-if="user === null" class="login-btn" @click="redirectToLogin">
    {{ t('buttons.login') }}
  </button>
  <a v-else class="inline-grid grid-flow-col gap-x-3 items-center" href="/app/dashboard">
    <div
      class="h-9 w-9 bg-contain rounded-full"
      :style="{
        backgroundImage: cssURL(user.profile_image_url),
      }"
    ></div>
  </a>
</template>

<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import type { AuthUser } from '@tsuwari/shared';
import { TswLoader } from '@tsuwari/ui-components';
import { useFetch, isClient } from '@vueuse/core';

import useTranslation from '@/hooks/useTranslation.js';
import { redirectToLogin } from '@/services/auth.service.js';
import { userStore, setUser } from '@/stores/user.js';
import { authFetch } from '@/utils/authFetch.js';
import { cssURL } from '@/utils/css.js';

const user = useStore(userStore);

const t = useTranslation<'landing'>();

const { onFetchResponse, isFetching, execute, data } = useFetch('/api/auth/profile', {
  fetch: authFetch,
  immediate: false,
}).json<AuthUser>();

onFetchResponse(async () => {
  if (data.value) setUser(data.value);
});

if (isClient) {
  execute();
}
</script>
