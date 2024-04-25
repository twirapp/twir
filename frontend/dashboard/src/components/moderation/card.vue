<script setup lang="ts">
import { IconSettings, IconTrash } from '@tabler/icons-vue';
import { type ItemWithId } from '@twir/api/messages/moderation/moderation';
import { NSwitch, NButton, NPopconfirm, useNotification } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { Icons } from './helpers.js';

import { useModerationManager, useUserAccessFlagChecker } from '@/api/index.js';
import Card from '@/components/card/card.vue';
import { ChannelRolePermissionEnum } from '@/gql/graphql';

const props = defineProps<{
	item: ItemWithId
}>();

const manager = useModerationManager();
const patcher = manager.patch!;
const deleter = manager.deleteOne;

const patchExecuting = ref(false);

defineEmits<{
	showSettings: []
}>();

const { t } = useI18n();

const userCanManageModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModeration);

const message = useNotification();

const switchState = async (id: string, v: boolean) => {
	patchExecuting.value = true;

	try {
		await patcher.mutateAsync({ id, enabled: v });
		props.item.data!.enabled = v;
	} catch (error) {
		console.error(error);
	} finally {
		patchExecuting.value = false;
	}
};

async function removeItem() {
	await deleter.mutateAsync({ id: props.item.id });
	message.success({
		title: t('sharedTexts.deleted'),
		duration: 2000,
	});
}
</script>

<template>
	<card
		:title="t(`moderation.types.${item.data!.type}.name`)"
		:icon="Icons[item.data!.type]"
		style="height:100%"
	>
		<template #headerExtra>
			<n-switch
				:disabled="!userCanManageModeration"
				:value="item.data!.enabled"
				:loading="patchExecuting"
				@update:value="(v) => switchState(item.id, v)"
			/>
		</template>

		<template #content>
			{{ t(`moderation.types.${item.data!.type}.description`) }}
		</template>

		<template #footer>
			<div class="flex gap-2">
				<n-button
					:disabled="!userCanManageModeration"
					secondary
					size="large"
					@click="$emit('showSettings')"
				>
					<div class="flex gap-1">
						<span>{{ t('sharedButtons.settings') }}</span>
						<IconSettings />
					</div>
				</n-button>
				<n-popconfirm
					:positive-text="t('deleteConfirmation.confirm')"
					:negative-text="t('deleteConfirmation.cancel')"
					@positive-click="removeItem"
				>
					<template #trigger>
						<n-button
							:disabled="!userCanManageModeration"
							secondary
							size="large"
							type="error"
						>
							<div class="flex gap-1">
								<span>{{ t('sharedButtons.delete') }}</span>
								<IconTrash />
							</div>
						</n-button>
					</template>
					{{ t('deleteConfirmation.text') }}
				</n-popconfirm>
			</div>
		</template>
	</card>
</template>
