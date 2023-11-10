<script lang="ts" setup>
import type {
	Settings,
} from '@twir/grpc/generated/api/api/overlays_kappagen';
import { useNotification } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import KappagenOverlay from 'kappagen'
import type { Emote, KappagenEmoteConfig } from 'kappagen'

const formValue = ref<Settings>({
	emotes: {
		time: 5,
		max: 0,
		queue: 0,
	},
	animations: [],
	enableRave: false,
	animation: {
		fadeIn: true,
		fadeOut: true,
		zoomIn: true,
		zoomOut: true,
	},
	cube: {
		speed: 6,
	},
	size: {
		// from 7 to 20
		ratioNormal: 7,
		// from 14 to 40
		ratioSmall: 14,
		min: 1,
		max: 256,
	},
});

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
