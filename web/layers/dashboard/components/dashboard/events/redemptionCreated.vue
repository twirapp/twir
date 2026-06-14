<script lang="ts" setup>

import Base from './base.vue'
import UserLink from './user-link.vue'

defineProps<{
	userName?: string | null
	userDisplayName?: string | null
	title?: string | null
	input?: string | null
	cost?: string | null
	createdAt: string
}>()
</script>

<template>
	<Base
		v-if="userName && userDisplayName && title && cost"
		:icon="'tabler:icon-award'"
		:icon-color="['#c28100', '#ffd37a']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<div class="flex flex-col">
				<span class="break-all">
					<UserLink :name="userName" :display-name="userDisplayName" />{{ '' }}
					<span class="font-bold">activated</span> {{ title }} ({{ cost }})
				</span>
				<span v-if="input?.length" class="text-xs">{{ input }}</span>
			</div>
		</template>
	</Base>
</template>
