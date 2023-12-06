<script setup lang="ts">
import { IconLogout } from '@tabler/icons-vue';
import { NButton, NSpin } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { useLogout, useProfile } from '@/api';

const { t } = useI18n();
const { isLoading: isProfileLoading } = useProfile();
const logout = useLogout();
</script>

<template>
	<n-button size="large" :disabled="isProfileLoading" @click="logout.mutate">
		<n-spin v-if="isProfileLoading" size="small" />
		<template v-else>
			<IconLogout style="width: 22px; height: 22px" />
			<span style="display: inline-block; margin-left: 5px">{{ t('navbar.logout') }}</span>
		</template>
	</n-button>
</template>
