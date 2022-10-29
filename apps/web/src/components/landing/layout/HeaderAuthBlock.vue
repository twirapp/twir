<template>
  <TswLoader v-if="isLoading" size="lg" />
  <slot v-else-if="isError" name="error"></slot>
  <TswAvatar v-else-if="user" :href="'/app/dashboard'" :src="user.profile_image_url" :lazy="true" />
</template>

<script lang="ts" setup>
import { TswLoader, TswAvatar } from '@tsuwari/ui-components';

import { useUserProfile } from '@/services/auth';

const { data: user, isError, isLoading } = useUserProfile();
</script>
