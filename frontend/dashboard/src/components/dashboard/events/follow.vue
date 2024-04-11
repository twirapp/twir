<script lang="ts" setup>
import { IconHeartFilled } from '@tabler/icons-vue';

import Base from './base.vue';

import { resolveUserName } from '@/helpers/index.js';

defineProps<{
	userName: string,
	userDisplayName: string,
	createdAt: string
}>();

</script>

<template>
	<Base :icon="IconHeartFilled" :created-at="createdAt">
		<template #leftContent>
			{{ resolveUserName(userName, userDisplayName) }}
		</template>
	</Base>
</template>
