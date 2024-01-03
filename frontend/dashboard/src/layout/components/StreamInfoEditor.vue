<script setup lang="ts">
import { refDebounced } from '@vueuse/core';
import {
	NButton,
	NForm,
	NFormItem,
	NInput,
	NModal,
	NSelect,
	type SelectOption,
	useNotification,
} from 'naive-ui';
import { computed, h, ref, VNodeChild, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
	twitchSetChannelInformationMutation,
	useTwitchSearchCategories,
	useUserAccessFlagChecker,
} from '@/api';

const { t } = useI18n();

const props = defineProps<{
	title?: string,
	categoryId?: string,
	categoryName?: string,
	opened: boolean,
}>();

const emits = defineEmits<{
	close: [],
}>();

const form = ref({
	title: '',
	categoryId: '',
});


const categoriesSearch = ref('');
const categoriesSearchDebounced = refDebounced(categoriesSearch, 500);

watch(props, (v) => {
	form.value = {
		title: v.title ?? '',
		categoryId: v.categoryId ?? '',
	};
	categoriesSearch.value = v.categoryName ?? '';
}, { immediate: true });

const {
	data: categoriesData,
	isLoading: isCategoriesLoading,
} = useTwitchSearchCategories(categoriesSearchDebounced);

const categoriesOptions = computed(() => {
	return categoriesData.value?.categories.map((c) => ({
		label: c.name,
		value: c.id,
		image: c.image,
	}));
});

const renderCategory = (o: SelectOption & { image?: string }): VNodeChild => {
	return [h(
		'div',
		{
			style: {
				display: 'flex',
				alignItems: 'center',
				height: '100px',
				gap: '10px',
			},
		},
		[
			h('img', {
				src: o.image?.replace('52x72', '144x192'),
				style: { height: '80px', width: '60px' },
			}),
			h('span', {}, o.label! as string),
		],
	)];
};

const informationUpdater = twitchSetChannelInformationMutation();

const messages = useNotification();

async function saveChannelInformation() {
	await informationUpdater.mutateAsync({
		categoryId: form.value.categoryId,
		title: form.value.title,
	});
	messages.success({
		title: t('sharedTexts.saved'),
		duration: 2500,
	});
	emits('close');
}

const userCanEditTitle = useUserAccessFlagChecker('UPDATE_CHANNEL_TITLE');
const userCanEditCategory = useUserAccessFlagChecker('UPDATE_CHANNEL_CATEGORY');
</script>

<template>
	<n-modal
		:show="opened"
		preset="card"
		:bordered="false"
		:segmented="true"
		style="width: 500px"
		:title="t('dashboard.statsWidgets.streamInfo.modalTitle')"
		@close="() => emits('close')"
	>
		<n-form>
			<n-form-item :label="t('dashboard.statsWidgets.streamInfo.title')">
				<n-input
					v-model:value="form.title"
					:disabled="!userCanEditTitle"
					:placeholder="t('dashboard.statsWidgets.streamInfo.title')"
				/>
			</n-form-item>

			<n-form-item :label="t('dashboard.statsWidgets.streamInfo.category')">
				<n-select
					v-model:value="form.categoryId"
					:disabled="!userCanEditCategory"
					filterable
					placeholder="Search..."
					:options="categoriesOptions"
					remote
					:render-label="renderCategory"
					:loading="isCategoriesLoading"
					:render-tag="(t) => t.option.label as string ?? ''"
					@search="(v) => categoriesSearch = v"
				/>
			</n-form-item>

			<n-button secondary block type="success" @click="saveChannelInformation">
				{{ t('sharedButtons.save') }}
			</n-button>
		</n-form>
	</n-modal>
</template>
