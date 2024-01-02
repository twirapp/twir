<script lang="ts" setup>
import { IconChevronRight } from '@tabler/icons-vue';
import { onClickOutside, onKeyStroke } from '@vueuse/core';
import {
	NAvatar,
	NInput,
	NSpin,
	NVirtualList,
	useThemeVars,
	NText,
	NPopover,
} from 'naive-ui';
import { computed, ref, watch } from 'vue';

import { useDashboards, useProfile, useSetDashboard, useTwitchGetUsers } from '@/api/index.js';

const emits = defineEmits<{
	dashboardSelected: []
}>();

const themeVars = useThemeVars();
const blockColor = computed(() => themeVars.value.buttonColor2);
const blockColor2 = computed(() => themeVars.value.buttonColor2Hover);

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

const menuOptions = computed(() => {
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
				icon: u.profileImageUrl,
			};
		}) ?? [];
});

const isSelectDashboardPopoverOpened = ref(false);

function togglePopover(value?: boolean) {
	isSelectDashboardPopoverOpened.value = value ?? !isSelectDashboardPopoverOpened.value;
}

function onSelectDashboard(key: string) {
	activeDashboard.value = key;
	togglePopover(false);
}

onKeyStroke('k', (event) => {
	if (event.ctrlKey || event.metaKey) {
		event.preventDefault();
		togglePopover();
	}
});

const refPopoverList = ref<HTMLElement | null>();
const refPopover = ref<HTMLElement | null>();
onClickOutside(refPopover, (event) => {
	if (isSelectDashboardPopoverOpened.value) {
		event.stopPropagation();
		togglePopover(false);
	}
}, { ignore: [refPopoverList] });
</script>

<template>
	<n-popover
		ref="refPopover"
		placement="bottom-start"
		trigger="manual"
		:show="isSelectDashboardPopoverOpened"
		:show-arrow="false"
	>
		<template #trigger>
			<div
				class="block selected-dashboard"
				style="cursor: pointer;"
				@click="isSelectDashboardPopoverOpened = true"
			>
				<n-avatar
					style="display: flex; align-self: center; border-radius: 111px;"
					:src="currentDashboard?.profileImageUrl"
				/>
				<div style="display: flex; flex-direction: column;">
					<n-text :depth="3" style="font-size: 11px">
						Managing user
					</n-text>
					<n-text>{{ currentDashboard?.displayName }}</n-text>
				</div>

				<IconChevronRight
					:style="{
						transition: '0.2s transform ease',
						transform: `rotate(${!isSelectDashboardPopoverOpened ? 90 : -90}deg)`
					}"
				/>
			</div>
		</template>
		<n-spin v-if="isProfileLoading || isDashboardsLoading"></n-spin>
		<div v-else ref="refPopoverList" class="dashboards-container">
			<n-text :depth="3" style="font-size: 11px">
				Channels you have access to
			</n-text>
			<n-virtual-list
				style="max-height: 400px;" :item-size="42" trigger="none"
				:items="menuOptions"
			>
				<template #default="{ item }">
					<div
						:key="item.key"
						class="item"
						style="height: 42px"
						@click="onSelectDashboard(item.key)"
					>
						<n-avatar :src="item.icon" round size="small" />
						<span> {{ item.label }}</span>
					</div>
				</template>
			</n-virtual-list>
			<template v-if="(usersForSelect.data.value?.users?.length ?? 0) > 10">
				<n-input v-model:value="filterValue" placeholder="Search" />
			</template>
		</div>
	</n-popover>
</template>

<style scoped>
.dashboards-container {
	-webkit-user-select: none;
	-ms-user-select: none;
	user-select: none;
}

.dashboards-container :deep(img) {
	-webkit-user-drag: none;
}

.item {
	display: flex;
	gap: 12px;
	align-items: center;
	width: 100%;
	background-color: v-bind(blockColor);
	padding: 6px;
	border-radius: 6px;
	cursor: pointer;
}

.dashboards-menu > .item:hover {
	background-color: v-bind(blockColor2);
}

.block {
	background-color: v-bind(blockColor);
	display: flex;
	gap: 16px;
	padding: 16px;
	border-radius: 10px;
	align-items: center;
}

.selected-dashboard {
	-webkit-user-select: none;
	-ms-user-select: none;
	user-select: none;
}

.selected-dashboard :deep(img) {
	-webkit-user-drag: none;
}
</style>
