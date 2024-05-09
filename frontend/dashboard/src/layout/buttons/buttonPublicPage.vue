<script setup lang="ts">
import { IconExternalLink } from '@tabler/icons-vue'
import { NButton, NTooltip } from 'naive-ui'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api'

const { data: profileData } = useProfile()
const { t } = useI18n()

const publicPageHref = computed(() => {
	const selectedDashboardLogin = profileData.value?.selectedDashboardTwitchUser?.login
	if (!selectedDashboardLogin) return null

	return `${window.location.origin}/p/${selectedDashboardLogin}`
})
</script>

<template>
	<NTooltip v-if="publicPageHref">
		<template #trigger>
			<NButton tag="a" circle quaternary target="_blank" :href="publicPageHref">
				<IconExternalLink />
			</NButton>
		</template>
		{{ t('navbar.publicPage') }}
	</NTooltip>
</template>
