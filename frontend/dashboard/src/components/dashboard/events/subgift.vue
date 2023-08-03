<script lang="ts" setup>
import { IconGift } from '@tabler/icons-vue';

import Base from './base.vue';

import { createUserName } from '@/helpers/index.js';

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
	<Base :icon="IconGift" icon-color="green" :created-at="createdAt">
		<template #leftContent>
			<span>
				{{ createUserName(userName, userDisplayName) }} <b style="font-size: 11px">gifted sub to</b> {{ createUserName(targetUserName, targetUserDisplayName) }}
			</span>
			<span>
				<span style="font-size: 11px">
					{{ level }} tier
				</span>
			</span>
		</template>
	</Base>
</template>
