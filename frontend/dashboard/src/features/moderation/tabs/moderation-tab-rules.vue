<script setup lang="ts">
import { useRouter } from 'vue-router'

import Card from '../ui/card.vue'

import { useModerationApi } from '@/features/moderation/composables/use-moderation-api.ts'

const { items } = useModerationApi()
const router = useRouter()

function showForm(itemId: string) {
	router.push({ name: 'ModerationForm', params: { id: itemId } })
}
</script>

<template>
	<div>
		<div v-if="!items.length">
			<div class="flex flex-col gap-2 items-center justify-center h-full">
				<h2 class="text-2xl font-bold">
					No rules
				</h2>
				<p class="text-sm">
					Create a new rule to start moderating your chat.
				</p>
			</div>
		</div>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<Card
				v-for="item of items"
				:key="item.id"
				:item="item"
				@show-settings="showForm(item.id)"
			/>
		</div>
	</div>
</template>
