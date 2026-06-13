<script lang="ts" setup>
import { IconBan, IconClockCancel } from '@tabler/icons-vue'
import { computed } from 'vue'

import Base from './base.vue'
import UserLink from './user-link.vue'

const props = defineProps<{
	userName?: string | null
	userLogin?: string | null
	reason?: string | null
	moderatorUserName?: string | null
	moderatorUserLogin?: string | null
	endsIn?: string | null
	createdAt: string
}>()

const iconBan = computed(() => {
	if (props.endsIn === 'permanent') return IconBan
	return IconClockCancel
})
</script>

<template>
	<Base
		v-if="moderatorUserLogin && moderatorUserName && userLogin && userName"
		:icon="iconBan"
		:icon-color="['#ff4f4d', '#ffaaa8']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<div class="flex flex-col">
				<span>
					<UserLink :name="moderatorUserLogin" :display-name="moderatorUserName" />{{ '' }}
					<span class="font-bold">banned</span>{{ '' }}
					<UserLink :name="userLogin" :display-name="userName" />
					for {{ endsIn }} {{ endsIn !== 'permanent' ? 'minutes' : '' }}</span>
				<span class="text-xs">{{ reason }}</span>
			</div>
		</template>
	</Base>
</template>
