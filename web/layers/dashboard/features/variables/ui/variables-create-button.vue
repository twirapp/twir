<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile, useUserAccessFlagChecker } from '@/api/auth'
import { useVariablesApi } from '@/api/variables'
import { Button } from '@/components/ui/button'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const variables = useVariablesApi()
const { t } = useI18n()
const { data: profile } = useProfile()

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
		<Button as="a" :href="href" :disabled="isCreateDisabled" @click="navigate">
			<PlusIcon class="size-4 mr-2" />
			{{ t('sharedButtons.create') }} ({{ variables.customVariables.value.length }}/{{ maxVariables }})
		</Button>
	</RouterLink>
</template>
