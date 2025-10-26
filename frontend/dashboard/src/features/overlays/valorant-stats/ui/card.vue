<script setup lang="ts">
import { SettingsIcon } from 'lucide-vue-next'

import ValorantIcon from '@/assets/integrations/valorant.svg?use'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
	<Card class="w-full h-full flex flex-col">
		<CardHeader>
			<CardTitle class="flex text-xl gap-4 items-center">
				<ValorantIcon class="size-11" />
				{{ t('overlays.valorant.secondTitle') }}
			</CardTitle>
		</CardHeader>

		<CardContent style="color: rgba(255, 255, 255, 0.52)">
			{{ t('overlays.valorant.description') }}
		</CardContent>

		<CardFooter class="mt-auto">
			<RouterLink v-slot="{ href, navigate }" custom :to="{ name: 'ValorantStatsOverlay' }">
				<Button
					as="a"
					:href="href"
					variant="outline"
					class="flex gap-2 items-center"
					@click="navigate"
				>
					<SettingsIcon class="size-4" />
					{{ t('overlays.buildWidget') }}
				</Button>
			</RouterLink>
		</CardFooter>
	</Card>
</template>
