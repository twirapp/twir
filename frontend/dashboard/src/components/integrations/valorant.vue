<script setup lang="ts">
import { useThemeVars } from 'naive-ui';

import { useValorantIntegration } from '@/api/index.js';
import IconValorant from '@/assets/integrations/valorant.svg?use';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const themeVars = useThemeVars();
const manager = useValorantIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
</script>

<template>
	<oauth-component
		title="Valorant"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconValorant"
	>
		<template #description>
			<i18n-t
				keypath="integrations.valorant.info"
			>
				<b class="variable">$(valorant.profile.elo)</b>
				<b class="variable">$(valorant.profile.tier)</b>
			</i18n-t>
		</template>
	</oauth-component>
</template>

<style scoped>
.variable {
	color: v-bind('themeVars.successColor');
}
</style>
