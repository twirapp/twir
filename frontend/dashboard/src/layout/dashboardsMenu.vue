<script lang="ts" setup>
import { useMagicKeys } from '@vueuse/core';
import { type MenuOption, NMenu, NAvatar, NSpin } from 'naive-ui';
import { computed, h, ref, watch } from 'vue';

import { useProfile, useTwitchGetUsers, useDashboards, useSetDashboard } from '@/api/index.js';

const keys = useMagicKeys();
const CmdK = keys['Meta+K'];

watch(CmdK, (v) => {
	if (v) {
		console.log('Meta + K has been pressed');
	}
});

const { data: profile, isLoading: isProfileLoading } = useProfile();
const { data: dashboards, isLoading: isDashboardsLoading } = useDashboards();
const setDashboard = useSetDashboard();

const twitchUsersIds = computed<string[]>(() => {
	return [
		profile.value?.id,
		...(dashboards.value?.dashboards.map(d => d.id) ?? []),
	].filter(Boolean) as string[] ?? [];
});

const usersForSelect = useTwitchGetUsers({
	ids: twitchUsersIds,
});

const currentDashboard = computed(() => {
	const dashboard = usersForSelect.data.value?.users.find(u => u.id === profile.value?.selectedDashboardId);
	if (!dashboard) return null;

	return dashboard;
});

const activeDashboard = ref('');
watch(currentDashboard, (v) => {
	if (!v) return;
	activeDashboard.value = v.id;
}, { immediate: true });

watch(activeDashboard, async (v) => {
	if (v === profile.value?.selectedDashboardId) return;

	await setDashboard.mutateAsync(v);
});

const menuOptions = computed<MenuOption[]>(() => {
	return usersForSelect.data.value?.users.map((u) => {
		return {
			key: u.id,
			label: u.login === u.displayName.toLocaleLowerCase()
				? u.displayName
				: `${u.displayName} (${u.login})`,
			icon: () => h(NAvatar, { src: u.profileImageUrl, round: true, size: 'small' }),
		};
	}).filter(Boolean) as MenuOption[] ?? [];
});
</script>

<template>
	<n-spin v-if="isProfileLoading || isDashboardsLoading"></n-spin>
	<n-menu
		v-else
		v-model:value="activeDashboard"
		:collapsed-width="64"
		:collapsed-icon-size="22"
		:options="menuOptions"
		:icon-size="35"
	/>
</template>
