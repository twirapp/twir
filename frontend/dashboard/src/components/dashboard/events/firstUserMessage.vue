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
			<div class="flex flex-col">
				<span>
					{{ createUserName(userName, userDisplayName) }}
					<b class="text-sx">first time joined chat</b>
				</span>
				<span class="text-sx">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
