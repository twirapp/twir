<script lang="ts" setup>
import { type MenuOption, NMenu, NAvatar, NSpin, NInput } from 'naive-ui';
import { computed, h, ref, watch } from 'vue';

import { useProfile, useTwitchGetUsers, useDashboards, useSetDashboard } from '@/api/index.js';

const emits = defineEmits<{
	dashboardSelected: []
}>();

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
	emits('dashboardSelected');
});

const filterValue = ref('');

const menuOptions = computed<MenuOption[]>(() => {
	return usersForSelect.data.value?.users
	.filter(u => {
		return u.displayName.includes(filterValue.value) || u.login.includes(filterValue.value);
	})
	.map((u) => {
		return {
			key: u.id,
			label: u.login === u.displayName.toLocaleLowerCase()
				? u.displayName
				: `${u.displayName} (${u.login})`,
			icon: () => h(NAvatar, { src: u.profileImageUrl, round: true, size: 'small' }),
		};
	}) ?? [];
});

</script>

<template>
	<n-spin v-if="isProfileLoading || isDashboardsLoading"></n-spin>
	<div v-else>
		<n-input v-model:value="filterValue" placeholder="Search" />
		<n-menu
			v-model:value="activeDashboard"
			:collapsed-width="64"
			:collapsed-icon-size="22"
			:options="menuOptions"
			:icon-size="35"
		/>
	</div>
</template>
