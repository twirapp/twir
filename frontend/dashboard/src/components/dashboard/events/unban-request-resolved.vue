<script lang="ts" setup>
import { IconUserCancel } from '@tabler/icons-vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	createdAt: string
	userName: string
	userLogin: string
	message: string
	moderatorUserName: string
	moderatorUserLogin: string
}>()
</script>

<template>
	<Base
		:icon="IconUserCancel"
		:icon-color="['#ff4f4d', '#ffaaa8']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<div class="flex flex-col">
				<span>
					<UserLink :name="moderatorUserLogin" :display-name="moderatorUserName" />{{ '' }}
					<span class="font-bold">resolved unban request from</span>{{ '' }}
					<UserLink :name="userLogin" :display-name="userName" />
				</span>
				<span class="text-xs">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
