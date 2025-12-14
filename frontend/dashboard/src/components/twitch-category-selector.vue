<script setup lang="ts">
import { Check, ChevronsUpDown } from 'lucide-vue-next'
import { refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'

import { useTwitchGetCategories, useTwitchSearchCategories } from '@/api'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '@/lib/utils'

withDefaults(
	defineProps<{
		disabled?: boolean
		placeholder?: string
	}>(),
	{
		disabled: false,
		placeholder: 'Select category...',
	}
)

const categoryId = defineModel<string | undefined>()

const open = ref(false)
const searchQuery = ref('')
const searchQueryDebounced = refDebounced(searchQuery, 300)

const { data: searchResults, isLoading: isSearching } =
	useTwitchSearchCategories(searchQueryDebounced)

// Load current category if categoryId is set
const initialIds = computed(() => {
	if (!categoryId.value) return []
	return [categoryId.value]
})

const { data: selectedCategories } = useTwitchGetCategories(initialIds)

const selectedCategory = computed(() => {
	if (!categoryId.value) return null
	return selectedCategories.value.find((c) => c.id === categoryId.value)
})

const displayedCategories = computed(() => {
	// Combine search results with selected category
	const categories = [...searchResults.value]

	if (selectedCategory.value && !categories.find((c) => c.id === selectedCategory.value!.id)) {
		categories.unshift(selectedCategory.value)
	}

	return categories
})

function selectCategory(category: (typeof searchResults.value)[0] | null) {
	categoryId.value = category?.id
	open.value = false
	searchQuery.value = ''
}
</script>

<template>
	<Popover v-model:open="open">
		<PopoverTrigger as-child>
			<Button
				variant="outline"
				role="combobox"
				:aria-expanded="open"
				:disabled="disabled"
				class="w-full justify-between"
			>
				<div v-if="selectedCategory" class="flex items-center gap-2">
					<img
						:src="selectedCategory.boxArtUrl.replace('{width}', '52').replace('{height}', '72')"
						:alt="selectedCategory.name"
						class="h-8 w-6 object-cover rounded"
					/>
					<span>{{ selectedCategory.name }}</span>
				</div>
				<span v-else class="text-muted-foreground">{{ placeholder }}</span>
				<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
			</Button>
		</PopoverTrigger>
		<PopoverContent class="w-[400px] p-0">
			<Command :filter-function="(list: any) => list">
				<div class="flex items-center border-b px-3">
					<input
						v-model="searchQuery"
						type="text"
						placeholder="Search categories..."
						class="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-hidden placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50"
					/>
				</div>
				<CommandList>
					<CommandEmpty v-if="!isSearching && displayedCategories.length === 0">
						{{ searchQuery ? 'No category found.' : 'Start typing to search...' }}
					</CommandEmpty>
					<CommandEmpty v-else-if="isSearching"> Searching... </CommandEmpty>
					<CommandGroup v-if="displayedCategories.length > 0">
						<CommandItem
							v-for="category in displayedCategories"
							:key="category.id"
							:value="category.id"
							@select="() => selectCategory(category)"
						>
							<div class="flex items-center gap-3 w-full">
								<img
									:src="category.boxArtUrl.replace('{width}', '52').replace('{height}', '72')"
									:alt="category.name"
									class="h-14 w-10 object-cover rounded shrink-0"
								/>
								<span class="flex-1 truncate">{{ category.name }}</span>
								<Check
									:class="cn('h-4 w-4', categoryId === category.id ? 'opacity-100' : 'opacity-0')"
								/>
							</div>
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
