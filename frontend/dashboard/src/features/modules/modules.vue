<script setup lang="ts">
import ModuleChatTranslations from '@/features/modules/ui/module-chat-translations.vue'
import PageLayout from '@/layout/page-layout.vue'
</script>

<template>
	<PageLayout>
		<template #title>
			Modules
		</template>

		<template #content>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<ModuleChatTranslations />
			</div>
		</template>
	</PageLayout>
</template>
