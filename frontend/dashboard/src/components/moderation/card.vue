<script setup lang="ts">
import { IconSettings, IconSword, IconTrash, IconLanguageOff, IconMessageOff, IconLinkOff, IconAbc, IconMoodOff, IconListLetters, IconAsteriskSimple } from '@tabler/icons-vue';
import { type ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NSwitch, NButton, NPopconfirm, useNotification } from 'naive-ui';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useModerationManager, useUserAccessFlagChecker } from '@/api';
import Card from '@/components/card/card.vue';

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

const userCanManageModeration = useUserAccessFlagChecker('MANAGE_MODERATION');

const icon = computed(() => {
	switch (props.item.data!.type) {
		case 'links': return IconLinkOff;
		case 'language': return IconLanguageOff;
		case 'deny_list': return IconListLetters;
		case 'long_message': return IconMessageOff;
		case 'caps': return IconAbc;
		case 'emotes': return IconMoodOff;
		case 'symbols': return IconAsteriskSimple;
		default: return IconSword;
	}
});

const message = useNotification();

const switchState = async (id: string, v: boolean) => {
	patchExecuting.value = true;

	try {
		await patcher.mutateAsync({ id, enabled: v });
		props.item.data!.enabled = v;

		const statusText = t(`sharedTexts.${v ? 'enabled' : 'disabled'}`).toLocaleLowerCase();

		message.success({
			title: `${t(`moderation.types.${props.item.data!.type}.name`)} ${statusText}`,
			duration: 1500,
		});
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
		:icon="icon"
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
			<div style="display: flex; gap: 8px;">
				<n-button
					:disabled="!userCanManageModeration"
					secondary
					size="large"
					@click="$emit('showSettings')"
				>
					<div style="display: flex; gap: 4px;">
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
							<div style="display: flex; gap: 4px;">
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
