<script setup lang="ts">
import { IconEdit, IconTrash, IconArrowUp, IconArrowDown } from '@tabler/icons-vue';
import {
	NCard,
	NInput,
	NFormItem,
	NButton,
	NForm,
	type FormItemRule,
} from 'naive-ui';
import type { FormRules, FormInst } from 'naive-ui';
import { OmitDeep } from 'type-fest';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserSettings } from '@/api';
import { useToast } from '@/components/ui/toast';
import { UserPublicSettingsQuery } from '@/gql/graphql';

type FormParams = OmitDeep<UserPublicSettingsQuery['userPublicSettings'], '__typename' | 'socialLinks'> & {
	socialLinks: Array<UserPublicSettingsQuery['userPublicSettings']['socialLinks'][number] & { isEditing?: boolean }>
}

const { t } = useI18n();
const toast = useToast();
const manager = useUserSettings();
const { data } = manager.usePublicQuery();
const updater = manager.usePublicMutation();

const formRef = ref<FormInst | null>(null);
const formData = ref<FormParams>({
	socialLinks: [],
	description: '',
});

watch(data, (v) => {
	if (!v) return;

	console.log(v);

	const rawData = toRaw(v).userPublicSettings;

	formData.value = {
		...rawData,
		socialLinks: rawData.socialLinks.map((link) => ({ ...link, isEditing: false })),
	};
}, { immediate: true });

async function save() {
	await updater.executeMutation({
		opts: formData.value,
	});

	toast.toast({
		title: t('sharedTexts.saved'),
		duration: 1500,
	});
}

const rules: FormRules = {
	title: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (value.length === 0) {
				return new Error(t('userSettings.public.errorEmpty'));
			}

			if (value.length > 30) {
				return new Error(t('userSettings.public.errorTooLong'));
			}

			return true;
		},
	},
	href: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (value.length === 0) {
				return new Error(t('userSettings.public.errorEmpty'));
			}

			if (value.length > 256) {
				return new Error(t('userSettings.public.errorTooLong'));
			}

			return true;
		},
	},
};

const linksLimitReached = computed(() => formData.value.socialLinks.length >= 10);

const newLinkForm = ref({
	title: '',
	href: '',
});

async function addLink() {
	await formRef.value?.validate();
	formData.value.socialLinks.push(newLinkForm.value);
	newLinkForm.value = {
		title: '',
		href: '',
	};
}

function removeLink(index: number) {
	formData.value.socialLinks = formData.value.socialLinks.filter((_, i) => i != index);
}

function changeSort(from: number, to: number) {
	const element = formData.value.socialLinks.splice(from, 1).at(0)!;
	formData.value.socialLinks.splice(to, 0, element);
}
</script>

<template>
	<div class="w-full flex flex-wrap gap-6">
		<div class="flex flex-col w-full gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
				{{ t('userSettings.public.description') }}
			</h4>

			<n-card size="small" bordered>
				<n-form-item :label="t('userSettings.public.description')" :show-feedback="false">
					<n-input
						v-model:value="formData.description"
						show-count
						maxlength="1000"
						type="textarea"
						placeholder=""
						:autosize="{ minRows: 3 }"
					/>
				</n-form-item>
			</n-card>
		</div>


		<div class="flex flex-col w-full gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
				{{ t('userSettings.public.socialLinks') }}
			</h4>

			<n-card size="small" bordered>
				<div class="flex flex-col gap-1">
					<n-card
						v-for="(link, idx) of formData.socialLinks"
						:key="idx"
						size="small"
						embedded
					>
						<template #header>
							<n-input v-if="link.isEditing" v-model:value="link.title" size="small" class="w-[30%]" :maxlength="30" />
							<template v-else>
								{{ link.title }}
							</template>
						</template>
						<template #header-extra>
							<div class="flex gap-2">
								<n-button
									text
									:disabled="!formData.socialLinks[idx+1]"
									@click="changeSort(idx, idx+1)"
								>
									<IconArrowDown />
								</n-button>
								<n-button
									text
									:disabled="idx === 0"
									@click="changeSort(idx, idx-1)"
								>
									<IconArrowUp />
								</n-button>
								<n-button text @click="link.isEditing = !link.isEditing">
									<IconEdit />
								</n-button>
								<n-button
									text
									@click="removeLink(idx)"
								>
									<IconTrash />
								</n-button>
							</div>
						</template>

						<n-input
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
					</n-card>
				</div>
				<n-form ref="formRef" :rules="rules" :model="newLinkForm" class="flex flex-wrap gap-2 items-center w-full mt-5">
					<n-form-item style="--n-label-height: 0px;" :label="t('userSettings.public.linkTitle')" class="flex-auto" path="title">
						<n-input
							v-model:value="newLinkForm.title"
							name="title"
							:maxlength="30"
							placeholder="Twir"
							:disabled=" linksLimitReached"
							:input-props="{ name: 'title' }"
						/>
					</n-form-item>
					<n-form-item style="--n-label-height: 0px;" :label="t('userSettings.public.linkLabel')" class="flex-auto" path="href">
						<n-input
							v-model:value="newLinkForm.href"
							:input-props="{ name: 'href', type: 'url', pattern: 'https?://.+' }"
							:maxlength="256"
							placeholder="https://twir.app"
							:disabled="linksLimitReached"
						/>
					</n-form-item>
					<n-button
						secondary
						attr-type="submit"
						type="success"
						:disabled="linksLimitReached"
						@click="addLink"
					>
						{{ t('sharedButtons.add') }}
					</n-button>
				</n-form>
			</n-card>
		</div>

		<div class="flex justify-start w-full">
			<n-button secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</div>
	</div>
</template>
