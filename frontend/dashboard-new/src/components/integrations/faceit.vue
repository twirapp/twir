<script setup lang='ts'>
import { useFaceitIntegration } from '@/api/index.js';
import IconFaceit from '@/assets/icons/integrations/faceit.svg?component';
import OauthComponent from '@/components/integrations/oauthComponent.vue';

const manager = useFaceitIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const getAuthLink = manager.useAuthLink();
</script>

<template>
  <oauth-component name="Faceit" :data="data" :logout="logout.mutateAsync" :getLoginLink="getAuthLink.refetch">
    <template #icon>
      <IconFaceit style="width: 30px;" />
    </template>
  </oauth-component>
</template>
