<script lang="ts" setup>
import { IconGift } from '@tabler/icons-vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	userName: string
	userDisplayName: string
	targetUserName: string
	targetUserDisplayName: string
	createdAt: string
	level: string
}>()
</script>

<template>
	<Base
		:icon="IconGift"
		:icon-color="['#1756d3', '#1f69ff']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<span>
				<UserLink :name="userName" :display-name="userDisplayName" />{{ '' }}
				<span class="font-bold">gifted sub to</span>{{ '' }}
				<UserLink :name="targetUserName" :display-name="targetUserDisplayName" />
			</span>
			<span>
				<span class="text-xs">
					{{ level }} tier
				</span>
			</span>
		</template>
	</Base>
</template>
