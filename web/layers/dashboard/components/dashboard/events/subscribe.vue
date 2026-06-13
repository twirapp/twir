<script lang="ts" setup>
import { IconStarFilled } from '@tabler/icons-vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	userName?: string | null
	userDisplayName?: string | null
	level?: string | null
	createdAt: string
}>()
</script>

<template>
	<Base
		v-if="userName && userDisplayName && level"
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
