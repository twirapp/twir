<script lang="ts" setup>
import type { AcceptableValue } from 'reka-ui'

import { useTwitchGetCategories, useTwitchSearchCategories } from '#layers/dashboard/api/twitch'

defineProps<{ id?: string }>()

const categories = defineModel<string[]>({ default: [] })

const categoriesSearch = ref('')
const categoriesSearchDebounced = refDebounced(categoriesSearch, 500)

const { data: searchCategoriesData } = useTwitchSearchCategories(categoriesSearchDebounced)

const { data: selectedCategories } = useTwitchGetCategories(categories)

interface SelectedCategoryValue {
	id: string
	label: string
	image: string
}

const selectedCategoriesValues = computed<Record<string, SelectedCategoryValue>>(() => {
	if (!selectedCategories.value) return {}

	return selectedCategories.value.reduce(
		(acc, val) => {
			acc[val.id] = {
				id: val.id,
				image: val.boxArtUrl.replace('{height}', '80').replace('{width}', '60'),
				label: val.name,
			}

			return acc
		},
		{} as Record<string, SelectedCategoryValue>
	)
})

function handleSelect(
	event: CustomEvent<{
		originalEvent: PointerEvent
		value?: AcceptableValue
	}>
) {
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
	<UiPopover :open="!!searchCategoriesData?.length">
		<UiPopoverTrigger as-child>
			<UiTagsInput v-model="categories">
				<UiTagsInputItem
					v-for="item in selectedCategoriesValues"
					:key="item.label"
					:value="item.id"
					class="flex gap-1 items-center rounded-full px-2"
				>
					<img :src="item.image" class="size-4 rounded-full" />
					{{ item.label }}
					<TagsInputItemDelete />
				</UiTagsInputItem>

				<input
					:id="id"
					v-model="categoriesSearch"
					type="text"
					placeholder="Search..."
					class="text-sm min-h-6 focus:outline-hidden flex-1 bg-transparent px-1"
				/>
			</UiTagsInput>
		</UiPopoverTrigger>
		<UiPopoverContent class="p-0">
			<UiCommand>
				<UiCommandList>
					<UiCommandGroup>
						<UiCommandItem
							v-for="option in searchCategoriesData"
							:key="option.id"
							:value="option.id"
							class="flex gap-2.5 h-24 items-center"
							@select="handleSelect"
						>
							<img
								:src="option.boxArtUrl.replace('{width}', '60').replace('{height}', '80')"
								class="h-[80px] w-[60px]"
							/>
							<span>{{ option.name }}</span>
						</UiCommandItem>
					</UiCommandGroup>
				</UiCommandList>
			</UiCommand>
		</UiPopoverContent>
	</UiPopover>
</template>
