<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'


import GreetingsDialog from './greetings-dialog.vue'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useGreetingsApi } from '#layers/dashboard/api/greetings'

import { ChannelRolePermissionEnum } from '~/gql/graphql'

const { t } = useI18n()
const { user: profile } = storeToRefs(useDashboardAuth())
const greetingsApi = useGreetingsApi()
const userCanManageGreetings = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageGreetings)

const { data: greetingsData } = greetingsApi.useQueryGreetings()

const maxGreetings = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxGreetings ?? 50
})

const greetingsLength = computed(() => greetingsData.value?.greetings.length ?? 0)

const isCreateDisabled = computed(() => {
	return greetingsLength.value >= maxGreetings.value || !userCanManageGreetings.value
})
</script>

<template>
	<div class="flex gap-2">
		<GreetingsDialog>
			<template #dialog-trigger>
				<UiButton :disabled="isCreateDisabled">
					<PlusIcon class="size-4 mr-2" />
					{{ greetingsLength >= maxGreetings ? t('greetings.limitExceeded') : t('greetings.create') }} ({{
						greetingsLength }}/{{ maxGreetings }})
				</UiButton>
			</template>
		</GreetingsDialog>
	</div>
</template>
