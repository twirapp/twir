<script setup lang="ts">
import { IconLogout } from '@tabler/icons-vue'
import { NButton, NSpin } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import { useLogout, useProfile } from '@/api'

const { t } = useI18n()
const { isLoading: isProfileLoading } = useProfile()
const logout = useLogout()
</script>

<template>
	<NButton
		secondary
		block
		type="error"
		:disabled="isProfileLoading"
		@click="logout.execute"
	>
		<NSpin v-if="isProfileLoading" size="small" />

		<template v-else>
			<IconLogout class="w-5 h-5" />
			<span class="inline-block ml-1">{{ t('navbar.logout') }}</span>
		</template>
	</NButton>
</template>
