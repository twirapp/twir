<script setup lang="ts">
import { IconLogout, IconMenu2, IconMoon, IconSun } from '@tabler/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import {
	type DropdownOption,
	NAvatar,
	NButton,
	NDivider,
	NDropdown,
	NSpace,
	NSpin,
	NText,
	useThemeVars,
} from 'naive-ui';
import { computed, defineAsyncComponent } from 'vue';
import { useI18n } from 'vue-i18n';

import Logo from '../../public/TwirInCircle.svg?component';

import { useLogout, useProfile, useTwitchGetUsers } from '@/api/index.js';
import DiscordLogo from '@/assets/icons/integrations/discord.svg?component';
import { renderIcon } from '@/helpers/index.js';
import { useTheme } from '@/hooks/index.js';

defineProps<{
	toggleSidebar: () => void;
}>();

const { theme, toggleTheme } = useTheme();

const themeVars = useThemeVars();
const discordIconColor = computed(() => themeVars.value.textColor2);

const logout = useLogout();
const profileOptions: DropdownOption[] = [{
	label: 'Logout',
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

const { t, availableLocales, locale } = useI18n({ useScope: 'global' });

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
		<div>
			<n-button text @click="toggleSidebar">
				<IconMenu2 />
			</n-button>
			<a href="/" style="text-decoration: none">
				<n-space align="center" style="gap: 8px; margin-left: 5px;">
					<Logo style="width: 25px; height: 25px; display: flex" />
					<n-text strong style="font-size: 20px;">
						Twir
					</n-text>
				</n-space>
			</a>
		</div>

		<div>
			<n-button tag="a" tertiary target="_blank" :href="publicPageHref">
				Public page
			</n-button>
			<n-button
				tertiary
				style="padding: 5px"
				@click="openDiscord"
			>
				<DiscordLogo :style="{ width: '20px', fill: discordIconColor }" />
			</n-button>
			<n-divider vertical style="height: 2em" />
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
				<n-button tertiary style="padding: 5px; font-size: 25px">
					<component :is="renderFlagIcon(localStorageLocale)" style="width: 35px; height: 50px;" />
				</n-button>
			</n-dropdown>
			<n-button tertiary style="padding: 5px" @click="toggleTheme">
				<IconSun v-if="theme === 'dark'" color="orange" />
				<IconMoon v-else />
			</n-button>
			<n-dropdown trigger="click" :options="profileOptions" size="large">
				<n-button tertiary>
					<n-spin v-if="isProfileLoading" size="small" />
					<div v-else class="profile">
						<n-avatar
							size="small"
							:src="profileData?.avatar"
							round
						/>
						{{ profileData?.displayName }}
					</div>
				</n-button>
			</n-dropdown>
		</div>
	</div>
</template>

<style scoped>
.header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-left: 10px;
	margin-right: 10px;
	height: 100%;
}

.header > div {
	display: flex;
	justify-content: flex-start;
	gap: 5px;
	align-items: center;
}

.profile {
	display: flex;
	gap: 5px;
	align-items: center;
}
</style>
