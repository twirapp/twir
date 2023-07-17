<script setup lang='ts'>
import { useVKIntegration } from '@/api/index.js';
import IconVk from '@/assets/icons/integrations/vk.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useVKIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const getAuthLink = manager.useAuthLink();
</script>

<template>
  <oauth-component
    name="VK"
    :data="data"
    :logout="() => logout.mutateAsync({})"
    :getLoginLink="getAuthLink.refetch"
    :icon="IconVk"
  />
</template>
