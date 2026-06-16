<script setup lang="ts">
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

import { useAlertsInformation } from '../composables/use-alerts-information'

const { t } = useI18n()
const { isCopied, copyOverlayLink, overlayLink, isShowOverlayLink, toggleShowOverlayLink } =
	useAlertsInformation()
</script>

<template>
	<Alert>
		<Icon
			name="lucide:info"
			class="size-5"
		/>
		<AlertTitle>
			{{ t('alerts.info') }}
		</AlertTitle>
		<AlertDescription>
			{{ t('alerts.overlayLabel') }}
			<div class="mt-4 flex w-full items-center gap-2">
				<div class="relative w-full">
					<Input
						class="w-full pr-12"
						:type="isShowOverlayLink ? 'text' : 'password'"
						:default-value="overlayLink"
						readonly
					/>
					<Button
						variant="ghost"
						size="icon"
						class="absolute top-1/2 right-0 -translate-y-1/2"
						@click="toggleShowOverlayLink"
					>
						<Icon
							name="lucide:eye"
							v-if="isShowOverlayLink"
							class="size-4"
						/>
						<Icon
							name="lucide:eye-off"
							v-else
							class="size-4"
						/>
					</Button>
				</div>
				<Button
					size="icon"
					@click="copyOverlayLink"
				>
					<Icon
						name="lucide:clipboard"
						v-if="!isCopied"
						class="size-4 min-w-10"
					/>
					<Icon
						name="lucide:check"
						v-else
						class="size-4 min-w-10"
					/>
				</Button>
			</div>
		</AlertDescription>
	</Alert>
</template>
