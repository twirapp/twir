<script setup lang="ts">
import { IconPlayerPlay, IconTrash } from '@tabler/icons-vue';
import {
	type FormInst,
	type FormItemRule,
	type FormRules,
	NButton,
	NForm,
	NFormItem,
	NInput,
	NModal,
	NSlider,
	NSpace,
} from 'naive-ui';
import { computed, onMounted, ref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { type EditableAlert } from './types.js';

import { useAlertsManager, useFiles, useProfile } from '@/api';
import FilesPicker from '@/components/files/files.vue';
import { playAudio } from '@/helpers/index.js';

const props = defineProps<{
	alert?: EditableAlert | null
}>();
const emits = defineEmits<{
	close: []
}>();

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableAlert>({
	id: '',
	name: '',
	audioId: undefined,
	audioVolume: 100,
});

onMounted(() => {
	if (!props.alert) return;
	formValue.value = structuredClone(toRaw(props.alert));
});

const { t } = useI18n();

const rules: FormRules = {
	name: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length || value.length > 30) {
				return new Error(t('alerts.validations.name'));
			}

			return true;
		},
	},
};

const manager = useAlertsManager();
const creator = manager.create;
const updater = manager.update;

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;

	if (data.id) {
		await updater.mutateAsync({
			...data,
			id: data.id!,
		});
	} else {
		await creator.mutateAsync(data);
	}

	emits('close');
}

const { data: files } = useFiles();
const selectedAudio = computed(() => files.value?.files.find(f => f.id === formValue.value.audioId));
const showAudioModal = ref(false);

const { data: profile } = useProfile();

async function testAudio() {
	if (!selectedAudio.value?.id || !profile.value) return;

	const query = new URLSearchParams({
		channel_id: profile.value.selectedDashboardId,
		file_id: selectedAudio.value.id,
	});

	const req = await fetch(`${window.location.origin}/api/files/?${query}`);
	if (!req.ok) {
		console.error(await req.text());
		return;
	}

	await playAudio(await req.arrayBuffer(), formValue.value.audioVolume);
}
</script>

<template>
	<n-form
		ref="formRef"
		:model="formValue"
		:rules="rules"
	>
		<n-space vertical style="width: 100%">
			<n-form-item label="Name" path="name" show-require-mark>
				<n-input v-model:value="formValue.name" :maxlength="30" />
			</n-form-item>

			<n-form-item label="Audio">
				<div style="display: flex; gap: 10px; width: 85%">
					<n-button block type="info" @click="showAudioModal = true">
						{{ selectedAudio?.name ?? t('sharedButtons.select') }}
					</n-button>
					<n-button
						:disabled="!formValue.audioId" text type="error"
						@click="formValue.audioId = undefined"
					>
						<IconTrash />
					</n-button>
					<n-button :disabled="!formValue.audioId" text type="info" @click="testAudio">
						<IconPlayerPlay />
					</n-button>
				</div>
			</n-form-item>

			<n-form-item label="Audio Volume">
				<n-slider
					v-model:value="formValue.audioVolume"
					show-tooltip
					:step="1"
					:min="1"
					:max="100"
					:marks="{ 1: '1', 50: '50', 100: '100' }"
				/>
			</n-form-item>

			<n-form-item label="Image">
				<n-button block type="info" disabled>
					Soon...
				</n-button>
			</n-form-item>

			<n-form-item label="Text">
				<n-button block type="info" disabled>
					Soon...
				</n-button>
			</n-form-item>
		</n-space>

		<n-button secondary type="success" block @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-form>

	<n-modal
		v-model:show="showAudioModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Select audio"
		class="modal"
		:style="{
			width: '1000px',
			top: '50px',
		}"
		:on-close="() => showAudioModal = false"
	>
		<files-picker
			mode="picker"
			tab="audios"
			@select="(id) => {
				formValue.audioId = id
				showAudioModal = false
			}"
			@delete="(id) => {
				if (id === formValue.audioId) {
					formValue.audioId = undefined
				}
			}"
		/>
	</n-modal>
</template>
