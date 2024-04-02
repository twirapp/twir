<script setup lang="ts">
import { refDebounced } from '@vueuse/core';
import { NSelect, SelectOption } from 'naive-ui';
import { computed, h, ref, VNodeChild, watch } from 'vue';

import { useTwitchGetCategories, useTwitchSearchCategories } from '@/api';

defineProps<{
	multiple?: boolean
}>();
const category = defineModel<string | string[]>();

const categoriesSearch = ref('');
const categoriesSearchDebounced = refDebounced(categoriesSearch, 500);

const {
	data: searchCategoriesData,
	isLoading: isSearchCategoriesLoading,
} = useTwitchSearchCategories(categoriesSearchDebounced);

const getCategoriesRef = ref<string[]>([]);
watch(() => category.value, (v) => {
	if (!v) return [];
	getCategoriesRef.value = Array.isArray(v) ? v : [v];
}, { immediate: true, once: true });
const {
	data: getCategoriesData,
	isLoading: isGetCategoriesLoading,
} = useTwitchGetCategories(getCategoriesRef);

const categoriesOptions = computed(() => {
	return [
		...searchCategoriesData.value?.categories.map((c) => ({
			label: c.name,
			value: c.id,
			image: c.image.replace('52x72', '144x192'),
		})) ?? [],
		...getCategoriesData.value?.categories.map((c) => ({
			label: c.name,
			value: c.id,
			image: c.image.replace('{width}', '144').replace('{height}', '192'),
		})) ?? [],
	];
});

const renderCategory = (o: SelectOption & { image?: string }): VNodeChild => {
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
	)];
};

</script>

<template>
	<n-select
		v-model:value="category"
		filterable
		placeholder="Search..."
		:options="categoriesOptions"
		remote
		:multiple
		:render-label="renderCategory"
		:loading="isSearchCategoriesLoading || isGetCategoriesLoading"
		:render-tag="(t) => t.option.label as string ?? ''"
		@search="(v) => categoriesSearch = v"
	/>
</template>
