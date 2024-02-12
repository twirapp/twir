<script setup lang="ts">
import { NButton } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useNowPlayingForm } from './use-now-playing-form';

import {
	useNowPlayingOverlayManager,
	useProfile,
	useUserAccessFlagChecker,
} from '@/api';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';
import { useNowPlayingIframe } from '@/pages/overlays/now-playing/now-playing-iframe';

const { t } = useI18n();

const discrete = useNaiveDiscrete();
const { copyOverlayLink } = useCopyOverlayLink('nowplaying');
const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const { data: profile } = useProfile();

const { data: formValue } = storeToRefs(useNowPlayingForm());

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays;
});

const manager = useNowPlayingOverlayManager();
const updater = manager.useUpdate();

async function save() {
	if (!formValue.value?.id) return;

	await updater.mutateAsync({
		id: formValue.value.id,
		preset: formValue.value.preset,
	});

	discrete.notification.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	});
}

const iframeStore = useNowPlayingIframe();

watch(formValue, (form) => {
	if (!form) return;
	console.log('sending settings');
	iframeStore.sendIframeMessage('settings', form);
}, { immediate: true, deep: true });
</script>

<template>
	<div v-if="formValue" class="card">
		<div class="card-header">
			<n-button
				secondary
				type="info"
				:disabled="!formValue.id || !canCopyLink"
				@click="copyOverlayLink({ id: formValue.id! })"
			>
				{{ t('overlays.copyOverlayLink') }}
			</n-button>
			<n-button secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</div>

		<div class="card-body-column">
			qwe
			{{ formValue }}
		</div>
	</div>
</template>
