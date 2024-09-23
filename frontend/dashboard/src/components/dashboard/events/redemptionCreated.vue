<script lang="ts" setup>
import { IconAward } from '@tabler/icons-vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	createdAt: string
	userName: string
	userDisplayName: string
	title: string
	input: string
	cost: string
}>()
</script>

<template>
	<Base
		:icon="IconAward"
		:icon-color="['#c28100', '#ffd37a']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<div class="flex flex-col">
				<span class="break-all">
					<UserLink :name="userName" :display-name="userDisplayName" />{{ '' }}
					<span class="font-bold">activated</span> {{ title }} ({{ cost }})
				</span>
				<span v-if="input.length" class="text-xs">{{ input }}</span>
			</div>
		</template>
	</Base>
</template>
