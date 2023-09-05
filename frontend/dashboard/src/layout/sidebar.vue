<script setup lang="ts">
import {
	IconActivity,
	IconBell,
	IconBox,
	IconCalendarEvent,
	IconClipboardCopy,
	IconClockHour7,
	IconCommand,
	IconDashboard,
	IconDeviceDesktop,
	IconDeviceDesktopAnalytics, IconDeviceGamepad2,
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
import {
	type MenuDividerOption,
	type MenuOption,
	NAvatar,
	NCard,
	NMenu,
	NScrollbar,
	NSpace,
	NSpin,
	NText,
	NBadge,
} from 'naive-ui';
import { computed, h, onMounted, ref, watch } from 'vue';
import { RouterLink, useRouter } from 'vue-router';

import DashboardMenu from './dashboardsMenu.vue';
import { renderIcon } from '../helpers/index.js';

import { useProfile, useTwitchGetUsers, useUserAccessFlagChecker } from '@/api/index.js';

defineProps<{
	isCollapsed: boolean
}>();

const router = useRouter();

const activeKey = ref<string | null>('/');
const menuOptions = computed<(MenuOption | MenuDividerOption)[]>(() => {
	const canViewIntegrations = useUserAccessFlagChecker('VIEW_INTEGRATIONS');
	const canViewEvents = useUserAccessFlagChecker('VIEW_EVENTS');
	const canViewOverlays = useUserAccessFlagChecker('VIEW_OVERLAYS');
	const canViewSongRequests = useUserAccessFlagChecker('VIEW_SONG_REQUESTS');
	const canViewCommands = useUserAccessFlagChecker('VIEW_COMMANDS');
	const canViewTimers = useUserAccessFlagChecker('VIEW_TIMERS');
	const canViewKeywords = useUserAccessFlagChecker('VIEW_KEYWORDS');
	const canViewVariabls = useUserAccessFlagChecker('VIEW_VARIABLES');
	const canViewGreetings = useUserAccessFlagChecker('VIEW_GREETINGS');
	const canViewRoles = useUserAccessFlagChecker('VIEW_ROLES');
	const canViewAlerts = useUserAccessFlagChecker('VIEW_ALERTS');
	const canViewGames = useUserAccessFlagChecker('VIEW_GAMES');

	return [
		{
			label: 'Dashboard',
			icon: renderIcon(IconDashboard),
			path: '/dashboard',
		},
		{
			label: 'Integrations',
			icon: renderIcon(IconBox),
			path: '/dashboard/integrations',
			disabled: !canViewIntegrations.value,
		},
		{
			label: 'Alerts',
			icon: renderIcon(IconBell),
			path: '/dashboard/alerts',
			disabled: !canViewAlerts.value,
			isNew: true,
		},
		{
			label: 'Events',
			icon: renderIcon(IconCalendarEvent),
			path: '/dashboard/events',
			disabled: !canViewEvents.value,
		},
		{
			label: 'OBS Overlays',
			icon: renderIcon(IconDeviceDesktop),
			path: '/dashboard/overlays',
			disabled: !canViewOverlays.value,
		},
		{
			label: 'Song Requests',
			icon: renderIcon(IconHeadphones),
			path: '/dashboard/song-requests',
			disabled: !canViewSongRequests.value,
		},
		{
			label: 'Games',
			icon: renderIcon(IconDeviceGamepad2),
			path: '/dashboard/games',
			disabled: !canViewSongRequests.value,
			isNew: true,
		},
		{
			label: 'Commands',
			icon: renderIcon(IconCommand),
			disabled: !canViewCommands.value,
			children: [
				{
					label: 'Custom',
					icon: renderIcon(IconPencilPlus),
					path: '/dashboard/commands/custom',
				},
				{
					label: 'Stats',
					icon: renderIcon(IconDeviceDesktopAnalytics),
					path: '/dashboard/commands/stats',
				},
				{
					label: 'Moderation',
					icon: renderIcon(IconSword),
					path: '/dashboard/commands/moderation',
				},
				{
					label: 'Songs',
					icon: renderIcon(IconPlaylist),
					path: '/dashboard/commands/songs',
				},
				{
					label: 'Manage',
					icon: renderIcon(IconClipboardCopy),
					path: '/dashboard/commands/manage',
				},
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
			disabled: !canViewRoles.value,
		},
		{
			label: 'Timers',
			icon: renderIcon(IconClockHour7),
			path: '/dashboard/timers',
			disabled: !canViewTimers.value,
		},
		{
			label: 'Keywords',
			icon: renderIcon(IconKey),
			path: '/dashboard/keywords',
			disabled: !canViewKeywords.value,
		},
		{
			label: 'Variables',
			icon: renderIcon(IconActivity),
			path: '/dashboard/variables',
			disabled: !canViewVariabls.value,
		},
		{
			label: 'Greetings',
			icon: renderIcon(IconSpeakerphone),
			path: '/dashboard/greetings',
			disabled: !canViewGreetings.value,
		},
		{
			type: 'divider',
		},
	].map((item) => ({
		...item,
		key: item.path ?? item.label,
		extra: item.disabled ? 'No perms' : undefined,
		label: !item.path || item.disabled ? item.label ?? undefined : () => h(
			RouterLink,
			{
				to: {
					path: item.path,
				},
			},
			{
				default: () => item.isNew
					? h(NBadge, {
						type: 'info',
						value: 'new',
						processing: true,
						offset: [17, 5],
					}, { default: () => item.label })
					: item.label,
			},
		),
		children: item.children?.map((child) => ({
			...child,
			key: child.path,
			label: item.disabled ? child.label : () => h(
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
});

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
	<div
		style="display: flex; flex-direction: column; justify-content: space-between; height: calc(100vh - 43px)"
	>
		<n-scrollbar trigger="none">
			<n-menu
				v-if="!isDashboardsMenu"
				v-model:value="activeKey"
				:collapsed-width="64"
				:collapsed-icon-size="22"
				:options="menuOptions"
			/>
			<dashboard-menu v-else @dashboard-selected="isDashboardsMenu = !isDashboardsMenu" />
		</n-scrollbar>

		<div style="padding: 5px">
			<n-card style="cursor: pointer;" size="small" @click="isDashboardsMenu = !isDashboardsMenu">
				<n-spin v-if="!selectedDashboard || isProfileLoading" />
				<n-space v-else align="center">
					<n-avatar
						style="display: flex; align-self: center;"
						:src="selectedDashboard.profileImageUrl"
					/>
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
