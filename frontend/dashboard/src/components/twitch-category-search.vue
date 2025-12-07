<script setup lang="ts">
import { refDebounced } from '@vueuse/core'
import { NSelect } from 'naive-ui'
import { computed, h, ref, watch } from 'vue'

import type { SelectOption } from 'naive-ui'
import type { VNodeChild } from 'vue'

import { useTwitchGetCategories, useTwitchSearchCategories } from '@/api'

defineProps<{
	multiple?: boolean
}>()
const category = defineModel<undefined | null | string | string[]>()

const categoriesSearch = ref('')
const categoriesSearchDebounced = refDebounced(categoriesSearch, 500)

const {
	data: searchCategoriesData,
	isLoading: isSearchCategoriesLoading,
} = useTwitchSearchCategories(categoriesSearchDebounced)

const getCategoriesRef = ref<string[]>([])
watch(() => category.value, (v) => {
	if (!v) return []
	getCategoriesRef.value = Array.isArray(v) ? v : [v]
}, { immediate: true, once: true })
const {
	data: getCategoriesData,
	isLoading: isGetCategoriesLoading,
} = useTwitchGetCategories(getCategoriesRef)

const categoriesOptions = computed(() => {
	return [
		...searchCategoriesData.value.map((c) => ({
			label: c.name,
			value: c.id,
			image: c.boxArtUrl.replace('{width}', '144').replace('{height}', '192'),
		})),
		...getCategoriesData.value.map((c) => ({
			label: c.name,
			value: c.id,
			image: c.boxArtUrl.replace('{width}', '144').replace('{height}', '192'),
		})),
	]
})

function renderCategory(o: SelectOption & { image?: string }): VNodeChild {
	return [h(
		'div',
		{ class: 'flex gap-2.5 h-24 items-center' },
		[
			h('img', {
				src: o.image,
				style: { height: '80px', width: '60px' },
			}),
			h('span', {}, o.label! as string),
		],
	)]
}
</script>

<template>
	<NSelect
		v-model:value="category"
		filterable
		placeholder="Search for category..."
		:options="categoriesOptions"
		remote
		:multiple
		:render-label="renderCategory"
		:loading="isSearchCategoriesLoading || isGetCategoriesLoading"
		:render-tag="(t) => t.option.label as string ?? ''"
		clearable
		@search="(v) => categoriesSearch = v"
	/>
</template>
