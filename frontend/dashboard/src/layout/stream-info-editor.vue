<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'

import {
	useUserAccessFlagChecker,
} from '@/api/auth'
import { twitchSetChannelInformationMutation } from '@/api/twitch'
import TwitchCategorySelector from '@/components/twitch-category-selector.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = defineProps<{
	title?: string
	categoryId?: string
	categoryName?: string
}>()

const open = defineModel<boolean>('open', { default: false })

const { t } = useI18n()

const form = ref({
	title: '',
	categoryId: undefined as string | undefined,
})

watch(() => props, (v) => {
	form.value = {
		title: v.title ?? '',
		categoryId: v.categoryId ?? undefined,
	}
}, { immediate: true, deep: true })

const informationUpdater = twitchSetChannelInformationMutation()

async function saveChannelInformation() {
	try {
		await informationUpdater.executeMutation({
			categoryId: form.value.categoryId,
			title: form.value.title,
		})
		toast.success(t('sharedTexts.saved'))
		open.value = false
	} catch (error) {
		toast.error(error instanceof Error ? error.message : 'Failed to save')
	}
}

const userCanEditTitle = useUserAccessFlagChecker(ChannelRolePermissionEnum.UpdateChannelTitle)
const userCanEditCategory = useUserAccessFlagChecker(ChannelRolePermissionEnum.UpdateChannelCategory)
</script>

<template>
	<Dialog v-model:open="open">
		<DialogContent class="sm:max-w-[500px]">
			<DialogHeader>
				<DialogTitle>{{ t('dashboard.statsWidgets.streamInfo.editStreamInfo') }}</DialogTitle>
				<DialogDescription>
					{{ t('dashboard.statsWidgets.streamInfo.editStreamInfoDescription') }}
				</DialogDescription>
			</DialogHeader>

			<div class="grid gap-4 py-4">
				<div class="grid gap-2">
					<Label for="title">
						{{ t('dashboard.statsWidgets.streamInfo.title') }}
					</Label>
					<Input
						id="title"
						v-model="form.title"
						:disabled="!userCanEditTitle"
						:placeholder="t('dashboard.statsWidgets.streamInfo.title')"
					/>
				</div>

				<div class="grid gap-2">
					<Label for="category">
						{{ t('dashboard.statsWidgets.streamInfo.category') }}
					</Label>
					<TwitchCategorySelector
						id="category"
						v-model="form.categoryId"
						:disabled="!userCanEditCategory"
					/>
				</div>
			</div>

			<DialogFooter>
				<Button
					type="button"
					variant="outline"
					@click="open = false"
				>
					{{ t('sharedButtons.cancel') }}
				</Button>
				<Button
					type="submit"
					:disabled="informationUpdater.fetching.value"
					@click="saveChannelInformation"
				>
					{{ informationUpdater.fetching.value ? t('sharedButtons.saving') : t('sharedButtons.save') }}
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
