<script lang="ts" setup>
import { IconAward } from '@tabler/icons-vue';

import Base from './base.vue';

import { createUserName } from '@/helpers/index.js';

defineProps<{
	createdAt: string,
	userName: string,
	userDisplayName: string,
	title: string,
	input: string,
	cost: string
}>();
</script>

<template>
	<Base :icon="IconAward" :created-at="createdAt">
		<template #leftContent>
			<div class="flex flex-col">
				<span class="break-all">{{ createUserName(userName, userDisplayName) }}
					<b class="text-xs">activated</b> {{ title }}({{ cost }})
				</span>
				<span v-if="input.length" class="text-xs">{{ input }}</span>
			</div>
		</template>
	</Base>
</template>
