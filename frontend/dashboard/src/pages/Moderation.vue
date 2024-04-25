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
import { ChannelRolePermissionEnum } from '@/gql/graphql';

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
const canEditModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModeration);

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
	<div class="flex items-center justify-center max-w-[1000px] mx-auto my-0">
		<n-grid cols="1 m:2" :x-gap="16" :y-gap="16" responsive="screen">
			<n-grid-item :span="1">
				<n-card
					class="select-none min-h-[300px] h-full hover:bg-[color:var(--hover-color)]"
					:style="{
						'--hover-color': theme.hoverColor,
						cursor: !canEditModeration ? 'not-allowed' : !isAddingNewItem ? 'pointer' : 'default'
					}"
				>
					<Transition mode="out-in">
						<div
							v-if="!isAddingNewItem"
							class="flex flex-col justify-center items-center h-full"
							@click="isAddingNewItem = true"
						>
							<IconSwords :size="45" />
							<span>{{ t('moderation.createNewRule') }}</span>
						</div>
						<div v-else class="flex flex-col gap-3">
							<div class="flex justify-between">
								<span>{{ t('moderation.createNewRule') }}</span>
								<n-button text size="tiny" @click="isAddingNewItem = false">
									<IconX />
								</n-button>
							</div>
							<div class="flex gap-2 flex-wrap">
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
											<div class="flex items-center gap-1">
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
			<div class="flex flex-col gap-[2px]">
				<span>
					{{ editableItem?.data ? t(`moderation.types.${editableItem.data.type}.name`) : 'Edit' }}
				</span>
				<span class="text-xs">
					{{ editableItem?.data ? t(`moderation.types.${editableItem.data.type}.description`) : '' }}
				</span>
			</div>
		</template>
		<modal v-if="editableItem" />
	</n-modal>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.1s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
