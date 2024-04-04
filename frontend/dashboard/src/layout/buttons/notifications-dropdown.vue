<script setup lang="ts">
import { Notification } from '@twir/api/messages/notifications/notifications';
import { NTime, NText, NScrollbar } from 'naive-ui';

// import { useI18n } from 'vue-i18n';
// const { t } = useI18n();

defineProps<{ notifications: Notification[] }>();
</script>

<template>
	<div class="flex flex-col px-3 py-2">
		<n-scrollbar style="max-height: 400px" trigger="none">
			<div class="flex flex-col gap-6 mr-4">
				<div v-for="notification of notifications" :key="notification.id" class="flex flex-col gap-2">
					<n-text class="w-full break-words" v-html="notification.message" />
					<n-text :title="new Date(notification.createdAt).toLocaleString()" class="flex text-xs justify-end" :depth="3">
						<n-time type="relative" :time="new Date(notification.createdAt)" />
					</n-text>
				</div>
			</div>
		</n-scrollbar>

		<!-- <n-button>{{ t('navbar.seeAll') }}</n-button> -->
	</div>
</template>
