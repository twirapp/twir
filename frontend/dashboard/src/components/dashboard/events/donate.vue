<script lang="ts" setup>
import { IconCurrencyDollar } from '@tabler/icons-vue';

import Base from './base.vue';

defineProps<{
	userName: string,
	createdAt: string,
	amount: string,
	currency: string,
	message: string
}>();
</script>

<template>
	<Base :icon="IconCurrencyDollar" :created-at="createdAt">
		<template #leftContent>
			<div class="flex flex-col">
				<span>{{ userName }} <b class="text-xs">donated</b> {{ amount }}{{ currency }}</span>
				<span class="text-xs">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
