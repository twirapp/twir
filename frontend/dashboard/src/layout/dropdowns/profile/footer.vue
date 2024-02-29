<script setup lang="ts">
import { IconLogout, IconSettings } from '@tabler/icons-vue';
import { NButton } from 'naive-ui';
import { useI18n } from 'vue-i18n';
import { RouterLink } from 'vue-router';

import { useLogout } from '@/api';

const logout = useLogout();

async function callLogout() {
	await logout.mutateAsync();
	window.location.replace('/');
}

const { t } = useI18n();
</script>

<template>
	<div class="footer">
		<router-link :to="{ name: 'Settings' }" #="{ navigate, href }" custom>
			<n-button
				:href="href"
				secondary
				block
				type="info"
				tag="a"
				@click="navigate"
			>
				<template #icon>
					<IconSettings />
				</template>

				Settings
			</n-button>
		</router-link>

		<n-button
			secondary
			block
			type="error"
			:loading="logout.isLoading.value"
			@click="callLogout"
		>
			<template #icon>
				<IconLogout />
			</template>

			{{ t('navbar.logout') }}
		</n-button>
	</div>
</template>


<style scoped>
.footer {
	display: flex;
	flex-direction: column;
	gap: 4px;
	padding: 8px;
}
</style>
