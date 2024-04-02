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
  NDivider,
  NSelect,
} from 'naive-ui';
import { computed, onMounted, ref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { type EditableAlert } from './types.js';
import rewardsSelector from '../rewardsSelector.vue';

import {
  useAlertsManager,
  useCommandsManager,
  useFiles,
  useGreetingsManager, useKeywordsManager,
  useProfile, useTwitchGetUsers,
} from '@/api';
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
  commandIds: [],
  rewardIds: [],
  greetingsIds: [],
  keywordsIds: [],
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

const commandsManager = useCommandsManager();
const { data: commands } = commandsManager.getAll({});
const commandsSelectOptions = computed(() => commands.value?.commands
    .map(c => ({ label: c.name, value: c.id })),
);

const greetingsManager = useGreetingsManager();
const { data: greetings } = greetingsManager.getAll({});
const greetingsUsersIds = computed(() => greetings.value?.greetings.map(g => g.userId) ?? []);
const { data: twitchUsers } = useTwitchGetUsers({ ids: greetingsUsersIds });
const greetingsSelectOptions = computed(() => {
  if (!greetingsUsersIds.value.length || !twitchUsers.value?.users.length) return [];
  return greetings.value?.greetings.map(g => {
    const twitchUser = twitchUsers.value.users.find(u => u.id === g.userId);
    return { label: twitchUser?.login ?? g.userId, value: g.id };
  });
});

const keywordsManager = useKeywordsManager();
const { data: keywords } = keywordsManager.getAll({});
const keywordsSelectOptions = computed(() => keywords.value?.keywords
    .map(k => ({ label: k.text, value: k.id })),
);
</script>

<template>
	<n-form
		ref="formRef"
		:model="formValue"
		:rules="rules"
	>
		<n-space vertical class="w-full">
			<n-form-item label="Name" path="name" show-require-mark>
				<n-input v-model:value="formValue.name" :maxlength="30" />
			</n-form-item>

			<n-divider />

			<n-form-item :label="t('alerts.trigger.commands')" path="commandIds">
				<n-select
					v-model:value="formValue.commandIds"
					:fallback-option="false"
					filterable
					multiple
					:options="commandsSelectOptions"
				/>
			</n-form-item>

			<n-form-item :label="t('alerts.trigger.rewards')" path="rewardIds">
				<rewardsSelector v-model="formValue.rewardIds" multiple />
			</n-form-item>

			<n-form-item :label="t('alerts.trigger.keywords')" path="rewardIds">
				<n-select
					v-model:value="formValue.keywordsIds"
					:fallback-option="false"
					filterable
					multiple
					:options="keywordsSelectOptions"
				/>
			</n-form-item>

			<n-form-item :label="t('alerts.trigger.greetings')" path="rewardIds">
				<n-select
					v-model:value="formValue.greetingsIds"
					:fallback-option="false"
					filterable
					multiple
					:options="greetingsSelectOptions"
				/>
			</n-form-item>

			<n-divider />

			<n-form-item :label="t('alerts.select.audio')">
				<div class="flex gap-2.5 w-[85%]">
					<n-button class="overflow-hidden text-nowrap" block type="info" @click="showAudioModal = true">
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

			<n-form-item :label="t('alerts.audioVolume', { volume: formValue.audioVolume })">
				<n-slider
					v-model:value="formValue.audioVolume"
					:step="1"
					:min="1"
					:max="100"
					:marks="{ 1: '1', 100: '100' }"
					:show-tooltip="false"
					:tooltip="false"
				/>
			</n-form-item>

			<n-form-item :label="t('alerts.select.image')">
				<n-button block type="info" disabled>
					Soon...
				</n-button>
			</n-form-item>

			<n-form-item :label="t('alerts.select.text')">
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
		:title="t('alerts.select.audio')"
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
