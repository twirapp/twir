<script lang="ts" setup>
import { IconRepeat } from '@tabler/icons-vue';

import Base from './base.vue';

import { createUserName } from '@/helpers/index.js';

defineProps<{
	userName: string,
	userDisplayName: string,
	createdAt: string,
	level: string,
	months: string,
	streak: string,
}>();

</script>

<template>
	<Base :icon="IconRepeat" :created-at="createdAt">
		<template #leftContent>
			<span>
				{{ createUserName(userName, userDisplayName) }} <b class="text-xs">resubscribed with</b> {{ level }} tier
			</span>
			<span>
				<span class="text-xs">
					{{ streak }} months • overall {{ months }} months
				</span>
			</span>
		</template>
	</Base>
</template>
