<script lang="ts" setup>
import type {
	Settings,
} from '@twir/grpc/generated/api/api/overlays_kappagen';
import { useNotification } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const formValue = ref<Settings>({});

const message = useNotification();
const { t } = useI18n();
const emote: Emote = {
	url: 'https://cdn.7tv.app/emote/6548b7074789656a7be787e1/4x.webp',
	zwe: [
		{
			url: 'https://cdn.7tv.app/emote/6128ed55a50c52b1429e09dc/4x.webp',
		},
	],
};
async function save() {
	if (!formValue.value) return;

	// await updater.mutateAsync(formValue.value);
	message.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	});
}
</script>

<template>
	<div>
		settings of kappagen
	</div>
</template>
