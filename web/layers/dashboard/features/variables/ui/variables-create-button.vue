<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'


import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useVariablesApi } from '#layers/dashboard/api/variables'

import { ChannelRolePermissionEnum } from '~/gql/graphql'

const variables = useVariablesApi()
const { t } = useI18n()
const { user: profile } = storeToRefs(useDashboardAuth())

const userCanManageVariables = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageVariables)

const maxVariables = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxVariables ?? 50
})

const isCreateDisabled = computed(() => {
	return variables.customVariables.value.length >= maxVariables.value || !userCanManageVariables.value
})
</script>

<template>
	<RouterLink v-slot="{ href, navigate }" custom to="/dashboard/variables/create">
		<UiButton as="a" :href="href" :disabled="isCreateDisabled" @click="navigate">
			<PlusIcon class="size-4 mr-2" />
			{{ t('sharedButtons.create') }} ({{ variables.customVariables.value.length }}/{{ maxVariables }})
		</UiButton>
	</RouterLink>
</template>
