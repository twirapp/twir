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
			<div style="display: flex; flex-direction: column">
				<span>{{ userName }} <b style="font-size: 11px">donated</b> {{ amount }}{{ currency }}</span>
				<span style="font-size: 11px">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
