<script lang="ts" setup>
import { IconCurrencyDollar } from '@tabler/icons-vue'

import Base from './base.vue'

defineProps<{
	userName: string
	createdAt: string
	amount: string
	currency: string
	message: string
}>()
</script>

<template>
	<Base
		:icon="IconCurrencyDollar"
		:icon-color="['#00a865', '#5cffbe']"
		:created-at="createdAt"
	>
		<template #leftContent>
			<div class="flex flex-col">
				<span>{{ userName }} <span class="font-bold">donated</span> {{ amount }}{{ currency }}</span>
				<span class="text-xs">{{ message }}</span>
			</div>
		</template>
	</Base>
</template>
