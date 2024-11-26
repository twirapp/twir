<script lang="ts" setup>
import { IconHeartFilled } from '@tabler/icons-vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	userName?: string | null
	userDisplayName?: string | null
	createdAt: string
}>()
</script>

<template>
	<Base
		v-if="userName && userDisplayName"
		:icon="IconHeartFilled"
		:icon-color="['#c516a5', '#ff38db']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<UserLink :name="userName" :display-name="userDisplayName" />
		</template>
	</Base>
</template>
