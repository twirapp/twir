<script setup lang="ts">
import { IconEdit, IconTrash } from '@tabler/icons-vue';
import type { Settings } from '@twir/api/messages/channels_public_settings/channels_public_settings';
import {
	NCard,
	NInput,
	NFormItem,
	NButton,
	NForm,
	type FormItemRule,
} from 'naive-ui';
import type { FormRules, FormInst } from 'naive-ui';
import { computed, ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { usePublicSettings } from '@/api/public-settings';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';
import { linkRegex } from '@/pages/user-settings/link-regex';

const { t } = useI18n();
const { notification } = useNaiveDiscrete();
const manager = usePublicSettings();
const { data } = manager.useGet();
const updater = manager.useUpdate();

const formRef = ref<FormInst | null>(null);
const formData = ref<Settings>({
	socialLinks: [],
	description: undefined,
});

watch(data, (v) => {
	if (!v) return;

	formData.value = toRaw(v);
}, { immediate: true });

async function save() {
	await updater.mutateAsync(formData.value);

	notification.create({
		title: t('sharedTexts.saved'),
		type: 'success',
		duration: 2500,
	});
}

const rules: FormRules = {
	href: {
		trigger: ['input', 'blur', 'focus'],
		validator: (_: FormItemRule, value: string) => {
			if (value.length === 0) {
				return new Error('Cannot be empty');
			}
			if (value.length > 500) {
				return new Error('Invalid length');
			}
			const isLink = linkRegex.test(value);
			if (!isLink) {
				return new Error('Invalid link');
			}

			return true;
		},
	},
};

const linksLimitReached = computed(() => formData.value.socialLinks.length >= 10);

const linkForm = ref<{
	title: string,
	href: string,
	existedIndex?: number
}>({
	title: '',
	href: '',
});
const isLinkEdit = computed(() => typeof linkForm.value.existedIndex != 'undefined');

function clearLinkForm() {
	linkForm.value = {
		title: '',
		href: '',
	};
}

function saveLink() {
	if (typeof linkForm.value.existedIndex !== 'undefined') {
		const link = formData.value.socialLinks.at(linkForm.value.existedIndex)!;
		link.href = linkForm.value.href;
		link.title = linkForm.value.title;
	} else {
		formData.value.socialLinks.push(linkForm.value);
	}
	clearLinkForm();
}

function setLinkFormEdit(index: number) {
	const existedLink = formData.value.socialLinks.at(index)!;
	linkForm.value = {
		href: existedLink.href,
		title: existedLink.title,
		existedIndex: index,
	};
}

function removeLink(index: number) {
	formData.value.socialLinks = formData.value.socialLinks.filter((_, i) => i != index);
}
</script>

<template>
	<div class="w-full flex flex-wrap gap-4">
		<n-card title="Description" size="small" bordered>
			<n-form-item label="Text in your profile">
				<n-input
					v-model:value="formData.description"
					type="textarea"
					placeholder=""
					:autosize="{ minRows: 3 }"
				/>
			</n-form-item>
		</n-card>

		<n-card title="Social links" size="small" bordered>
			<n-form ref="formRef" :rules="rules" :model="linkForm" class="flex flex-wrap gap-2 items-center w-full">
				<n-form-item style="--n-label-height: 0px;" label="Title" class="flex-auto" path="title">
					<n-input
						v-model:value="linkForm.title"
						:maxlength="30"
						placeholder="Enter title"
						:disabled="!isLinkEdit && linksLimitReached"
					/>
				</n-form-item>
				<n-form-item style="--n-label-height: 0px;" label="Href" class="flex-auto" path="href">
					<n-input
						v-model:value="linkForm.href"
						:maxlength="500"
						placeholder="Enter link"
						:disabled="!isLinkEdit && linksLimitReached"
					/>
				</n-form-item>
				<n-button
					secondary
					type="success"
					:disabled="!isLinkEdit && linksLimitReached"
					@click="saveLink"
				>
					{{ isLinkEdit ? 'Save' : 'Add' }}
				</n-button>
				<n-button
					v-if="isLinkEdit"
					secondary
					type="warning"
					@click="clearLinkForm"
				>
					Cancel
				</n-button>
			</n-form>

			<div class="flex flex-col gap-1">
				<n-card
					v-for="(link, idx) of formData.socialLinks"
					:key="idx"
					size="small"
					embedded
				>
					<template #header>
						{{ link.title }}
					</template>
					<template #header-extra>
						<div class="flex gap-2">
							<n-button text @click="setLinkFormEdit(idx)">
								<IconEdit style="height: 18px; width: 18px;" />
							</n-button>
							<n-button
								text
								@click="removeLink(idx)"
							>
								<IconTrash style="height: 18px; width: 18px;" />
							</n-button>
						</div>
					</template>
					{{ link.href }}
				</n-card>
			</div>
		</n-card>

		<div class="flex justify-start w-full">
			<n-button secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</div>
	</div>
</template>
