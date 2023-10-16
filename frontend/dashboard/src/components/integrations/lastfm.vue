<script setup lang='ts'>
import { useI18n } from 'vue-i18n';

import { useLastfmIntegration } from '@/api/index.js';
import IconLastFM from '@/assets/icons/integrations/lastfm.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useLastfmIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();

const { t } = useI18n();
</script>

<template>
	<oauth-component
		title="Last.fm"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconLastFM"
		:description="t('integrations.songServicesInfo')"
	/>
</template>
