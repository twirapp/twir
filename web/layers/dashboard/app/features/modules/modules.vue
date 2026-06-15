<script setup lang="ts">
import ModuleChatTranslations from '~~/layers/dashboard/app/features/modules/ui/module-chat-translations.vue'
import ModuleWebhookNotifications from '~~/layers/dashboard/app/features/modules/ui/module-webhook-notifications.vue'
import PageLayout from '~~/layers/dashboard/app/layout/page-layout.vue'
</script>

<template>
	<PageLayout>
		<template #title>
			Modules
		</template>

		<template #content>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<ModuleChatTranslations />
				<ModuleWebhookNotifications />
			</div>
		</template>
	</PageLayout>
</template>
