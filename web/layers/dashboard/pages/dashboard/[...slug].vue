<script setup lang="ts">
import { Button } from '@/components/ui/button'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { t } = useI18n()
const router = useRouter()
const localePath = useLocalePath()

function goToDashboard() {
	router.push(localePath('/dashboard'))
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
				<Icon name="lucide:home" class="h-4 w-4 mr-2" />
				{{ t('sharedButtons.goToDashboard') }}
			</Button>
		</div>
	</div>
</template>
