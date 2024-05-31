<script setup lang="ts">
import {
	NButton,
	NForm,
	NFormItem,
	NInput,
} from 'naive-ui'
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import {
	twitchSetChannelInformationMutation,
	useUserAccessFlagChecker,
} from '@/api'
import TwitchCategorySearch from '@/components/twitch-category-search.vue'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = defineProps<{
	title?: string
	categoryId?: string
	categoryName?: string
}>()

const { t } = useI18n()

const form = ref({
	title: '',
	categoryId: '',
})

watch(props, (v) => {
	form.value = {
		title: v.title ?? '',
		categoryId: v.categoryId ?? '',
	}
}, { immediate: true })

const informationUpdater = twitchSetChannelInformationMutation()

const discrete = useNaiveDiscrete()

async function saveChannelInformation() {
	await informationUpdater.mutateAsync({
		categoryId: form.value.categoryId,
		title: form.value.title,
	})
	discrete.notification.success({
		title: t('sharedTexts.saved'),
		duration: 2500,
	})
	discrete.dialog.destroyAll()
}

const userCanEditTitle = useUserAccessFlagChecker(ChannelRolePermissionEnum.UpdateChannelTitle)
const userCanEditCategory = useUserAccessFlagChecker(ChannelRolePermissionEnum.UpdateChannelCategory)
</script>

<template>
	<NForm>
		<NFormItem :label="t('dashboard.statsWidgets.streamInfo.title')">
			<NInput
				v-model:value="form.title"
				:disabled="!userCanEditTitle"
				:placeholder="t('dashboard.statsWidgets.streamInfo.title')"
			/>
		</NFormItem>

		<NFormItem :label="t('dashboard.statsWidgets.streamInfo.category')">
			<TwitchCategorySearch v-model="form.categoryId" :disabled="!userCanEditCategory" />
		</NFormItem>

		<NButton secondary block type="success" @click="saveChannelInformation">
			{{ t('sharedButtons.save') }}
		</NButton>
	</NForm>
</template>
