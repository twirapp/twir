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
	IconPencilPlus, IconPlayerPlay,
	IconPlaylist, IconSettings,
	IconShieldHalfFilled,
	IconSpeakerphone,
	IconSword,
	IconUsers,
} from '@tabler/icons-vue';
import { NMenu, MenuOption, NAvatar, MenuDividerOption } from 'naive-ui';
import { h, ref, onMounted } from 'vue';
import { RouterLink, useRouter } from 'vue-router';

import { renderIcon } from '../helpers/index.js';

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
		label: 'Overlays',
		icon: renderIcon(IconDeviceDesktop),
		path: '/dashboard/overlays',
	},
	{
		label: 'Song Requests',
		icon: renderIcon(IconHeadphones),
		children: [
			{
				label: 'Player',
				icon: renderIcon(IconPlayerPlay),
				path: '/dashboard/song-requests/player',
			},
			{
				label: 'Settings',
				icon: renderIcon(IconSettings),
				path: '/dashboard/song-requests/settings',
			},
		],
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
		label: 'Community',
		icon: renderIcon(IconUsers),
		children: [
			{ label: 'Users', icon: renderIcon(IconUsers), path: '/dashboard/community/users' },
			{ label: 'Roles', icon: renderIcon(IconShieldHalfFilled), path: '/dashboard/community/roles' },
		],
	},
	{
		label: 'Timers',
		icon: renderIcon(IconClockHour7),
		path: '/dashboard/timers',
	},
	{
		label: 'Moderation',
		icon: renderIcon(IconSword),
		path: '/dashboard/moderation' },
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

</script>

<template>
  <div style="display: flex; flex-direction: column; justify-content: space-between; height: calc(100vh - 43px)">
    <n-menu
      v-model:value="activeKey"
      :collapsed-width="64"
      :collapsed-icon-size="22"
      :options="menuOptions"
    />

    <n-menu
      :options="[{
        label: 'Logout',
        key: 'logout',
        onClick: () => {
          window.location.href = '/logout';
        },
      }]"
    >
    </n-menu>
  </div>
</template>

<style scoped>

</style>
