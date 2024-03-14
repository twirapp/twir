<script lang="ts" setup>

import { IconUserCancel } from '@tabler/icons-vue';

import Base from './base.vue';

import { createUserName } from '@/helpers';

defineProps<{
	createdAt: string,
	userName: string,
	userLogin: string
	message: string
	moderatorUserName: string
	moderatorUserLogin: string
}>();
</script>

<template>
	<Base :icon="IconUserCancel" :created-at="createdAt">
		<template #leftContent>
			<div style="display: flex; flex-direction: column">
				<span>
					{{ createUserName(moderatorUserLogin, moderatorUserName) }}
					<b style="font-size: 11px">resolved unban request from</b> {{ createUserName(userLogin, userName) }}
				</span>
				<span style="font-size: 11px">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
