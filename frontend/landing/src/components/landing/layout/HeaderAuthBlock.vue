<template>
  <div v-if="isLoading" class="w-9 h-9 inline-flex items-center justify-center">
    <TswLoader size="md" />
  </div>
  <slot v-else-if="isError" name="error"></slot>
  <TswAvatar v-else-if="user" :href="'/dashboard'" :src="user.profile_image_url" :lazy="true" />
</template>

<script lang="ts" setup>
import { TswLoader, TswAvatar } from '@twir/ui-components';

import { useUserProfile } from '@/services/auth';

const { data: user, isError, isLoading } = useUserProfile();
</script>
