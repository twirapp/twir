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
			<div style="display: flex; flex-direction: column">
				<span>{{ createUserName(moderatorUserLogin, moderatorUserName) }} <b style="font-size: 11px">banned</b> {{ createUserName(userLogin, userName) }} for {{ endsIn }}</span>
				<span style="font-size: 11px">{{ reason }}</span>
			</div>
		</template>
	</Base>
</template>
