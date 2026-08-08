<script setup lang="ts">
import ModuleChatTranslations from '~~/layers/dashboard/features/modules/ui/module-chat-translations.vue'
import ModuleDota from '~~/layers/dashboard/features/modules/ui/module-dota.vue'
import ModuleWebhookNotifications from '~~/layers/dashboard/features/modules/ui/module-webhook-notifications.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'
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
				<ModuleDota />
			</div>
		</template>
	</PageLayout>
</template>
