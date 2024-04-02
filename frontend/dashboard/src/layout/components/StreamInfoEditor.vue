<script setup lang="ts">
import {
	NButton,
	NForm,
	NFormItem,
	NInput,
} from 'naive-ui';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import {
	twitchSetChannelInformationMutation,
	useUserAccessFlagChecker,
} from '@/api';
import TwitchCategorySearch from '@/components/twitch-category-search.vue';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';

const { t } = useI18n();

const props = defineProps<{
	title?: string,
	categoryId?: string,
	categoryName?: string,
}>();

const form = ref({
	title: '',
	categoryId: '',
});


watch(props, (v) => {
	form.value = {
		title: v.title ?? '',
		categoryId: v.categoryId ?? '',
	};
}, { immediate: true });

const informationUpdater = twitchSetChannelInformationMutation();

const discrete = useNaiveDiscrete();

async function saveChannelInformation() {
	await informationUpdater.mutateAsync({
		categoryId: form.value.categoryId,
		title: form.value.title,
	});
	discrete.notification.success({
		title: t('sharedTexts.saved'),
		duration: 2500,
	});
	discrete.dialog.destroyAll();
}

const userCanEditTitle = useUserAccessFlagChecker('UPDATE_CHANNEL_TITLE');
const userCanEditCategory = useUserAccessFlagChecker('UPDATE_CHANNEL_CATEGORY');
</script>

<template>
	<n-form>
		<n-form-item :label="t('dashboard.statsWidgets.streamInfo.title')">
			<n-input
				v-model:value="form.title"
				:disabled="!userCanEditTitle"
				:placeholder="t('dashboard.statsWidgets.streamInfo.title')"
			/>
		</n-form-item>

		<n-form-item :label="t('dashboard.statsWidgets.streamInfo.category')">
			<twitch-category-search v-model="form.categoryId" :disabled="!userCanEditCategory" />
		</n-form-item>

		<n-button secondary block type="success" @click="saveChannelInformation">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-form>
</template>
