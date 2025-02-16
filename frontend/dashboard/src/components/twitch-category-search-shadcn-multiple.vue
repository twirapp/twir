<script lang="ts" setup>
import { refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'

import Command from './ui/command/Command.vue'
import CommandGroup from './ui/command/CommandGroup.vue'
import CommandItem from './ui/command/CommandItem.vue'
import CommandList from './ui/command/CommandList.vue'
import Popover from './ui/popover/Popover.vue'
import PopoverContent from './ui/popover/PopoverContent.vue'
import PopoverTrigger from './ui/popover/PopoverTrigger.vue'
import TagsInput from './ui/tags-input/TagsInput.vue'
import TagsInputItem from './ui/tags-input/TagsInputItem.vue'
import TagsInputItemDelete from './ui/tags-input/TagsInputItemDelete.vue'

import { useTwitchGetCategories, useTwitchSearchCategories } from '@/api'

defineProps<{ id?: string }>()

const categories = defineModel<string[]>({ default: [] })

const categoriesSearch = ref('')
const categoriesSearchDebounced = refDebounced(categoriesSearch, 500)

const {
	data: searchCategoriesData,
} = useTwitchSearchCategories(categoriesSearchDebounced)

const {
	data: selectedCategories,
} = useTwitchGetCategories(categories, { keepPreviousData: true })

interface SelectedCategoryValue {
	id: string
	label: string
	image: string
}

const selectedCategoriesValues = computed<Record<string, SelectedCategoryValue>>(() => {
	if (!selectedCategories.value) return {}

	return selectedCategories.value.categories.reduce((acc, val) => {
		acc[val.id] = {
			id: val.id,
			image: val.image.replace('{height}', '80').replace('{width}', '60'),
			label: val.name,
		}

		return acc
	}, {} as Record<string, SelectedCategoryValue>)
})

function handleSelect(event: CustomEvent<{
	originalEvent: PointerEvent
	value?: string | number | boolean | Record<string, any>
}>) {
	if (typeof event.detail.value !== 'string') return
	if (categories.value.includes(event.detail.value)) {
		return
	} else {
		categories.value?.push(event.detail.value)
	}

	categoriesSearch.value = ''
}
</script>

<template>
	<Popover :open="!!searchCategoriesData?.categories.length">
		<PopoverTrigger as-child>
			<TagsInput v-model="categories">
				<TagsInputItem v-for="item in selectedCategoriesValues" :key="item.label" :value="item.id" class="flex gap-1 items-center rounded-full px-2">
					<img :src="item.image" class="size-4 rounded-full" />
					{{ item.label }}
					<TagsInputItemDelete />
				</TagsInputItem>

				<input
					:id="id"
					v-model="categoriesSearch"
					type="text"
					placeholder="Search..."
					class="text-sm min-h-6 focus:outline-none flex-1 bg-transparent px-1"
				/>
			</TagsInput>
		</PopoverTrigger>
		<PopoverContent class="p-0">
			<Command>
				<CommandList>
					<CommandGroup>
						<CommandItem
							v-for="option in searchCategoriesData?.categories"
							:key="option.id"
							:value="option.id"
							class="flex gap-2.5 h-24 items-center"
							@select="handleSelect"
						>
							<img :src="option.image" class="h-[80px] w-[60px]" />
							<span>{{ option.name }}</span>
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
