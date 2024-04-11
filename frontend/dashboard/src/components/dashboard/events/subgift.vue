<script lang="ts" setup>
import { IconGift } from '@tabler/icons-vue';

import Base from './base.vue';

import { resolveUserName } from '@/helpers/index.js';

defineProps<{
	userName: string,
	userDisplayName: string,
	targetUserName: string,
	targetUserDisplayName: string,
	createdAt: string,
	level: string,
}>();

</script>

<template>
	<Base :icon="IconGift" :created-at="createdAt">
		<template #leftContent>
			<span>
				{{ resolveUserName(userName, userDisplayName) }} <b class="text-xs">gifted sub to</b> {{ resolveUserName(targetUserName, targetUserDisplayName) }}
			</span>
			<span>
				<span class="text-xs">
					{{ level }} tier
				</span>
			</span>
		</template>
	</Base>
</template>
