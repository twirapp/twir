<script setup lang="ts">
import { IconSquare, IconSquareCheck } from '@tabler/icons-vue'
import chunk from 'lodash.chunk'
import {
	NButton,
	NButtonGroup,
	NDivider,
	NFormItem,
	NInput,
	NInputNumber,
	useNotification
} from 'naive-ui'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { useEditableItem } from './helpers.js'
import ModalCaps from './modal-caps.vue'
import ModalDenylist from './modal-denylist.vue'
import ModalEmotes from './modal-emotes.vue'
import ModalLanguage from './modal-language.vue'
import ModalLongMessage from './modal-longmessage.vue'
import ModalSymbols from './modal-symbols.vue'

import { useModerationManager } from '@/api'
import { useRoles } from '@/api/roles'

const { t } = useI18n()

const manager = useModerationManager()
const updater = manager.update
const creator = manager.create

const { editableItem } = useEditableItem()

const rolesManager = useRoles()
const { data: availableRoles } = rolesManager.useRolesQuery()

const rolesSelectOptions = computed(() => {
	if (!availableRoles.value?.roles) return []
	return availableRoles.value.roles
		.filter(r => !['BROADCASTER', 'MODERATOR'].includes(r.type))
		.map((role) => ({
			label: role.name,
			value: role.id
		}))
})

const message = useNotification()

async function saveSettings() {
	if (!editableItem.value) return

	if (!editableItem.value.id) {
		await creator.mutateAsync({
			data: editableItem.value.data
		})
	} else {
		await updater.mutateAsync({
			id: editableItem.value.id,
			data: editableItem.value.data
		})
	}
	message.success({
		title: t('sharedTexts.saved'),
		duration: 2000
	})
}
</script>

<template>
	<div class="flex flex-col gap-3">
		<ModalSymbols
			v-if="editableItem?.data?.type === 'symbols'"
			class="form-block"
		/>

		<ModalLanguage
			v-if="editableItem?.data?.type === 'language'"
			class="form-block"
		/>

		<ModalLongMessage
			v-if="editableItem?.data?.type === 'long_message'"
			class="form-block"
		/>

		<ModalCaps
			v-if="editableItem?.data?.type === 'caps'"
			class="form-block"
		/>

		<ModalEmotes
			v-if="editableItem?.data?.type === 'emotes'"
			class="form-block"
		/>

		<div class="form-block">
			<NFormItem v-if="editableItem?.data" label="Timeout message">
				<NInput
					v-model:value="editableItem.data.banMessage"
					type="textarea"
					:maxLength="500"
					autosize
				/>
			</NFormItem>

			<NFormItem v-if="editableItem?.data" :label="t('moderation.banTime')" :feedback="t('moderation.banDescription')">
				<NInputNumber
					v-model:value="editableItem.data.banTime"
					:min="0"
					:max="86400"
				/>
			</NFormItem>
		</div>

		<NDivider class="m-0 p-0" />

		<div class="form-block">
			<NFormItem v-if="editableItem?.data" :label="t('moderation.warningMessage')">
				<NInput
					v-model:value="editableItem.data.warningMessage"
					type="textarea"
					:maxLength="500"
					autosize
				/>
			</NFormItem>

			<NFormItem v-if="editableItem?.data" :label="t('moderation.warningMaxCount')">
				<NInputNumber
					v-model:value="editableItem.data.maxWarnings"
					:min="0"
					:max="10"
				/>
			</NFormItem>
		</div>

		<NDivider class="m-0 p-0" />

		<div class="form-block">
			<span>{{ t('moderation.excludedRoles') }}</span>
			<div v-if="editableItem?.data" class="flex flex-col gap-1">
				<NButtonGroup
					v-for="(group, index) of chunk(rolesSelectOptions.sort(), 5)"
					:key="index"
				>
					<NButton
						v-for="option of group"
						:key="option.value"
						:type="editableItem?.data?.excludedRoles.includes(option.value) ? 'success' : 'default'"
						secondary
						@click="() => {
							if (editableItem!.data!.excludedRoles.includes(option.value)) {
								editableItem!.data!.excludedRoles = editableItem!.data!.excludedRoles.filter(r => r !== option.value)
							}
							else {
								editableItem!.data!.excludedRoles.push(option.value)
							}
						}"
					>
						<template #icon>
							<IconSquareCheck v-if="editableItem?.data?.excludedRoles.includes(option.value)" />
							<IconSquare v-else />
						</template>
						{{ option.label }}
					</NButton>
				</NButtonGroup>
			</div>
		</div>

		<ModalDenylist
			v-if="editableItem?.data?.type === 'deny_list'"
			class="form-block"
		/>

		<NDivider class="m-0 p-0" />

		<NButton type="success" secondary @click="saveSettings">
			{{ t('sharedButtons.save') }}
		</NButton>
	</div>
</template>

<style scoped>
.form-block {
	@apply flex flex-col gap-2;
}

.form-block :deep(.n-input-number) {
	@apply w-full;
}
</style>
