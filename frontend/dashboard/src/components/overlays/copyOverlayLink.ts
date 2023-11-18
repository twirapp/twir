import { useMessage } from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useProfile } from '@/api';

export const useCopyOverlayLink = (overlayPath: string) => {
	const { data: profile } = useProfile();
	const { t } = useI18n();
	const messages = useMessage();

	const overlayLink = computed(() => {
		return `${window.location.origin}/overlays/${profile.value?.apiKey}/${overlayPath}`;
	});

	const copyOverlayLink = () => {
		navigator.clipboard.writeText(overlayLink.value);
		messages.success(t('overlays.copied'));
	};

	return {
		overlayLink,
		copyOverlayLink,
	};
};
