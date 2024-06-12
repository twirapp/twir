<script setup lang="ts">
import { NScrollbar, NText, NTime } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import BlocksRender from '@/components/ui/editorjs/blocks-render.vue'
import { useNotifications } from '@/composables/use-notifications'

const { t } = useI18n()

const { notifications } = useNotifications()
</script>

<template>
	<div class="flex flex-col pl-3 py-2">
		<NScrollbar style="max-height: 400px" trigger="none">
			<div class="flex flex-col gap-6 mr-4">
				<div v-if="notifications.length === 0">
					<NText>{{ t('adminPanel.notifications.emptyNotifications') }}</NText>
				</div>
				<div v-for="notification of notifications" :key="notification.id" class="flex flex-col gap-2">
					<NText v-if="notification.text" class="w-full break-words" v-html="notification.text" />
					<BlocksRender v-if="notification.editorJsJson" :data="notification.editorJsJson" />

					<NText :title="new Date(notification.createdAt).toLocaleString()" class="flex text-xs justify-end" :depth="3">
						<NTime type="relative" :time="new Date(notification.createdAt)" />
					</NText>
				</div>
			</div>
		</NScrollbar>

		<!-- <n-button>{{ t('navbar.seeAll') }}</n-button> -->
	</div>
</template>
