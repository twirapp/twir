<script lang="ts" setup>
import { IconCalendarStar } from '@tabler/icons-vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	userName: string
	userDisplayName: string
	createdAt: string
	level: string
	months: string
	streak: string
}>()
</script>

<template>
	<Base
		:icon="IconCalendarStar"
		:icon-color="['#1756d3', '#1f69ff']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<span>
				<UserLink :name="userName" :display-name="userDisplayName" />{{ '' }}
				<span class="font-bold">resubscribed with</span> {{ level }} tier
			</span>
			<span>
				<span class="text-xs">
					{{ streak }} months • overall {{ months }} months
				</span>
			</span>
		</template>
	</Base>
</template>
