<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useUserAccessFlagChecker } from '@/api'
import { Button } from '@/components/ui/button'
import { useTimersEdit } from '@/features/timers/composables/use-timers-edit'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const { t } = useI18n()
const userCanManageTimers = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageTimers)

const { timers } = useTimersEdit()
const timersLength = computed(() => timers.data?.value?.timers.length ?? 0)
</script>

<template>
	<div class="flex gap-2">
		<RouterLink v-slot="{ navigate, href }" custom to="/dashboard/timers/create">
			<Button
				as="a"
				:href="href"
				:disabled="!userCanManageTimers || timersLength >= 10"
				@click="navigate"
			>
				<PlusIcon class="size-4 mr-2" />
				{{ timersLength >= 10 ? t('timers.limitExceeded') : t('sharedButtons.create') }} ({{
					timersLength }}/10)
			</Button>
		</RouterLink>
	</div>
</template>
