<script setup lang="ts">
import { IconMoon, IconSun } from '@tabler/icons-vue';
import { NButton } from 'naive-ui';

import { useTheme } from '@/hooks';

const { theme, toggleTheme } = useTheme();
</script>

<template>
	<n-button circle quaternary style="padding: 5px" @click="toggleTheme">
		<IconSun v-if="theme === 'dark'" color="orange" />
		<IconMoon v-else />
	</n-button>
</template>
