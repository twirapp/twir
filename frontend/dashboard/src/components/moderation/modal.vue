<script setup lang="ts">
import { IconSquare, IconSquareCheck } from '@tabler/icons-vue';
import { type ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import chunk from 'lodash.chunk';
import { NButton, NFormItem, NInput, NInputNumber, NDivider, NButtonGroup, useNotification } from 'naive-ui';
import { computed, ref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

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

const props = defineProps<{
	item: ItemWithId
}>();
const formValue = ref(structuredClone(toRaw(props.item)));

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
	await updater.mutateAsync(formValue.value);
	message.success({
		title: t('sharedTexts.saved'),
	});
}
</script>

<template>
	<div style="display: flex; flex-direction: column; gap: 12px;">
		<modal-symbols
			v-if="formValue.data!.type === 'symbols'"
			class="form-block"
			:item="formValue"
		/>

		<modal-language
			v-if="formValue.data!.type === 'language'"
			class="form-block"
			:item="formValue"
		/>

		<modal-long-message
			v-if="formValue.data!.type === 'long_message'"
			class="form-block"
			:item="formValue"
		/>

		<modal-caps
			v-if="formValue.data!.type === 'caps'"
			class="form-block"
			:item="formValue"
		/>

		<modal-emotes
			v-if="formValue.data!.type === 'emotes'"
			class="form-block"
			:item="formValue"
		/>

		<div class="form-block">
			<n-form-item label="Timeout message" feedback="qwe">
				<n-input
					v-model:value="formValue.data!.banMessage"
					type="textarea"
					:maxLength="500"
				/>
			</n-form-item>

			<n-form-item label="Ban time" feedback="qwe">
				<n-input-number
					v-model:value="formValue.data!.banTime"
					:min="0"
					:max="86400"
				/>
			</n-form-item>
		</div>

		<n-divider style="margin: 0; padding: 0" />

		<div class="form-block">
			<n-form-item label="Warning message">
				<n-input
					v-model:value="formValue.data!.warningMessage"
					type="textarea"
					:maxLength="500"
				/>
			</n-form-item>

			<n-form-item label="Warnins count">
				<n-input-number
					v-model:value="formValue.data!.maxWarnings"
					:min="0"
					:max="10"
				/>
			</n-form-item>
		</div>

		<n-divider style="margin: 0; padding: 0" />

		<div class="form-block">
			<span>Excluded for moderation roles</span>
			<div style="display: flex; flex-direction: column; gap: 5px;">
				<n-button-group
					v-for="(group, index) of chunk(rolesSelectOptions.sort(), 5)"
					:key="index"
				>
					<n-button
						v-for="option of group"
						:key="option.value"
						:type="formValue.data!.excludedRoles.includes(option.value) ? 'success' : 'default'"
						secondary
						@click="() => {
							if (formValue.data!.excludedRoles.includes(option.value)) {
								formValue.data!.excludedRoles = formValue.data!.excludedRoles.filter(r => r !== option.value)
							} else {
								formValue.data!.excludedRoles.push(option.value)
							}
						}"
					>
						<template #icon>
							<IconSquareCheck v-if="formValue.data!.excludedRoles.includes(option.value)" />
							<IconSquare v-else />
						</template>
						{{ option.label }}
					</n-button>
				</n-button-group>
			</div>
		</div>

		<modal-denylist
			v-if="formValue.data!.type === 'deny_list'"
			class="form-block"
			:item="formValue"
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
