<script setup lang='ts'>

import { IconMenu2, IconSun, IconMoon, IconLogout } from '@tabler/icons-vue';
import { NButton, NDropdown, NAvatar, NSpin, DropdownOption } from 'naive-ui';

import { useProfile } from '@/api/index.js';
import { renderIcon } from '@/helpers/index.js';
import { useTheme } from '@/hooks/index.js';

defineProps<{
	toggleSidebar: () => void;
}>();

const theme = useTheme();
const changeTheme = () => {
	if (theme.value === 'dark') {
		theme.value = 'light';
	} else {
		theme.value = 'dark';
	}
};

const profileOptions: DropdownOption[] = [{
	label: 'Logout',
	key: 'logout',
	icon: renderIcon(IconLogout),
	props: {
		onClick: () => {
			console.log('logout');
		},
		style: {
			// 'background-color': 'red',
		},
	},
}];

const { data: profileData, isLoading: isProfileLoading } = useProfile();
</script>

<template>
  <div class="header">
    <div>
      <n-button text @click="toggleSidebar">
        <IconMenu2 />
      </n-button>
      TwirApp
    </div>

    <div>
      <n-button tertiary @click="changeTheme">
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
