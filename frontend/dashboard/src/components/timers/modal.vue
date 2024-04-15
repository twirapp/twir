<script setup lang='ts'>
import {
	IconTrash,
	IconPlus,
	IconArrowNarrowUp,
	IconArrowNarrowDown,
} from '@tabler/icons-vue';
import {
	type FormInst,
	type FormRules,
	type FormItemRule,
	NForm,
	NFormItem,
	NInput,
	NInputNumber,
	NButton,
	NSlider,
	NGrid,
	NGridItem,
	NDynamicInput,
	NDivider,
	NSpace,
	NTimeline,
	NTimelineItem,
	NCheckbox,
	useThemeVars,
} from 'naive-ui';
import { ref, onMounted, toRaw, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useTimersApi } from '@/api/timers.js';
import type { EditableTimer, EditableTimerResponse } from '@/components/timers/types.js';

const props = defineProps<{
	timer?: EditableTimer | null
}>();
const emits = defineEmits<{
	close: []
}>();

const nowTime = new Date();

const { t } = useI18n();

const themeVars = useThemeVars();
const configBackground = computed(() => themeVars.value.actionColor);

const formRef = ref<FormInst | null>(null);
const formValue = ref<EditableTimer>({
	name: '',
	enabled: true,
	messageInterval: 0,
	timeInterval: 5,
	responses: [],
});
const rules: FormRules = {
	name: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('timers.modal.validations.nameRequired'));
			}
			if (value.length > 25) {
				return new Error(t('timers.modal.validations.nameLong'));
			}
			return true;
		},
	},
	timeInterval: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: number) => {
			if (value < 1) {
				return new Error('Time interval is too short');
			}
			if (value > 100) {
				return new Error('Time interval is too long');
			}
			return true;
		},
	},
	responses: {
		trigger: ['input', 'blur'],
		validator: (_: FormItemRule, value: string) => {
			if (!value || !value.length) {
				return new Error(t('timers.modal.validations.responseRequired'));
			}
			if (value.length > 500) {
				return new Error(t('timers.modal.validations.responseLong'));
			}
			return true;
		},
	},
};

onMounted(() => {
	if (!props.timer) return;
	formValue.value = structuredClone(toRaw(props.timer));
});

const timersApi = useTimersApi();
const timersUpdate = timersApi.useMutationUpdateTimer();
const timersCreate = timersApi.useMutationCreateTimer();

async function save() {
	if (!formRef.value || !formValue.value) return;
	await formRef.value.validate();

	const data = formValue.value;

	if (data.id) {
		await timersUpdate.executeMutation({
			id: data.id,
			opts: {
				name: data.name,
				enabled: data.enabled,
				messageInterval: data.messageInterval,
				timeInterval: data.timeInterval,
				responses: data.responses,
			},
		});
	} else {
		await timersCreate.executeMutation({ opts: data });
	}

	emits('close');
}

const sliderMarks = {
	5: '5',
	15: '15',
	30: '30',
	45: '45',
	60: '60',
	75: '75',
	90: '90',
	100: '100',
};
</script>

<template>
	<n-form
		ref="formRef"
		:model="formValue"
		:rules="rules"
	>
		<n-space vertical class="w-full">
			<n-form-item :label="t('sharedTexts.name')" path="name" show-require-mark>
				<n-input
					v-model:value="formValue.name" :placeholder="t('sharedTexts.name')"
					:maxlength="25"
				/>
			</n-form-item>

			<div
				:style="{ 'background-color': configBackground }"
				style="border-radius: 11px; padding: 8px"
			>
				<n-form-item
					:label="t('timers.table.columns.intervalInMinutes')" path="timeInterval"
					show-require-mark
				>
					<n-grid :cols="12" :x-gap="5">
						<n-grid-item :span="10">
							<n-slider
								v-model:value="formValue.timeInterval"
								:step="1"
								:marks="sliderMarks"
								:min="1"
							/>
						</n-grid-item>
						<n-grid-item :span="2">
							<n-input-number
								v-model:value="formValue.timeInterval" size="small" :min="1"
								:max="100"
							/>
						</n-grid-item>
					</n-grid>
				</n-form-item>

				<n-divider dashed class="p-0 m-0">
					AND
				</n-divider>

				<n-form-item :label="t('timers.table.columns.intervalInMessages')">
					<n-input-number v-model:value="formValue.messageInterval" :min="0" :max="5000" />
				</n-form-item>
			</div>


			<!--      <n-form-item label="Interval in messages" path="messageInterval">-->
			<!--        <n-input-number v-model:value="formValue.messageInterval" :min="0" :max="5000" />-->
			<!--      </n-form-item>-->

			<n-divider>
				{{ t('sharedTexts.responses') }}
			</n-divider>

			<n-dynamic-input
				v-model:value="formValue.responses"
				class="groups"
				:create-button-props="({ class: 'create-button' } as any)"
			>
				<template #default="{ value, index }: { value: EditableTimerResponse, index: number }">
					<n-space vertical class="w-full">
						<n-form-item
							:path="`responses[${index}].text`" :rule="rules.responses"
							show-require-mark
						>
							<n-input
								v-model:value="value.text"
								type="textarea"
								placeholder="text"
								:autosize="{
									minRows: 1,
									maxRows: 5
								}"
							/>
						</n-form-item>
						<n-checkbox v-model:checked="value.isAnnounce">
							Use twitch announce
						</n-checkbox>
					</n-space>
				</template>

				<template #action="{ index, remove, move }">
					<div class="flex gap-x-[5px] items-center ml-2.5">
						<n-button size="small" type="error" quaternary @click="() => remove(index)">
							<IconTrash />
						</n-button>
						<n-button
							size="small"
							type="info"
							quaternary
							:disabled="index == 0"
							@click="() => move('up', index)"
						>
							<IconArrowNarrowUp />
						</n-button>
						<n-button
							size="small"
							type="info"
							quaternary
							:disabled="!!formValue.responses.length && index === formValue.responses.length-1"
							@click="() => move('down', index)"
						>
							<IconArrowNarrowDown />
						</n-button>
					</div>
				</template>
			</n-dynamic-input>

			<n-button
				dashed
				block
				@click="() => formValue.responses.push({ text: '', isAnnounce: false })"
			>
				<IconPlus />
				{{ t('sharedButtons.create') }}
			</n-button>

			<n-divider />

			{{ t('timers.modal.timelineDescription') }}
			<n-timeline>
				<n-timeline-item
					v-for="(response, index) in formValue.responses"
					:key="index"
					type="info"
					:title="`Response #${index+1}`"
					:content="response.text.slice(0, 50) + '...'"
					:time="new Date(nowTime.getTime() + (formValue.timeInterval * index * 60000)).toLocaleTimeString()"
					line-type="dashed"
				/>
			</n-timeline>

			<n-divider />

			<n-button secondary type="success" block class="mt-2.5" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</n-space>
	</n-form>
</template>

<style scoped>
.groups :deep(.create-button) {
	display: none;
}
</style>
