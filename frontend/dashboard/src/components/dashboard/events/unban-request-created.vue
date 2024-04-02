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
			<div class="flex flex-col">
				<span>
					{{ createUserName(userLogin, userName) }}
					<b class="text-xs">requested unban</b>
				</span>
				<span class="text-xs">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
