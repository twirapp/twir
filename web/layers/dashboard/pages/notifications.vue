<script setup lang="ts">
import { UseTimeAgo } from '@vueuse/components'
import { onMounted } from 'vue'

import { useNotifications } from '~/composables/use-notifications'
import PageLayout from '~/layout/page-layout.vue'

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
					<p class="text-muted-foreground">
						{{ t('adminPanel.notifications.emptyNotifications') }}
					</p>
				</div>

				<UiCard v-for="notification of notifications" :key="notification.id">
					<UiCardContent class="pt-6">
						<div
							v-if="notification.text"
							class="w-full break-words"
							v-html="notification.text"
						/>
						<UiBlocksRender v-if="notification.editorJsJson" :data="notification.editorJsJson" />

						<p
							:title="new Date(notification.createdAt).toLocaleString()"
							class="flex text-xs justify-end text-muted-foreground mt-2"
						>
							<UseTimeAgo v-slot="{ timeAgo }" :time="new Date(notification.createdAt)">
								{{ timeAgo }}
							</UseTimeAgo>
						</p>
					</UiCardContent>
				</UiCard>
			</div>
		</template>
	</PageLayout>
</template>
