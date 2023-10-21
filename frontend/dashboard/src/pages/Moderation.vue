<script lang="ts" setup>
import { type ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NGrid, NGridItem, NModal } from 'naive-ui';
import { ref } from 'vue';

import { useModerationManager } from '@/api';
import Card from '@/components/moderation/card.vue';
import Modal from '@/components/moderation/modal.vue';

const manager = useModerationManager();
const { data: settings } = manager.getAll({});

const editableItem = ref<ItemWithId | undefined>();
const settingsOpened = ref(false);
function showSettings(id: string) {
	editableItem.value = settings.value?.body.find(i => i.id === id);
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
		<modal :item="editableItem" />
	</n-modal>
</template>
