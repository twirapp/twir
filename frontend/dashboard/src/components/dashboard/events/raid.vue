<script lang="ts" setup>
import { IconPick } from '@tabler/icons-vue';

import Base from './base.vue';

import { createUserName } from '@/helpers/index.js';

defineProps<{
	userName: string,
	userDisplayName: string,
	createdAt: string,
	viewers: string,
}>();

</script>

<template>
	<Base :icon="IconPick" :created-at="createdAt">
		<template #leftContent>
			<span>

				{{ createUserName(userName, userDisplayName) }} <b class="text-xs">raided with</b> {{ viewers }} viewers
			</span>
		</template>
	</Base>
</template>
