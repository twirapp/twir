<script setup lang="ts">
import {
	IconActivity,
	IconBox,
	IconCalendarEvent,
	IconClipboardCopy,
	IconClockHour7,
	IconCommand,
	IconDashboard,
	IconDeviceDesktop,
	IconDeviceDesktopAnalytics,
	IconHeadphones,
	IconKey,
	IconPencilPlus,
	IconPlaylist,
	IconShieldHalfFilled,
	IconSpeakerphone,
	IconSword,
	IconUsers,
} from '@tabler/icons-vue';
import { useMagicKeys } from '@vueuse/core';
import { type MenuOption, type MenuDividerOption, NMenu, NCard, NSpin, NSpace, NAvatar, NText } from 'naive-ui';
import { h, ref, onMounted, computed, watch } from 'vue';
import { RouterLink, useRouter } from 'vue-router';

import DashboardMenu from './dashboardsMenu.vue';
import { renderIcon } from '../helpers/index.js';

import { useProfile, useTwitchGetUsers } from '@/api/index.js';

defineProps<{
	isCollapsed: boolean
}>();

const activeKey = ref<string | null>('/');
const menuOptions: (MenuOption | MenuDividerOption)[] = [
	{
		label: 'Dashboard',
		icon: renderIcon(IconDashboard),
		path: '/dashboard',
	},
	{
		label: 'Integrations',
		icon: renderIcon(IconBox),
		path: '/dashboard/integrations',
	},
	{
		label: 'Events',
		icon: renderIcon(IconCalendarEvent),
		path: '/dashboard/events',
	},
	{
		label: 'OBS Overlays',
		icon: renderIcon(IconDeviceDesktop),
		path: '/dashboard/overlays',
	},
	{
		label: 'Song Requests',
		icon: renderIcon(IconHeadphones),
		path: '/dashboard/song-requests',
	},
	{
		label: 'Commands',
		icon: renderIcon(IconCommand),
		children: [
			{ label: 'Custom', icon: renderIcon(IconPencilPlus), path: '/dashboard/commands/custom' },
			{ label: 'Stats', icon: renderIcon(IconDeviceDesktopAnalytics), path: '/dashboard/commands/stats' },
			{ label: 'Moderation', icon: renderIcon(IconSword), path: '/dashboard/commands/moderation' },
			{ label: 'Songs', icon: renderIcon(IconPlaylist), path: '/dashboard/commands/songs' },
			{ label: 'Manage', icon: renderIcon(IconClipboardCopy), path: '/dashboard/commands/manage' },
		],
	},
	{
		label: 'Users',
		icon: renderIcon(IconUsers),
		path: '/dashboard/community/users',
	},
	{
		label: 'Permissions',
		icon: renderIcon(IconShieldHalfFilled),
		path: '/dashboard/community/roles',
	},
	{
		label: 'Timers',
		icon: renderIcon(IconClockHour7),
		path: '/dashboard/timers',
	},
	// {
	// 	label: 'Moderation',
	// 	icon: renderIcon(IconSword),
	// 	path: '/dashboard/moderation'
	// },
	{
		label: 'Keywords',
		icon: renderIcon(IconKey),
		path: '/dashboard/keywords' },
	{
		label: 'Variables',
		icon: renderIcon(IconActivity),
		path: '/dashboard/variables' },
	{
		label: 'Greetings',
		icon: renderIcon(IconSpeakerphone),
		path: '/dashboard/greetings',
	},
	{
		type: 'divider',
	},
].map((item) => ({
	...item,
	key: item.path ?? item.label,
	label: !item.path ? item.label ?? undefined : () => h(
		RouterLink,
		{
			to: {
				path: item.path,
			},
		},
		{ default: () => item.label },
	),
	children: item.children?.map((child) => ({
		...child,
		key: child.path,
		label: () => h(
			RouterLink,
			{
				to: {
					path: child.path,
				},
			},
			{ default: () => child.label },
		),
	})),
}));

const router = useRouter();

onMounted(async () => {
	await router.isReady();
	activeKey.value = router.currentRoute.value.path;
});

const isDashboardsMenu = ref(false);

const { Ctrl_k } = useMagicKeys({
	passive: false,
	onEventFired(e) {
		if (e.ctrlKey && e.key === 'k' && e.type === 'keydown') {
      e.preventDefault();
		}
  },
});

watch(Ctrl_k, (v) => {
	if (v) {
		isDashboardsMenu.value = !isDashboardsMenu.value;
	}
});

const { data: profile, isLoading: isProfileLoading } = useProfile();
const selectedDashboardId = computed(() => profile.value?.selectedDashboardId ?? '');
const foundTwitchUsers = useTwitchGetUsers({
	ids: selectedDashboardId,
});
const selectedDashboard = computed(() => {
	const twitchUser = foundTwitchUsers.data.value?.users.find(u => u.id === profile.value?.selectedDashboardId);
	if (!twitchUser) return null;

	return twitchUser;
});
</script>

<template>
	<div style="display: flex; flex-direction: column; justify-content: space-between; height: calc(100vh - 43px)">
		<n-menu
			v-if="!isDashboardsMenu"
			v-model:value="activeKey"
			:collapsed-width="64"
			:collapsed-icon-size="22"
			:options="menuOptions"
		/>
		<dashboard-menu v-else @dashboard-selected="isDashboardsMenu = !isDashboardsMenu" />

		<div style="padding: 5px">
			<n-card style="cursor: pointer;" size="small" @click="isDashboardsMenu = !isDashboardsMenu">
				<n-spin v-if="!selectedDashboard || isProfileLoading" />
				<n-space v-else align="center">
					<n-avatar style="display: flex; align-self: center;" :src="selectedDashboard.profileImageUrl" />
					<n-space v-if="!isCollapsed" vertical style="gap: 0; width: 100%">
						<n-text>{{ selectedDashboard.displayName }}</n-text>
						<n-text style="font-size: 12px;">
							Manage dashboard
						</n-text>
					</n-space>
				</n-space>
			</n-card>
		</div>
	</div>
</template>

<style scoped>

</style>
