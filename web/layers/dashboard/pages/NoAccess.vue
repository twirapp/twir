<script setup lang="ts">
import { HomeIcon, ShieldAlertIcon } from 'lucide-vue-next'

import { useRouter } from 'vue-router'



const { t } = useI18n()
const router = useRouter()

function goToDashboard() {
	router.push('/dashboard')
}
</script>

<template>
	<div class="flex items-center justify-center h-full">
		<div class="flex flex-col items-center gap-4 text-center">
			<ShieldAlertIcon class="h-24 w-24 text-destructive" />
			<h1 class="text-2xl font-semibold">{{ t('haveNoAccess.title') }}</h1>
			<p class="text-muted-foreground max-w-md">
				{{ t('haveNoAccess.description') }}
			</p>
			<UiButton @click="goToDashboard" class="mt-4">
				<HomeIcon class="h-4 w-4 mr-2" />
				{{ t('sharedButtons.goToDashboard') }}
			</UiButton>
		</div>
	</div>
</template>
