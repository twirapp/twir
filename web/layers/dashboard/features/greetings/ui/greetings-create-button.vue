<script setup lang="ts">
import { computed } from 'vue'

import GreetingsDialog from './greetings-dialog.vue'

import { useProfile, useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import { useGreetingsApi } from '~~/layers/dashboard/api/greetings'
import { Button } from '@/components/ui/button'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

const { t } = useI18n()
const { data: profile } = useProfile()
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
				<Button :disabled="isCreateDisabled">
					<Icon name="lucide:plus" class="size-4 mr-2" />
					{{ greetingsLength >= maxGreetings ? t('greetings.limitExceeded') : t('greetings.create') }} ({{
						greetingsLength }}/{{ maxGreetings }})
				</Button>
			</template>
		</GreetingsDialog>
	</div>
</template>
