<script setup lang="ts">
import { IconSettings, IconUser } from '@tabler/icons-vue';
import { NButton } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';
import { RouterLink } from 'vue-router';

import { useProfile } from '@/api';
import ButtonLogout from '@/layout/buttons/buttonLogout.vue';

const profile = storeToRefs(useProfile());
const isAdmin = computed(() => profile.data.value?.isBotAdmin);
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

		<button-logout />
	</div>
</template>
