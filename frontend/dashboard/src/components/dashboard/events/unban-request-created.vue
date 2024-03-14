<script lang="ts" setup>

import { IconUserCancel } from '@tabler/icons-vue';

import Base from './base.vue';

import { createUserName } from '@/helpers';

defineProps<{
	userName: string,
	userLogin: string
	createdAt: string,
	message: string
}>();
</script>

<template>
	<Base :icon="IconUserCancel" :created-at="createdAt">
		<template #leftContent>
			<div style="display: flex; flex-direction: column">
				<span>
					{{ createUserName(userLogin, userName) }}
					<b style="font-size: 11px">requested unban</b>
				</span>
				<span style="font-size: 11px">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
