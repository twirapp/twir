<script setup lang='ts'>

import { IconMenu2, IconSun, IconMoon, IconLogout } from '@tabler/icons-vue';
import type { DropdownOption } from 'naive-ui';

import { useProfile, useLogout } from '@/api/index.js';
import { renderIcon } from '@/helpers/index.js';
import { useTheme } from '@/hooks/index.js';

defineProps<{
	toggleSidebar: () => void;
}>();

const { theme, toggleTheme } = useTheme();

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
</script>

<template>
  <div class="header">
    <div>
      <n-button text @click="toggleSidebar">
        <IconMenu2 />
      </n-button>
      <n-text strong>
        TwirApp
      </n-text>
    </div>

    <div>
      <n-button tertiary @click="toggleTheme">
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
