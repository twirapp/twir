<script setup lang="ts">
import { HomeIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { Button } from '@/components/ui/button'

const { t } = useI18n()
const router = useRouter()

function goToDashboard() {
	router.push('/dashboard')
}
</script>

<template>
	<div class="flex items-center justify-center h-full">
		<div class="flex flex-col items-center gap-4 text-center">
			<div class="text-9xl font-bold text-muted-foreground">404</div>
			<h1 class="text-2xl font-semibold">{{ t('pageNotFound.title') }}</h1>
			<p class="text-muted-foreground max-w-md">
				{{ t('pageNotFound.description') }}
			</p>
			<Button @click="goToDashboard" class="mt-4">
				<HomeIcon class="h-4 w-4 mr-2" />
				{{ t('sharedButtons.goToDashboard') }}
			</Button>
		</div>
	</div>
</template>
