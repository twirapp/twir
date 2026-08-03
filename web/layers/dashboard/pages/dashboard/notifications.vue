<script setup lang="ts">
import { UseTimeAgo } from '@vueuse/components'
import { onMounted } from 'vue'
import { useNotifications } from '~~/layers/dashboard/composables/use-notifications'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

import { Card, CardContent } from '@/components/ui/card'
import BlocksRender from '@/components/ui/editorjs/blocks-render.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

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
			<div class="mr-4 flex flex-col gap-6">
				<div v-if="notifications.length === 0">
					<p class="text-muted-foreground">
						{{ t('adminPanel.notifications.emptyNotifications') }}
					</p>
				</div>

				<Card
					v-for="notification of notifications"
					:key="notification.id"
				>
					<CardContent class="pt-6">
						<div
							v-if="notification.text"
							class="w-full break-words"
							v-html="notification.text"
						/>
						<BlocksRender
							v-if="notification.editorJsJson"
							:data="notification.editorJsJson"
						/>

						<p
							:title="new Date(notification.createdAt).toLocaleString()"
							class="text-muted-foreground mt-2 flex justify-end text-xs"
						>
							<UseTimeAgo
								v-slot="{ timeAgo }"
								:time="new Date(notification.createdAt)"
							>
								{{ timeAgo }}
							</UseTimeAgo>
						</p>
					</CardContent>
				</Card>
			</div>
		</template>
	</PageLayout>
</template>
