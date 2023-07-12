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
	IconUsers, SVGProps,
} from '@tabler/icons-vue';
import { NMenu, MenuOption, NAvatar, MenuDividerOption } from 'naive-ui';
import { h, ref, FunctionalComponent, onMounted } from 'vue';
import { RouterLink, useRouter } from 'vue-router';

function renderIcon(icon: (props: SVGProps) => FunctionalComponent<SVGProps>) {
	return () => h(icon, null, { default: () => h(icon) });
}

const activeKey = ref<string | null>('/');
const menuOptions: (MenuOption | MenuDividerOption)[] = [
	{
		label: 'Dashboard',
		icon: renderIcon(IconDashboard),
		path: '/',
	},
	{
		label: 'Integrations',
		icon: renderIcon(IconBox),
		path: '/integrations',
	},
	{
		label: 'Events',
		icon: renderIcon(IconCalendarEvent),
		path: '/events',
	},
	{
		label: 'Overlays',
		icon: renderIcon(IconDeviceDesktop),
		path: '/overlays',
	},
	{
		label: 'Song Requests',
		icon: renderIcon(IconHeadphones),
		path: '/song-requests',
		children: [
			{
				label: 'Player',
				icon: renderIcon(IconPlayerPlay),
				path: '/song-requests/player',
			},
			{
				label: 'Settings',
				icon: renderIcon(IconSettings),
				path: '/song-requests/settings',
			},
		],
	},
	{
		label: 'Commands',
		icon: renderIcon(IconCommand),
		path: '/commands',
		children: [
			{ label: 'Custom', icon: renderIcon(IconPencilPlus), path: '/commands/custom' },
			{ label: 'Stats', icon: renderIcon(IconDeviceDesktopAnalytics), path: '/commands/stats' },
			{ label: 'Moderation', icon: renderIcon(IconSword), path: '/commands/moderation' },
			{ label: 'Songs', icon: renderIcon(IconPlaylist), path: '/commands/songs' },
			{ label: 'Manage', icon: renderIcon(IconClipboardCopy), path: '/commands/manage' },
		],
	},
	{
		label: 'Community',
		icon: renderIcon(IconUsers),
		path: '/community',
		children: [
			{ label: 'Users', icon: renderIcon(IconUsers), path: '/community/users' },
			{ label: 'Roles', icon: renderIcon(IconShieldHalfFilled), path: '/community/roles' },
		],
	},
	{
		label: 'Timers',
		icon: renderIcon(IconClockHour7),
		path: '/timers',
	},
	{
		label: 'Moderation',
		icon: renderIcon(IconSword),
		path: '/moderation' },
	{
		label: 'Keywords',
		icon: renderIcon(IconKey),
		path: '/keywords' },
	{
		label: 'Variables',
		icon: renderIcon(IconActivity),
		path: '/variables' },
	{
		label: 'Greetings',
		icon: renderIcon(IconSpeakerphone),
		path: '/greetings',
	},
	{
		type: 'divider',
	},
].map((item) => ({
	...item,
	key: item.path,
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
