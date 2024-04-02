<script lang="ts" setup>

import { IconUserCancel } from '@tabler/icons-vue';

import Base from './base.vue';

import { createUserName } from '@/helpers';

defineProps<{
	userName: string,
	userLogin: string
	createdAt: string,
	reason: string
	moderatorUserName: string
	moderatorUserLogin: string
	endsIn: string
}>();
</script>

<template>
	<Base :icon="IconUserCancel" :created-at="createdAt">
		<template #leftContent>
			<div class="flex flex-col">
				<span>
					{{ createUserName(moderatorUserLogin, moderatorUserName) }}
					<b class="text-xs">banned</b> {{ createUserName(userLogin, userName) }}
					for {{ endsIn }} {{ endsIn != 'permanent' ? 'minutes' : '' }}</span>
				<span class="text-xs">{{ reason }}</span>
			</div>
		</template>
	</Base>
</template>
