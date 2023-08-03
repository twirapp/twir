<script lang="ts" setup>
import { IconStarFilled } from '@tabler/icons-vue';

import Base from './base.vue';

import { createUserName } from '@/helpers/index.js';

defineProps<{
	userName: string,
	userDisplayName: string,
	createdAt: string,
	level: string,
}>();

</script>

<template>
	<Base :icon="IconStarFilled" icon-color="green" :created-at="createdAt">
		<template #leftContent>
			<span>
				{{ createUserName(userName, userDisplayName) }} <b style="font-size: 11px">subscribed with</b> {{ level }} tier
			</span>
		</template>
	</Base>
</template>
