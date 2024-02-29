<script setup lang="ts">
import { type DropdownOption, NAvatar, NButton, NDropdown, NSpin } from 'naive-ui';
import { h } from 'vue';

import DropdownFooter from './profile/footer.vue';
import DropdownHeader from './profile/header.vue';

import { useProfile } from '@/api';

const { data: profileData, isLoading: isProfileLoading } = useProfile();

const profileOptions: DropdownOption[] = [
	{
		key: 'header',
		type: 'render',
		render: () => h(DropdownHeader),
	},
	{
		key: 'header-divider',
		type: 'divider',
	},
	{
		key: 'footer',
		type: 'render',
		render: () => h(DropdownFooter),
	},
];
</script>

<template>
	<n-dropdown trigger="click" :options="profileOptions" size="large" style="width: 400px;">
		<n-button text>
			<n-spin v-if="isProfileLoading" size="small" />
			<div v-else class="profile">
				<n-avatar
					size="small"
					:src="profileData?.avatar"
					round
				/>
			</div>
		</n-button>
	</n-dropdown>
</template>

<style scoped>
.profile {
	display: flex;
	gap: 5px;
	align-items: center;
}
</style>
