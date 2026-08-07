<script setup lang="ts">
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardTitle } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'

const overlayName = defineModel<string>('overlayName', {
	type: String,
	required: true,
})
const instaSave = defineModel<boolean>('instaSave', {
	type: Boolean,
	required: true,
})
const { t } = useI18n()
</script>

<template>
	<Card class="border-0 shadow-none p-0">
		<div class="border-b p-2">
			<CardTitle class="text-sm font-medium">{{ t('overlayBuilder.overlaySettings.title') }}</CardTitle>
		</div>
		<CardContent class="px-3 pb-3 space-y-3">
			<!-- Overlay Name -->
			<div class="space-y-1.5">
				<Label for="overlay-name" class="text-xs">
					{{ t('overlayBuilder.overlaySettings.name') }} <span class="text-destructive">*</span>
				</Label>
				<Input
					id="overlay-name"
					v-model="overlayName"
					:placeholder="t('overlayBuilder.overlaySettings.namePlaceholder')"
					maxlength="30"
					class="h-8 text-sm"
					@keydown.stop
				/>
				<p class="text-xs text-muted-foreground">
					{{ t('overlayBuilder.overlaySettings.characterCount', { count: overlayName.length }) }}
				</p>
			</div>

			<!-- Canvas Size (Fixed) -->
			<div class="space-y-1.5">
				<Label class="text-xs">{{ t('overlayBuilder.overlaySettings.canvasSize') }}</Label>
				<p class="text-sm text-muted-foreground">
					{{ t('overlayBuilder.overlaySettings.fullHd') }}
				</p>
			</div>

			<!-- Insta Save -->
			<div class="space-y-1.5">
				<div class="flex items-center justify-between">
					<Label for="insta-save" class="text-xs">
						{{ t('overlayBuilder.overlaySettings.instantSave') }}
					</Label>
					<Switch
						id="insta-save"
						:model-value="instaSave"
						@update:model-value="(val: boolean) => (instaSave = val)"
					/>
				</div>
				<p class="text-xs text-muted-foreground">
					{{ t('overlayBuilder.overlaySettings.instantSaveDescription') }}
				</p>
			</div>
		</CardContent>
	</Card>
</template>
