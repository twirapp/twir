<script setup lang="ts">
import { IconLogout } from '@tabler/icons-vue';
import { NButton } from 'naive-ui';
import { useI18n } from 'vue-i18n';

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
	padding: 8px;
}
</style>
