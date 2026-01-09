<script setup lang="ts">
import { CheckIcon, ClipboardIcon, EyeIcon, EyeOffIcon, InfoIcon } from 'lucide-vue-next'


import { useAlertsInformation } from '../composables/use-alerts-information'





const { t } = useI18n()
const { isCopied, copyOverlayLink, overlayLink, isShowOverlayLink, toggleShowOverlayLink } =
	useAlertsInformation()
</script>

<template>
	<UiAlert>
		<InfoIcon class="size-5" />
		<UiAlertTitle>
			{{ t('alerts.info') }}
		</UiAlertTitle>
		<UiAlertDescription>
			{{ t('alerts.overlayLabel') }}
			<div class="flex w-full items-center gap-2 mt-4">
				<div class="relative w-full">
					<UiInput
						class="pr-12 w-full"
						:type="isShowOverlayLink ? 'text' : 'password'"
						:default-value="overlayLink"
						readonly
					/>
					<UiButton
						variant="ghost"
						size="icon"
						class="absolute right-0 top-1/2 -translate-y-1/2"
						@click="toggleShowOverlayLink"
					>
						<EyeIcon v-if="isShowOverlayLink" class="size-4" />
						<EyeOffIcon v-else class="size-4" />
					</UiButton>
				</div>
				<UiButton size="icon" @click="copyOverlayLink">
					<ClipboardIcon v-if="!isCopied" class="size-4 min-w-10" />
					<CheckIcon v-else class="size-4 min-w-10" />
				</UiButton>
			</div>
		</UiAlertDescription>
	</UiAlert>
</template>
