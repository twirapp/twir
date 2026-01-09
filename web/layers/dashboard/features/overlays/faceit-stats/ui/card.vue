<script setup lang="ts">
import { SettingsIcon } from 'lucide-vue-next'

import FaceitIcon from '~/assets/integrations/faceit.svg?use'




const { t } = useI18n()
</script>

<template>
	<UiCard class="w-full h-full flex flex-col">
		<UiCardHeader>
			<UiCardTitle class="flex text-xl gap-3 items-center">
				<FaceitIcon class="size-12" />
				{{ t('overlays.faceit.secondTitle') }}
			</UiCardTitle>
		</UiCardHeader>

		<UiCardContent style="color: rgba(255, 255, 255, 0.52)">
			{{ t('overlays.faceit.description') }}
		</UiCardContent>

		<UiCardFooter class="mt-auto">
			<RouterLink v-slot="{ href, navigate }" custom to="/dashboard/overlays/faceit-stats">
				<UiButton
					as="a"
					:href="href"
					variant="outline"
					class="flex gap-2 items-center"
					@click="navigate"
				>
					<SettingsIcon class="size-4" />
					{{ t('overlays.buildWidget') }}
				</UiButton>
			</RouterLink>
		</UiCardFooter>
	</UiCard>
</template>
