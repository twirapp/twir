<script lang="ts" setup>

import { IconUserCancel } from '@tabler/icons-vue';

import Base from './base.vue';

import { resolveUserName } from '@/helpers';

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
			<div class="flex flex-col">
				<span>
					{{ resolveUserName(moderatorUserLogin, moderatorUserName) }}
					<b class="text-xs">resolved unban request from</b> {{ resolveUserName(userLogin, userName) }}
				</span>
				<span class="text-xs">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
