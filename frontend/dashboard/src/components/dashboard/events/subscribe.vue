<script lang="ts" setup>
import { IconStarFilled } from '@tabler/icons-vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	userName: string
	userDisplayName: string
	createdAt: string
	level: string
}>()
</script>

<template>
	<Base
		:icon="IconStarFilled"
		:icon-color="['#1756d3', '#1f69ff']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<span>
				<UserLink :name="userName" :display-name="userDisplayName" />{{ '' }}
				<span class="font-bold">subscribed with</span> {{ level }} tier
			</span>
		</template>
	</Base>
</template>
