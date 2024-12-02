<script setup lang="ts">
import { SettingsIcon } from 'lucide-vue-next'

import FaceitIcon from '@/assets/integrations/faceit.svg?use'
import { Button } from '@/components/ui/button'
import { Card, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
</script>

<template>
	<Card class="w-full h-full flex flex-col">
		<CardHeader>
			<CardTitle class="flex gap-2 flex-col ">
				<FaceitIcon class="size-12" />
				Faceit Stats
			</CardTitle>
		</CardHeader>

		<CardFooter class="mt-auto">
			<RouterLink v-slot="{ href, navigate }" custom to="/dashboard/overlays/faceit-stats">
				<Button
					as="a"
					:href="href"
					variant="outline"
					class="flex gap-2 items-center"
					@click="navigate"
				>
					<SettingsIcon class="size-4" />
					Build widget
				</Button>
			</RouterLink>
		</CardFooter>
	</Card>
</template>
