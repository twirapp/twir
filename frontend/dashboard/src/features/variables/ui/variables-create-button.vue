<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useUserAccessFlagChecker } from '@/api'
import { useVariablesApi } from '@/api/variables'
import { Button } from '@/components/ui/button'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const variables = useVariablesApi()
const { t } = useI18n()

const userCanManageVariables = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageVariables)
const isCreateDisabled = computed(() => {
	return variables.customVariables.value.length >= 50 || !userCanManageVariables.value
})
</script>

<template>
	<RouterLink v-slot="{ href, navigate }" custom to="/dashboard/variables/create">
		<Button as="a" :href="href" :disabled="isCreateDisabled" @click="navigate">
			<PlusIcon class="size-4 mr-2" />
			{{ t('sharedButtons.create') }} ({{ variables.customVariables.value.length }}/50)
		</Button>
	</RouterLink>
</template>
