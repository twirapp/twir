<script setup lang="ts">
import { CheckIcon, ClipboardIcon, EyeIcon, EyeOffIcon, InfoIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { useAlertsInformation } from '../composables/use-alerts-information'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

const { t } = useI18n()
const { isCopied, copyOverlayLink, overlayLink, isShowOverlayLink, toggleShowOverlayLink } =
	useAlertsInformation()
</script>

<template>
	<Alert>
		<InfoIcon class="size-5" />
		<AlertTitle>
			{{ t('alerts.info') }}
		</AlertTitle>
		<AlertDescription>
			{{ t('alerts.overlayLabel') }}
			<div class="flex w-full items-center gap-2 mt-4">
				<div class="relative w-full">
					<Input
						class="pr-12 w-full"
						:type="isShowOverlayLink ? 'text' : 'password'"
						:default-value="overlayLink"
						readonly
					/>
					<Button
						variant="ghost"
						size="icon"
						class="absolute right-0 top-1/2 -translate-y-1/2"
						@click="toggleShowOverlayLink"
					>
						<EyeIcon v-if="isShowOverlayLink" class="size-4" />
						<EyeOffIcon v-else class="size-4" />
					</Button>
				</div>
				<Button size="icon" @click="copyOverlayLink">
					<ClipboardIcon v-if="!isCopied" class="size-4 min-w-10" />
					<CheckIcon v-else class="size-4 min-w-10" />
				</Button>
			</div>
		</AlertDescription>
	</Alert>
</template>
