<script lang="ts" setup>
import { IconSwords, IconX } from '@tabler/icons-vue';
import { NGrid, NGridItem, NModal, NCard, useThemeVars, NButton, NTooltip } from 'naive-ui';
import { ref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

import { useModerationManager, useUserAccessFlagChecker } from '@/api/index.js';
import Card from '@/components/moderation/card.vue';
import { availableSettingsTypes, availableSettings, useEditableItem } from '@/components/moderation/helpers.js';
import { Icons } from '@/components/moderation/helpers.js';
import Modal from '@/components/moderation/modal.vue';

const manager = useModerationManager();
const { data: settings } = manager.getAll({});

const { editableItem } = useEditableItem();

const settingsOpened = ref(false);
function showSettings(id: string) {
	const item = settings.value?.body.find(i => i.id === id);
	if (!item) return;

	editableItem.value = toRaw(item);
	settingsOpened.value = true;
}

const theme = useThemeVars();
const { t } = useI18n();

const isAddingNewItem = ref(false);
const canEditModeration = useUserAccessFlagChecker('MANAGE_MODERATION');

async function createNewItem(itemType: string) {
	const defaultSettings = availableSettings.find(s => s.type === itemType);
	if (!defaultSettings) return;
	editableItem.value = {
		data: structuredClone(defaultSettings),
	};
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
					<Transition mode="out-in">
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
								<n-tooltip
									v-for="itemType of availableSettingsTypes"
									:key="itemType"
								>
									<template #trigger>
										<n-button
											secondary
											:disabled="!canEditModeration"
											@click="createNewItem(itemType)"
										>
											<div style="display: flex; align-items: center; gap: 4px">
												<component
													:is="Icons[itemType]"
													:size="20"
												/>
												<span>{{ t(`moderation.types.${itemType}.name`) }}</span>
											</div>
										</n-button>
									</template>
									{{ t(`moderation.types.${itemType}.description`) }}
								</n-tooltip>
							</div>
						</div>
					</Transition>
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
		:style="{
			width: '40vw',
			top: '0px',
		}"
		:on-close="() => settingsOpened = false"
	>
		<template #header>
			<div style="display: flex; flex-direction: column; gap: 2px">
				<span>
					{{ editableItem?.data ? t(`moderation.types.${editableItem.data.type}.name`) : 'Edit' }}
				</span>
				<span style="font-size: 12px;">
					{{ editableItem?.data ? t(`moderation.types.${editableItem.data.type}.description`) : '' }}
				</span>
			</div>
		</template>
		<modal v-if="editableItem" />
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

.v-enter-active,
.v-leave-active {
  transition: opacity 0.1s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
