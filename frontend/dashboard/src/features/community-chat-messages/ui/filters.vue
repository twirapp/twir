<script setup lang="ts">
import { useChatMessagesFilters } from '../composables/use-filters'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const filters = useChatMessagesFilters()
</script>

<template>
	<div>
		<div class="flex flex-row flex-wrap gap-4 justify-between">
			<div class="flex flex-col gap-2 w-[48%]">
				<Label for="username">
					Username
				</Label>
				<Input id="username" v-model="filters.userSearchInput.value" placeholder="Search by username" />
			</div>
			<div class="flex flex-col gap-2 w-[48%]">
				<Label for="text">
					Text
				</Label>
				<Input id="text" v-model="filters.textSearchInput.value" placeholder="Search by text" />
			</div>
		</div>
	</div>
</template>
