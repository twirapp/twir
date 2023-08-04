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
			<div style="display: flex; flex-direction:column">
				<span style="word-break:break-all;">{{ createUserName(userName, userDisplayName) }}
					<b style="font-size: 11px">activated</b> {{ title }}({{ cost }})
				</span>
				<span v-if="input.length" style="font-size: 11px">{{ input }}</span>
			</div>
		</template>
	</Base>
</template>
