<script setup lang="ts">
import { NFormItem, NSelect } from 'naive-ui';
import { ref } from 'vue';

import { useFaceitIntegration } from '@/api/index.js';
import IconFaceit from '@/assets/integrations/faceit.svg?use';
import OauthComponent from '@/components/integrations/variants/oauth.vue';

const manager = useFaceitIntegration();
const { data } = manager.useData();
const logout = manager.useLogout();
const { data: authLink } = manager.useAuthLink();
const updater = manager.update!();

const game = ref('');

async function save() {
	await updater.mutateAsync({ game: game.value });
}
</script>

<template>
	<oauth-component
		title="Faceit"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconFaceit"
		:withSettings="true"
		:save="save"
	>
		<template #description>
			<i18n-t
				keypath="integrations.faceit.info"
			>
				<b class="variable">$(faceit.elo)</b>
				<b class="variable">$(faceit.lvl)</b>
			</i18n-t>
		</template>
		<template #settings>
			<NFormItem label="Game">
				<n-select
					v-model:value="game"
					defaultValue="cs2"
					:options="[
						{ label: 'Counter-Strike', value: 'cs2' },
						{ label: 'Counter-Strike: Global Offensive', value: 'csgo' },
						{ label: 'Dota 2', values: 'dota2' }
					]"
				/>
			</NFormItem>
		</template>
	</oauth-component>
</template>

<style scoped>
.variable {
	@apply text-[color:var(--n-color-target)];
}
</style>
