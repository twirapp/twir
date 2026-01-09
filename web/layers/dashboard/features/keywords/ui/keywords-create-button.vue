<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'


import KeywordsDialog from './keywords-dialog.vue'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useKeywordsApi } from '#layers/dashboard/api/keywords'

import { ChannelRolePermissionEnum } from '~/gql/graphql'

const { t } = useI18n()
const { user: profile } = storeToRefs(useDashboardAuth())
const keywordsApi = useKeywordsApi()
const keywords = keywordsApi.useQueryKeywords()
const userCanManageKeywords = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageKeywords)

const maxKeywords = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxKeywords ?? 50
})

const keywordsLength = computed(() => keywords.data.value?.keywords.length ?? 0)

const isCreateDisabled = computed(() => {
	return keywordsLength.value >= maxKeywords.value || !userCanManageKeywords.value
})
</script>

<template>
	<div class="flex gap-2">
		<KeywordsDialog>
			<template #dialog-trigger>
				<UiButton :disabled="isCreateDisabled">
					<PlusIcon class="size-4 mr-2" />
					{{ keywordsLength >= maxKeywords ? t('keywords.limitExceeded') : t('keywords.create') }} ({{
						keywordsLength }}/{{ maxKeywords }})
				</UiButton>
			</template>
		</KeywordsDialog>
	</div>
</template>
