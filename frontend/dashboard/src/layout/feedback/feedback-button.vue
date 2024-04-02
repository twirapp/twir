<script setup lang="ts">

import { IconMessageShare } from '@tabler/icons-vue';
import { NButton } from 'naive-ui';
import { type ButtonProps } from 'naive-ui/es/button/src/Button';
import { storeToRefs } from 'pinia';
import { computed, h } from 'vue';
import { useI18n } from 'vue-i18n';

import FeedbackModal from './feedback-modal.vue';

import { useNaiveDiscrete } from '@/composables/use-naive-discrete';
import { useFeedbackForm } from '@/layout/feedback/feedback';
import { useSidebarCollapseStore } from '@/layout/use-sidebar-collapse';

withDefaults(defineProps<{ isDrawer: boolean }>(), {
	isDrawer: false,
});

const discrete = useNaiveDiscrete();
const feedbackFormStore = useFeedbackForm();
const { t } = useI18n();

const { form } = storeToRefs(feedbackFormStore);

const positiveButtonProps = computed<ButtonProps>(() => ({
	disabled: !form.value.message.length,
}));

function openFeedbackModal() {
	discrete.dialog.create({
		showIcon: false,
		content: () => h(FeedbackModal),
		negativeText: t('sharedButtons.cancel'),
		positiveText: t('sharedButtons.confirm'),
		// this works, but typings saying it not. Feel free to provide better and safe solution
		positiveButtonProps: positiveButtonProps as ButtonProps,
		async onPositiveClick() {
			await feedbackFormStore.$submit();
			discrete.notification.success({
				title: t('feedback.notification'),
				duration: 2500,
			});
			feedbackFormStore.$clearForm();
		},
	});
}

const collapsedStore = useSidebarCollapseStore();
const { isCollapsed } = storeToRefs(collapsedStore);
</script>

<template>
	<div class="flex px-2 mb-2">
		<n-button block secondary type="success" size="large" @click="openFeedbackModal">
			<template #icon>
				<IconMessageShare />
			</template>
			<template v-if="isDrawer || !isCollapsed">
				{{ t('feedback.button') }}
			</template>
		</n-button>
	</div>
</template>
