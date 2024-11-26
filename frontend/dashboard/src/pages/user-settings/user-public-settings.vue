<script setup lang="ts">
import { ArrowDownIcon, ArrowUpIcon, EditIcon, TrashIcon } from 'lucide-vue-next'
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NButton,
	NCard,
	NForm,
	NFormItem,
	NInput
} from 'naive-ui'
import { computed, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import type { UserPublicSettingsQuery } from '@/gql/graphql'
import type { OmitDeep } from 'type-fest'

import { useUserSettings } from '@/api'
import { useToast } from '@/components/ui/toast'

type Link = Omit<UserPublicSettingsQuery['userPublicSettings']['socialLinks'][number], '__typename'>

type FormParams = OmitDeep<UserPublicSettingsQuery['userPublicSettings'], '__typename' | 'socialLinks'> & {
	socialLinks: Array<Link & { isEditing?: boolean }>
}

const { t } = useI18n()
const toast = useToast()
const manager = useUserSettings()
const { data } = manager.usePublicQuery()
const updater = manager.usePublicMutation()

const formRef = ref<FormInst | null>(null)
const formData = ref<FormParams>({
	socialLinks: [],
	description: ''
})

watch(data, (v) => {
	if (!v) return
	const rawData = toRaw(v).userPublicSettings
	formData.value = {
		...rawData,
		socialLinks: rawData.socialLinks.map((link) => ({ ...link, isEditing: false }))
	}
}, { immediate: true })

async function save() {
	await updater.executeMutation({
		opts: {
			...formData.value,
			socialLinks: formData.value.socialLinks.map((link) => ({
				title: link.title,
				href: link.href
			}))
		}
	})

	toast.toast({
		title: t('sharedTexts.saved'),
		duration: 1500
	})
}

const rules: FormRules = {
	title: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (value.length === 0) {
				return new Error(t('userSettings.public.errorEmpty'))
			}

			if (value.length > 30) {
				return new Error(t('userSettings.public.errorTooLong'))
			}

			return true
		}
	},
	href: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (value.length === 0) {
				return new Error(t('userSettings.public.errorEmpty'))
			}

			if (value.length > 256) {
				return new Error(t('userSettings.public.errorTooLong'))
			}

			return true
		}
	}
}

const linksLimitReached = computed(() => formData.value.socialLinks.length >= 10)

const newLinkForm = ref({
	title: '',
	href: ''
})

async function addLink() {
	await formRef.value?.validate()
	formData.value.socialLinks.push(newLinkForm.value)
	newLinkForm.value = {
		title: '',
		href: ''
	}
}

function removeLink(index: number) {
	formData.value.socialLinks = formData.value.socialLinks.filter((_, i) => i !== index)
}

function changeSort(from: number, to: number) {
	const element = formData.value.socialLinks.splice(from, 1).at(0)!
	formData.value.socialLinks.splice(to, 0, element)
}
</script>

<template>
	<div class="w-full flex flex-wrap gap-6">
		<div class="flex flex-col w-full gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
				{{ t('userSettings.public.description') }}
			</h4>

			<NCard size="small" bordered>
				<NFormItem :label="t('userSettings.public.description')" :show-feedback="false">
					<NInput
						v-model:value="formData.description"
						show-count
						maxlength="1000"
						type="textarea"
						placeholder=""
						:autosize="{ minRows: 3 }"
					/>
				</NFormItem>
			</NCard>
		</div>

		<div class="flex flex-col w-full gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
				{{ t('userSettings.public.socialLinks') }}
			</h4>

			<NCard size="small" bordered>
				<div class="flex flex-col gap-1">
					<NCard
						v-for="(link, idx) of formData.socialLinks"
						:key="idx"
						size="small"
						embedded
					>
						<template #header>
							<NInput v-if="link.isEditing" v-model:value="link.title" size="small" class="w-[30%]" :maxlength="30" />
							<template v-else>
								{{ link.title }}
							</template>
						</template>
						<template #header-extra>
							<div class="flex gap-2">
								<NButton
									text
									:disabled="!formData.socialLinks[idx + 1]"
									@click="changeSort(idx, idx + 1)"
								>
									<ArrowDownIcon />
								</NButton>
								<NButton
									text
									:disabled="idx === 0"
									@click="changeSort(idx, idx - 1)"
								>
									<ArrowUpIcon />
								</NButton>
								<NButton text @click="link.isEditing = !link.isEditing">
									<EditIcon />
								</NButton>
								<NButton
									text
									@click="removeLink(idx)"
								>
									<TrashIcon />
								</NButton>
							</div>
						</template>

						<NInput
							v-if="link.isEditing"
							v-model:value="link.href"
							size="small"
							type="textarea"
							autosize
							:maxlength="500"
						/>
						<template v-else>
							{{ link.href }}
						</template>
					</NCard>
				</div>
				<NForm ref="formRef" :rules="rules" :model="newLinkForm" class="flex flex-wrap gap-2 items-center w-full mt-5">
					<NFormItem style="--n-label-height: 0px;" :label="t('userSettings.public.linkTitle')" class="flex-auto" path="title">
						<NInput
							v-model:value="newLinkForm.title"
							name="title"
							:maxlength="30"
							placeholder="Twir"
							:disabled=" linksLimitReached"
							:input-props="{ name: 'title' }"
						/>
					</NFormItem>
					<NFormItem style="--n-label-height: 0px;" :label="t('userSettings.public.linkLabel')" class="flex-auto" path="href">
						<NInput
							v-model:value="newLinkForm.href"
							:input-props="{ name: 'href', type: 'url', pattern: 'https?://.+' }"
							:maxlength="256"
							placeholder="https://twir.app"
							:disabled="linksLimitReached"
						/>
					</NFormItem>
					<NButton
						secondary
						attr-type="submit"
						type="success"
						:disabled="linksLimitReached"
						@click="addLink"
					>
						{{ t('sharedButtons.add') }}
					</NButton>
				</NForm>
			</NCard>
		</div>

		<div class="flex justify-start w-full">
			<NButton secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</NButton>
		</div>
	</div>
</template>
