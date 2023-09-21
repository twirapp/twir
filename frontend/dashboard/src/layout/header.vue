<script setup lang="ts">
import { IconLogout, IconShare, IconMoon, IconSun } from '@tabler/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import {
	type DropdownOption,
	NAvatar,
	NButton,
	NDropdown,
	NTooltip,
	NSpin,
	useThemeVars,
} from 'naive-ui';
import { computed, defineAsyncComponent } from 'vue';
import { useI18n } from 'vue-i18n';

import DashboardsMenu from './dashboardsMenu.vue';
import Logo from '../../public/TwirInCircle.svg?component';

import { useLogout, useProfile, useTwitchGetUsers } from '@/api/index.js';
import DiscordLogo from '@/assets/icons/integrations/discord.svg?component';
import { renderIcon } from '@/helpers/index.js';
import { useTheme } from '@/hooks/index.js';

defineProps<{
	toggleSidebar: () => void;
}>();

const { t, availableLocales, locale } = useI18n({ useScope: 'global' });
const { theme, toggleTheme } = useTheme();

const themeVars = useThemeVars();
const blockColor = computed(() => themeVars.value.buttonColor2);
const discordIconColor = computed(() => themeVars.value.textColor2);

const logout = useLogout();
const profileOptions: DropdownOption[] = [{
	label: t('navbar.logout'),
	key: 'logout',
	icon: renderIcon(IconLogout),
	props: {
		onClick: () => {
			logout.mutate();
		},
		style: {
			// 'background-color': 'red',
		},
	},
}];

const { data: profileData, isLoading: isProfileLoading } = useProfile();

const localStorageLocale = useLocalStorage('twirLocale', 'en');

const selectedUserId = computed(() => {
	return (profileData.value?.selectedDashboardId ?? profileData?.value?.id) || '';
});
const selectedDashboardTwitchUser = useTwitchGetUsers({
	ids: selectedUserId,
});

const openDiscord = () => window.open('https://discord.gg/Q9NBZq3zVV', '_blank');
const publicPageHref = computed<string>(() => {
	if (!profileData.value || !selectedDashboardTwitchUser.data.value?.users.length) return '';

	return `${window.location.origin}/p/${selectedDashboardTwitchUser.data.value.users.at(0)!.login}`;
});

const renderFlagIcon = (code: string) => defineAsyncComponent(() => import(`@/assets/icons/flags/${code}.svg?component`));
</script>

<template>
	<div class="header">
		<div class="content">
			<div>
				<a href="/" class="logo">
					<Logo style="width: 36px; height: 36px; display: flex" />
				</a>

				<DashboardsMenu />
			</div>

			<div style="display: flex; gap: 12px;">
				<div class="block">
					<n-tooltip>
						<template #trigger>
							<n-button tag="a" circle quaternary target="_blank" :href="publicPageHref">
								<IconShare />
							</n-button>
						</template>
						{{ t('navbar.publicPage') }}
					</n-tooltip>
					<n-button
						quaternary
						circle
						style="padding: 5px"
						@click="openDiscord"
					>
						<DiscordLogo :style="{ width: '20px', fill: discordIconColor }" />
					</n-button>
				</div>

				<div class="block">
					<n-dropdown
						trigger="click"
						:options="availableLocales.map(l => ({
							title: t('languageName', {}, { locale: l }),
							key: l as string,
						}))"
						size="medium"
						@select="(l) => {
							locale = l
							localStorageLocale = l
						}"
					>
						<n-button circle quaternary style="padding: 5px; font-size: 25px">
							<component :is="renderFlagIcon(localStorageLocale)" style="width: 35px; height: 50px;" />
						</n-button>
					</n-dropdown>
					<n-button circle quaternary style="padding: 5px" @click="toggleTheme">
						<IconSun v-if="theme === 'dark'" color="orange" />
						<IconMoon v-else />
					</n-button>
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
				</div>
			</div>
		</div>
	</div>
</template>

<style>
.header {
	height: 100%;
	padding-top: 10px;
	padding-left: 10px;
	padding-right: 10px;
}

.content .logo {
	display: flex;
	align-items: center;
	text-decoration: none;
	padding-right: 14px;
}

.content > div {
	display: flex;
	justify-content: flex-start;
	gap: 5px;
	height: 45px;
}

.content {
	display: flex;
	justify-content: space-between;
	align-items: center;
	height: 45px;
}
.profile {
	display: flex;
	gap: 5px;
	align-items: center;
}

.block {
	background-color: v-bind(blockColor);
	display: flex;
	gap: 16px;
	padding: 16px;
	border-radius: 10px;
	align-items: center;
}
</style>
