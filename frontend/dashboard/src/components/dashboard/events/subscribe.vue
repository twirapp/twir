<script lang="ts" setup>
import { IconStarFilled } from '@tabler/icons-vue';

import Base from './base.vue';

import { resolveUserName } from '@/helpers/index.js';

defineProps<{
	userName: string,
	userDisplayName: string,
	createdAt: string,
	level: string,
}>();

</script>

<template>
	<Base :icon="IconStarFilled" :created-at="createdAt">
		<template #leftContent>
			<span>
				{{ resolveUserName(userName, userDisplayName) }} <b class="text-xs">subscribed with</b> {{ level }} tier
			</span>
		</template>
	</Base>
</template>
