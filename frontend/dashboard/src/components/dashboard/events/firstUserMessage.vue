<script lang="ts" setup>
import { IconMoodPlus } from '@tabler/icons-vue';

import Base from './base.vue';

import { createUserName } from '@/helpers/index.js';

defineProps<{
	createdAt: string,
	userName: string,
	userDisplayName: string,
	message: string
}>();

</script>

<template>
	<Base :icon="IconMoodPlus" :created-at="createdAt">
		<template #leftContent>
			<div style="display: flex; flex-direction:column">
				<span>{{ createUserName(userName, userDisplayName) }} <b style="font-size: 11px">first time joined chat</b></span>
				<span style="font-size: 11px">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
