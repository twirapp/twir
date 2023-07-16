<script setup lang='ts'>
import { useDonationAlertsIntegration } from '@/api/index.js';
import IconDonationAlerts from '@/assets/icons/integrations/donationalerts.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useDonationAlertsIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const getAuthLink = manager.useAuthLink();
</script>

<template>
  <oauth-component name="DonationAlerts" :data="data" :logout="logout.mutateAsync" :getLoginLink="getAuthLink.refetch">
    <template #icon>
      <IconDonationAlerts style="width: 30px;" />
    </template>
  </oauth-component>
</template>
