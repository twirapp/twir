<template>
  <TswLoader v-if="isLoading" size="lg" />
  <button v-else-if="isError" class="login-btn" @click="redirectToLogin">
    {{ t('buttons.login') }}
  </button>
  <a v-else-if="user" class="inline-grid grid-flow-col gap-x-3 items-center" href="/app/dashboard">
    <div
      class="h-9 w-9 bg-contain rounded-full"
      :style="{
        backgroundImage: cssURL(user.profile_image_url),
      }"
    ></div>
  </a>
</template>

<script lang="ts" setup>
import { TswLoader } from '@tsuwari/ui-components';

import useTranslation from '@/hooks/useTranslation.js';
import { redirectToLogin, useUserProfile } from '@/services/auth';
import { cssURL } from '@/utils/css.js';

const t = useTranslation<'landing'>();

const { data: user, isError, isLoading } = useUserProfile();
</script>
