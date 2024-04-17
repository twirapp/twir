<script setup lang="ts">
import { IconLogout, IconSettings, IconUser } from '@tabler/icons-vue';
import { NButton } from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { RouterLink } from 'vue-router';

import { useLogout, useProfile } from '@/api';
import { storeToRefs } from 'pinia';

const { t } = useI18n();

const logout = useLogout();

const profile = storeToRefs(useProfile());
const isAdmin = computed(() => profile.data.value?.isBotAdmin);

async function callLogout() {
	await logout.mutateAsync();
	window.location.replace('/');
}
</script>

<template>
	<div class="flex flex-col gap-1 p-2">
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

		<router-link v-if="isAdmin" :to="{ name: 'AdminPanel' }" #="{ navigate, href }" custom>
			<n-button
				:href="href"
				secondary
				block
				type="info"
				tag="a"
				@click="navigate"
			>
				<template #icon>
					<IconUser />
				</template>

				Admin Panel
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
