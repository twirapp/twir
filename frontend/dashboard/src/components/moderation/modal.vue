<script setup lang="ts">
import { IconSquare, IconSquareCheck } from '@tabler/icons-vue';
import chunk from 'lodash.chunk';
import {
	NButton,
	NFormItem,
	NInput,
	NInputNumber,
	NDivider,
	NButtonGroup,
	useNotification,
} from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useEditableItem } from './helpers.js';
import ModalCaps from './modal-caps.vue';
import ModalDenylist from './modal-denylist.vue';
import ModalEmotes from './modal-emotes.vue';
import ModalLanguage from './modal-language.vue';
import ModalLongMessage from './modal-longmessage.vue';
import ModalSymbols from './modal-symbols.vue';

import { useModerationManager, useRolesManager } from '@/api';

const { t } = useI18n();

const manager = useModerationManager();
const updater = manager.update;
const creator = manager.create;

const { editableItem } = useEditableItem();

const { data: availableRoles } = useRolesManager().getAll({});
const rolesSelectOptions = computed(() => {
	if (!availableRoles.value) return [];
	return availableRoles.value.roles
		.filter(r => !['BROADCASTER', 'MODERATOR'].includes(r.type))
		.map((role) => ({
			label: role.name,
			value: role.id,
		}));
});

const message = useNotification();

async function saveSettings() {
	if (!editableItem.value) return;

	if (!editableItem.value.id) {
		await creator.mutateAsync({
			data: editableItem.value.data,
		});
	} else {
		await updater.mutateAsync({
			id: editableItem.value.id,
			data: editableItem.value.data,
		});
	}
	message.success({
		title: t('sharedTexts.saved'),
		duration: 2000,
	});
}
</script>

<template>
	<div style="display: flex; flex-direction: column; gap: 12px;">
		<modal-symbols
			v-if="editableItem?.data?.type === 'symbols'"
			class="form-block"
		/>

		<modal-language
			v-if="editableItem?.data?.type === 'language'"
			class="form-block"
		/>

		<modal-long-message
			v-if="editableItem?.data?.type === 'long_message'"
			class="form-block"
		/>

		<modal-caps
			v-if="editableItem?.data?.type === 'caps'"
			class="form-block"
		/>

		<modal-emotes
			v-if="editableItem?.data?.type === 'emotes'"
			class="form-block"
		/>

		<div class="form-block">
			<n-form-item v-if="editableItem?.data" label="Timeout message">
				<n-input
					v-model:value="editableItem.data.banMessage"
					type="textarea"
					:maxLength="500"
					autosize
				/>
			</n-form-item>

			<n-form-item v-if="editableItem?.data" :label="t('moderation.banTime')" :feedback="t('moderation.banDescription')">
				<n-input-number
					v-model:value="editableItem.data.banTime"
					:min="0"
					:max="86400"
				/>
			</n-form-item>
		</div>

		<n-divider style="margin: 0; padding: 0" />

		<div class="form-block">
			<n-form-item v-if="editableItem?.data" :label="t('moderation.warningMessage')">
				<n-input
					v-model:value="editableItem.data.warningMessage"
					type="textarea"
					:maxLength="500"
					autosize
				/>
			</n-form-item>

			<n-form-item v-if="editableItem?.data" :label="t('moderation.warningMaxCount')">
				<n-input-number
					v-model:value="editableItem.data.maxWarnings"
					:min="0"
					:max="10"
				/>
			</n-form-item>
		</div>

		<n-divider style="margin: 0; padding: 0" />

		<div class="form-block">
			<span>{{ t('moderation.excludedRoles') }}</span>
			<div v-if="editableItem?.data" style="display: flex; flex-direction: column; gap: 5px;">
				<n-button-group
					v-for="(group, index) of chunk(rolesSelectOptions.sort(), 5)"
					:key="index"
				>
					<n-button
						v-for="option of group"
						:key="option.value"
						:type="editableItem?.data?.excludedRoles.includes(option.value) ? 'success' : 'default'"
						secondary
						@click="() => {
							if (editableItem!.data!.excludedRoles.includes(option.value)) {
								editableItem!.data!.excludedRoles = editableItem!.data!.excludedRoles.filter(r => r !== option.value)
							} else {
								editableItem!.data!.excludedRoles.push(option.value)
							}
						}"
					>
						<template #icon>
							<IconSquareCheck v-if="editableItem?.data?.excludedRoles.includes(option.value)" />
							<IconSquare v-else />
						</template>
						{{ option.label }}
					</n-button>
				</n-button-group>
			</div>
		</div>

		<modal-denylist
			v-if="editableItem?.data?.type === 'deny_list'"
			class="form-block"
		/>

		<n-divider style="margin: 0; padding: 0" />

		<n-button type="success" secondary @click="saveSettings">
			{{ t('sharedButtons.save') }}
		</n-button>
	</div>
</template>

<style scoped>
.form-block {
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.form-block .content {
	padding: 8px;
}

.form-block :deep(.n-input-number) {
	width: 100%;
}
</style>
