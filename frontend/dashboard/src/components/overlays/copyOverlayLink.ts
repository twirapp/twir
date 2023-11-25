import { useNotification } from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useProfile } from '@/api';

export const useCopyOverlayLink = (overlayPath: string) => {
	const { data: profile } = useProfile();
	const { t } = useI18n();
	const messages = useNotification();

	const overlayLink = computed(() => {
		return `${window.location.origin}/overlays/${profile.value?.apiKey}/${overlayPath}`;
	});

	const copyOverlayLink = () => {
		navigator.clipboard.writeText(overlayLink.value);
		messages.success({
			title: t('overlays.copied'),
			duration: 5000,
		});
	};

	return {
		overlayLink,
		copyOverlayLink,
	};
};
