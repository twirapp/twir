<script setup lang="ts">
import { IconEdit, IconTrash, IconArrowUp, IconArrowDown } from '@tabler/icons-vue';
import type {
	Settings,
	SocialLink,
} from '@twir/api/messages/channels_public_settings/channels_public_settings';
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

type SocialLinkWithEdit = SocialLink & { isEditing?: boolean };
type FormData = Omit<Settings, 'socialLinks'> & { socialLinks: SocialLinkWithEdit[] }

const formRef = ref<FormInst | null>(null);
const formData = ref<FormData>({
	socialLinks: [],
	description: undefined,
});

watch(data, (v) => {
	if (!v) return;

	const rawData = toRaw(v);

	formData.value = {
		...rawData,
		socialLinks: rawData.socialLinks.map((link) => ({ ...link, isEditing: false })),
	};
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

const newLinkForm = ref({
	title: '',
	href: '',
});

function addLink() {
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
	<div class="w-full flex flex-wrap gap-4">
		<n-card title="Description" size="small" bordered>
			<n-form-item :label="t('userSettings.public.description')">
				<n-input
					v-model:value="formData.description"
					type="textarea"
					placeholder=""
					:autosize="{ minRows: 3 }"
				/>
			</n-form-item>
		</n-card>

		<n-card title="Social links" size="small" bordered>
			<div class="flex flex-col gap-1">
				<n-card
					v-for="(link, idx) of formData.socialLinks"
					:key="idx"
					size="small"
					embedded
				>
					<template #header>
						<n-input v-if="link.isEditing" v-model:value="link.title" size="small" style="width: 30%" :maxlength="30" />
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
								<IconArrowDown class="header-button" />
							</n-button>
							<n-button
								text
								:disabled="idx === 0"
								@click="changeSort(idx, idx-1)"
							>
								<IconArrowUp class="header-button" />
							</n-button>
							<n-button text @click="link.isEditing = !link.isEditing">
								<IconEdit class="header-button" />
							</n-button>
							<n-button
								text
								@click="removeLink(idx)"
							>
								<IconTrash class="header-button" />
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
				<n-form-item style="--n-label-height: 0px;" label="Title" class="flex-auto" path="title">
					<n-input
						v-model:value="newLinkForm.title"
						:maxlength="30"
						placeholder="Twir"
						:disabled=" linksLimitReached"
					/>
				</n-form-item>
				<n-form-item style="--n-label-height: 0px;" label="Href" class="flex-auto" path="href">
					<n-input
						v-model:value="newLinkForm.href"
						:maxlength="500"
						placeholder="https://twir.app"
						:disabled="linksLimitReached"
					/>
				</n-form-item>
				<n-button
					secondary
					type="success"
					:disabled="linksLimitReached"
					@click="addLink"
				>
					{{ t('sharedButtons.add') }}
				</n-button>
			</n-form>
		</n-card>

		<div class="flex justify-start w-full">
			<n-button secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</div>
	</div>
</template>

<style scoped>
.header-button {
	height: 18px;
	width: 18px;
}
</style>
