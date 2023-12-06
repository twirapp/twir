<script setup lang="ts">
import { IconLogout } from '@tabler/icons-vue';
import { DropdownOption, NAvatar, NButton, NDropdown, NSpin } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useLogout, useProfile } from '@/api';
import { renderIcon } from '@/helpers';

const { t } = useI18n();
const { data: profileData, isLoading: isProfileLoading } = useProfile();
const logout = useLogout();

const profileOptions: DropdownOption[] = [{
	label: t('navbar.logout'),
	key: 'logout',
	icon: renderIcon(IconLogout),
	props: {
		onClick: () => {
			logout.mutate();
		},
	},
}];
</script>

<template>
	<n-dropdown trigger="click" :options="profileOptions" size="large">
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
