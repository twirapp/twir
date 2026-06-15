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
			<Icon name="lucide:shield-alert" class="h-24 w-24 text-destructive" />
			<h1 class="text-2xl font-semibold">{{ t('haveNoAccess.title') }}</h1>
			<p class="text-muted-foreground max-w-md">
				{{ t('haveNoAccess.description') }}
			</p>
			<Button @click="goToDashboard" class="mt-4">
				<Icon name="lucide:home" class="h-4 w-4 mr-2" />
				{{ t('sharedButtons.goToDashboard') }}
			</Button>
		</div>
	</div>
</template>
