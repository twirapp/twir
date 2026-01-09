<script setup lang="ts">
import { useChatMessagesFilters } from '../composables/use-filters'




const filters = useChatMessagesFilters()
</script>

<template>
	<div>
		<div class="flex flex-row flex-wrap gap-4 justify-between">
			<div class="flex flex-col gap-2 w-[48%]">
				<UiLabel for="username">
					Username
				</UiLabel>
				<UiInput id="username" v-model="filters.userSearchInput.value" placeholder="Search by username" />
			</div>
			<div class="flex flex-col gap-2 w-[48%]">
				<UiLabel for="text">
					Text
				</UiLabel>
				<UiInput id="text" v-model="filters.textSearchInput.value" placeholder="Search by text" />
			</div>
		</div>
	</div>
</template>
