<script setup lang="ts">
import { NText, NTime } from 'naive-ui'
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

import { Card, CardContent } from '@/components/ui/card'
import BlocksRender from '@/components/ui/editorjs/blocks-render.vue'
import { useNotifications } from '@/composables/use-notifications'
import PageLayout from '@/layout/page-layout.vue'

const { t } = useI18n()

const { notifications, notificationsCounter } = useNotifications()

onMounted(() => {
	notificationsCounter.value.onRead()
})
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('adminPanel.notifications.title') }}
		</template>

		<template #content>
			<div class="flex flex-col gap-6 mr-4">
				<div v-if="notifications.length === 0">
					<NText>{{ t('adminPanel.notifications.emptyNotifications') }}</NText>
				</div>

				<Card v-for="notification of notifications" :key="notification.id">
					<CardContent class="pt-6">
						<NText
							v-if="notification.text"
							class="w-full break-words"
							v-html="notification.text"
						/>
						<BlocksRender v-if="notification.editorJsJson" :data="notification.editorJsJson" />

						<NText
							:title="new Date(notification.createdAt).toLocaleString()"
							class="flex text-xs justify-end"
							:depth="3"
						>
							<NTime type="relative" :time="new Date(notification.createdAt)" />
						</NText>
					</CardContent>
				</Card>
			</div>
		</template>
	</PageLayout>
</template>
