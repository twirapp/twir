<script setup lang='ts'>
import { useStreamlabsIntegration } from '@/api/index.js';
import IconStreamLabs from '@/assets/icons/integrations/streamlabs.svg?component';
import OauthComponent from '@/components/integrations/oauthComponent.vue';

const manager = useStreamlabsIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const getAuthLink = manager.useAuthLink();
</script>

<template>
  <oauth-component name="StreamLabs" :data="data" :logout="logout.mutateAsync" :getLoginLink="getAuthLink.refetch">
    <template #icon>
      <IconStreamLabs style="width: 30px;" />
    </template>
  </oauth-component>
</template>
