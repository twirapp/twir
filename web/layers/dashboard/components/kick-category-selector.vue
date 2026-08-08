<script setup lang="ts">
import { refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'

import { useKickSearchCategories } from '~~/layers/dashboard/api/kick.js'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { cn } from '~~/layers/dashboard/lib/utils'

const props = withDefaults(
	defineProps<{
		categoryName?: string
		disabled?: boolean
		placeholder?: string
	}>(),
	{
		categoryName: '',
		disabled: false,
		placeholder: 'Select category...',
	}
)

const categoryId = defineModel<string | undefined>()

const open = ref(false)
const searchQuery = ref('')
const searchQueryDebounced = refDebounced(searchQuery, 300)
const { data: searchResults, isLoading: isSearching } = useKickSearchCategories(searchQueryDebounced)

interface Category {
	readonly id: string
	readonly name: string
	readonly thumbnail: string
}

const selectedCategory = computed<Category | null>(() => {
	if (!categoryId.value) return null
	return (
		searchResults.value.find((category) => category.id === categoryId.value) ??
		(props.categoryName
			? { id: categoryId.value, name: props.categoryName, thumbnail: '' }
			: null)
	)
})

const displayedCategories = computed<readonly Category[]>(() => {
	const categories = [...searchResults.value]
	if (selectedCategory.value && !categories.some((category) => category.id === selectedCategory.value?.id)) {
		categories.unshift(selectedCategory.value)
	}
	return categories
})

function selectCategory(category: Category | null) {
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
				:disabled="props.disabled"
				class="w-full justify-between"
			>
				<div v-if="selectedCategory" class="flex items-center gap-2">
					<img
						v-if="selectedCategory.thumbnail"
						:src="selectedCategory.thumbnail"
						:alt="selectedCategory.name"
						class="h-8 w-6 rounded object-cover"
					/>
					<span>{{ selectedCategory.name }}</span>
				</div>
				<span v-else class="text-muted-foreground">{{ props.placeholder }}</span>
				<Icon name="lucide:chevrons-up-down" class="ml-2 h-4 w-4 shrink-0 opacity-50" />
			</Button>
		</PopoverTrigger>
		<PopoverContent class="w-[400px] p-0">
			<Command>
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
					<CommandEmpty v-else-if="isSearching">Searching...</CommandEmpty>
					<CommandGroup v-if="displayedCategories.length > 0">
						<CommandItem
							v-for="category in displayedCategories"
							:key="category.id"
							:value="category.id"
							@select="() => selectCategory(category)"
						>
							<div class="flex w-full items-center gap-3">
								<img
									v-if="category.thumbnail"
									:src="category.thumbnail"
									:alt="category.name"
									class="h-14 w-10 shrink-0 rounded object-cover"
								/>
								<span class="flex-1 truncate">{{ category.name }}</span>
								<Icon
									name="lucide:check"
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
