<script setup lang='ts'>
import { useI18n } from 'vue-i18n';

import { useVKIntegration } from '@/api/index.js';
import IconVk from '@/assets/icons/integrations/vk.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useVKIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();

const { t } = useI18n();
</script>

<template>
	<oauth-component
		title="VK"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconVk"
		:description="t('integrations.songServicesInfo')"
	/>
</template>
