<script lang="ts" setup>
import { IconSwords, IconX } from '@tabler/icons-vue';
import { type ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NGrid, NGridItem, NModal, NCard, useThemeVars, NButton } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useModerationManager, useUserAccessFlagChecker } from '@/api/index.js';
import Card from '@/components/moderation/card.vue';
import { availableSettingsTypes, availableSettings } from '@/components/moderation/helpers.js';
import Modal from '@/components/moderation/modal.vue';

const manager = useModerationManager();
const { data: settings } = manager.getAll({});
const creator = manager.create;

const editableItem = ref<ItemWithId | undefined>();
const settingsOpened = ref(false);
function showSettings(id: string) {
	editableItem.value = settings.value?.body.find(i => i.id === id);
	settingsOpened.value = true;
}

const theme = useThemeVars();
const { t } = useI18n();

const isAddingNewItem = ref(false);
const canEditModeration = useUserAccessFlagChecker('MANAGE_MODERATION');

async function createNewItem(itemType: string) {
	const defaultSettings = availableSettings.find(s => s.type === itemType);
	if (!defaultSettings) return;
	const newItem = await creator.mutateAsync({ data: defaultSettings });
	editableItem.value = newItem;
	settingsOpened.value = true;
}
</script>

<template>
	<div
		style="
			display: flex;
			align-items: center;
			justify-content: center;
			max-width: 1000px;
			margin: 0 auto;
		"
	>
		<n-grid cols="1 m:2" :x-gap="16" :y-gap="16" responsive="screen">
			<n-grid-item :span="1">
				<n-card
					style="min-height: 300px; height: 100%"
					class="new-item-card"
					:style="{
						cursor: !canEditModeration ? 'not-allowed' : !isAddingNewItem ? 'pointer' : 'default'
					}"
				>
					<div
						v-if="!isAddingNewItem"
						class="new-item-block"
						@click="isAddingNewItem = true"
					>
						<IconSwords :size="45" />
						<span>{{ t('moderation.createNewRule') }}</span>
					</div>
					<div v-else style="display: flex; flex-direction: column; gap: 12px;">
						<div style="display: flex; justify-content: space-between;">
							<span>{{ t('moderation.createNewRule') }}</span>
							<n-button text size="tiny" @click="isAddingNewItem = false">
								<IconX />
							</n-button>
						</div>
						<div style="display: flex; gap: 8px; flex-wrap: wrap;">
							<n-button
								v-for="itemType of availableSettingsTypes"
								:key="itemType"
								secondary
								:disabled="!canEditModeration"
								@click="createNewItem(itemType)"
							>
								<div style="display: flex; align-items: center; gap: 4px">
									<component
										:is="availableSettings.find(i => i.type === itemType)?.icon"
										:size="20"
									/>
									<span>{{ t(`moderation.types.${itemType}.name`) }}</span>
								</div>
							</n-button>
						</div>
					</div>
				</n-card>
			</n-grid-item>
			<n-grid-item v-for="item of settings?.body" :key="item.id" :span="1">
				<card
					:item="item"
					@show-settings="showSettings(item.id)"
				/>
			</n-grid-item>
		</n-grid>
	</div>

	<n-modal
		v-model:show="settingsOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Edit settings"
		:style="{
			width: '400px',
			top: '0px',
		}"
		:on-close="() => settingsOpened = false"
	>
		<modal v-if="editableItem" :item="editableItem" />
	</n-modal>
</template>

<style scoped>
.new-item-block {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	height: 100%;
}

.new-item-card {
	-webkit-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

.new-item-card:hover {
	background-color: v-bind('theme.hoverColor');
}
</style>
