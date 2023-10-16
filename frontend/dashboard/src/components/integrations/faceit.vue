<script setup lang='ts'>
import { useI18n } from 'vue-i18n';

import { useFaceitIntegration } from '@/api/index.js';
import IconFaceit from '@/assets/icons/integrations/faceit.svg?component';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useFaceitIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();

const { t } = useI18n();
</script>

<template>
	<oauth-component
		title="Faceit"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconFaceit"
		:description="t('integrations.faceit.info')"
	/>
</template>
