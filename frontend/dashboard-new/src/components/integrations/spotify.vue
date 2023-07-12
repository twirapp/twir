<script setup lang='ts'>
import { useSpotifyIntegration } from '@/api/index.js';
import IconSpotify from '@/assets/icons/integrations/spotify.svg?component';
import OauthComponent from '@/components/integrations/oauthComponent.vue';

const manager = useSpotifyIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const getAuthLink = manager.useAuthLink();
</script>

<template>
  <oauth-component :data="data" :logout="logout.mutateAsync" :getLoginLink="getAuthLink">
    <template #icon>
      <IconSpotify style="width: 30px;" />
    </template>
  </oauth-component>
</template>
