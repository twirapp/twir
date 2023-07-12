<script setup lang='ts'>

import { IconMenu2, IconSun, IconMoon } from '@tabler/icons-vue';
import { NButton, NDropdown, NAvatar } from 'naive-ui';

import { useTheme } from '../hooks/index.js';

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

const profileOptions = [{
	label: 'Marina Bay Sands',
	key: 'Marina Bay Sands',
	disabled: false,
}];
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
      <n-dropdown trigger="click" :options="profileOptions">
        <n-button>
          <n-avatar
            size="small"
            src="https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg?t=1"
          />
          Satont
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
</style>
