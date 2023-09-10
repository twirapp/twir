<script lang="ts" setup>
import { IconEdit } from '@tabler/icons-vue';
import { refDebounced } from '@vueuse/core';
import {
	type SelectOption,
	NCard,
	NModal,
	NForm,
	NFormItem,
	NInput,
	NButton,
	NSelect,
	NAvatar,
	useMessage,
useThemeVars,
 } from 'naive-ui';
import { VNodeChild, computed, ref, h } from 'vue';
import { useI18n } from 'vue-i18n';

import { twitchSetChannelInformationMutation, useTwitchSearchCategories } from '@/api/index.js';

const { t } = useI18n();

const theme = useThemeVars();

const props = defineProps<{
	title?: string,
	categoryId?: string,
	categoryName?: string,
}>();

const isEditInformationModalShowed = ref(false);

const form = ref({
	title: '',
	categoryId: '',
});

const categoriesSearch = ref('');
const categoriesSearchDebounced = refDebounced(categoriesSearch, 500);

const openEditInformationModalModal = () => {
	isEditInformationModalShowed.value = true;
	form.value = {
		categoryId: props.categoryId ?? '',
		title: props.title ?? '',
	};
	categoriesSearch.value = props.categoryName || '';
};

const { data: categoriesData, isLoading: isCategoriesLoading } = useTwitchSearchCategories(categoriesSearchDebounced);

const categoriesOptions = computed(() => {
	return categoriesData.value?.categories.map((c) => ({ label: c.name, value: c.id, image: c.image }));
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
			h(NAvatar, {
				src: o.image?.replace('52x72', '1280x720'),
				style: { height: '90px', width: '90px', marginLeft: '10px' },
			}),
			h('span', {}, o.label! as string),
		],
	)];
};

const informationUpdater = twitchSetChannelInformationMutation();

const messages = useMessage();
async function saveChannelInformation() {
	const { title, categoryId } = form.value;
	await informationUpdater.mutateAsync({ categoryId, title });
	isEditInformationModalShowed.value = false;
	messages.success(t('sharedTexts.saved'));
}

</script>

<template>
	<div style="display: flex; flex-direction: row; flex-wrap: wrap; gap: 5px; height: 100%">
		<n-card
			class="card"
			:bordered="false"
			embedded
			content-style="padding: 5px;"
			:style="{ 'background-color': theme.actionColor }"
		>
			<div style="display: flex; justify-content: space-between; align-items: center">
				<div style="display: flex; font-size:15px; flex-direction: column; white-space: nowrap; overflow: hidden;text-overflow: ellipsis">
					<span style="font-size:15px;">
						{{ props?.title || t('dashboard.statsWidgets.streamInfo.noTitle') }}
					</span>
					<span style="font-size:15px;">
						{{ props?.categoryName || t('dashboard.statsWidgets.streamInfo.noCategory') }}
					</span>
				</div>

				<IconEdit style="display: flex; width: 35px; height: 35px" @click="openEditInformationModalModal" />
			</div>
		</n-card>
	</div>

	<n-modal
		v-model:show="isEditInformationModalShowed"
		preset="card"
		:bordered="false"
		:segmented="true"
		style="width: 500px"
		:title="t('dashboard.statsWidgets.streamInfo.modalTitle')"
	>
		<n-form>
			<n-form-item :label="t('dashboard.statsWidgets.streamInfo.title')">
				<n-input v-model:value="form.title" :placeholder="t('dashboard.statsWidgets.streamInfo.title')" />
			</n-form-item>

			<n-form-item :label="t('dashboard.statsWidgets.streamInfo.category')">
				<n-select
					v-model:value="form.categoryId"
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

<style scoped>
.card {
	flex: 1 1 200px;
	cursor: pointer;
}
</style>
