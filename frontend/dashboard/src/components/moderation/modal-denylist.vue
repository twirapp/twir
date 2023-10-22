<script setup lang="ts">
import { IconTrash } from '@tabler/icons-vue';
import { type ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NAlert, NInput, NDivider, NButton } from 'naive-ui';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
	item: ItemWithId
}>();

const { t } = useI18n();

const addItem = () => props.item.data!.denyList.push('');
const removeItem = (i: number) => props.item.data!.denyList = props.item.data!.denyList.filter((_, idx) => idx != i);
</script>

<template>
	<div>
		<n-divider style="margin: 0; padding: 0" />

		<n-alert v-if="!item.data?.denyList.length" type="warning">
			Add new badword
		</n-alert>

		<n-alert v-else type="info">
			Supports regex
		</n-alert>

		<div
			v-for="(_, i) of item.data!.denyList"
			:key="i"
			style="display: flex; gap: 4px;"
		>
			<n-input
				v-model:value="item.data!.denyList[i]"
				type="textarea"
				autosize
				:maxlength="500"
			/>
			<n-button text type="error" @click="removeItem(i)">
				<IconTrash />
			</n-button>
		</div>

		<n-button
			:disabled="item.data!.denyList.length >= 100"
			type="success"
			dashed
			block
			@click="addItem"
		>
			{{ t('sharedButtons.create') }}
		</n-button>
	</div>
</template>
