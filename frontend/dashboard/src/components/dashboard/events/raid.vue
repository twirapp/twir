<script lang="ts" setup>
import { IconPick } from '@tabler/icons-vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	userName: string
	userDisplayName: string
	createdAt: string
	viewers: string
}>()
</script>

<template>
	<Base
		:icon="IconPick"
		:icon-color="['#949400', '#ebeb00']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<span>
				<UserLink :name="userName" :display-name="userDisplayName" />{{ '' }}
				<span class="font-bold">raided with</span> {{ viewers }} viewers
			</span>
		</template>
	</Base>
</template>
