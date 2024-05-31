<script setup lang="ts">
import { IconUsers } from '@tabler/icons-vue'
import { NButton, NTooltip } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import { usePublicPageHref } from '@/layout/use-public-page-href'

const { t } = useI18n()
const publicPageHref = usePublicPageHref()
</script>

<template>
	<NTooltip v-if="publicPageHref">
		<template #trigger>
			<NButton tag="a" circle quaternary target="_blank" :href="publicPageHref">
				<IconUsers />
			</NButton>
		</template>

		{{ t('navbar.publicPage') }}
	</NTooltip>
</template>
