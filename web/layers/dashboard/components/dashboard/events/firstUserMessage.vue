<script lang="ts" setup>
import { IconMoodPlus } from '@tabler/icons-vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	userName?: string | null
	userDisplayName?: string | null
	message?: string | null
	createdAt: string
}>()
</script>

<template>
	<Base
		v-if="userName && userDisplayName && message"
		:icon="IconMoodPlus"
		:icon-color="['#772ce8', '#9147ff']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<div class="flex flex-col">
				<span>
					<UserLink :name="userName" :display-name="userDisplayName" />{{ '' }}
					<span class="font-bold">first time message</span>
				</span>
				<span class="text-sx">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
