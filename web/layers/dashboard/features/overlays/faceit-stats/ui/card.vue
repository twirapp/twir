<script setup lang="ts">
import FaceitIcon from '~~/layers/dashboard/assets/integrations/faceit.svg'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'

const { t } = useI18n()
</script>

<template>
	<Card class="flex h-full w-full flex-col">
		<CardHeader>
			<CardTitle class="flex items-center gap-3 text-xl">
				<FaceitIcon class="size-12" />
				{{ t('overlays.faceit.secondTitle') }}
			</CardTitle>
		</CardHeader>

		<CardContent style="color: rgba(255, 255, 255, 0.52)">
			{{ t('overlays.faceit.description') }}
		</CardContent>

		<CardFooter class="mt-auto">
			<RouterLink
				v-slot="{ href, navigate }"
				custom
				to="/dashboard/overlays/faceit-stats"
			>
				<Button
					as="a"
					:href="href"
					variant="outline"
					class="flex items-center gap-2"
					@click="navigate"
				>
					<Icon
						name="lucide:settings"
						class="size-4"
					/>
					{{ t('overlays.buildWidget') }}
				</Button>
			</RouterLink>
		</CardFooter>
	</Card>
</template>
