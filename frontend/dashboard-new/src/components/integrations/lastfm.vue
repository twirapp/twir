<script setup lang='ts'>
import { useLastfmIntegration } from '@/api/index.js';
import IconLastFM from '@/assets/icons/integrations/lastfm.svg?component';
import OauthComponent from '@/components/integrations/oauthComponent.vue';

const manager = useLastfmIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const getAuthLink = manager.useAuthLink();
</script>

<template>
  <oauth-component :data="data" :logout="logout.mutateAsync" :getLoginLink="getAuthLink.refetch">
    <template #icon>
      <IconLastFM style="width: 30px;" />
    </template>
  </oauth-component>
</template>
