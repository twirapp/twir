<script setup lang="ts">


const overlayName = defineModel<string>('overlayName', {
	type: String,
	required: true,
})
const instaSave = defineModel<boolean>('instaSave', {
	type: Boolean,
	required: true,
})
</script>

<template>
	<UiCard class="border-0 shadow-none p-0">
		<div class="border-b p-2">
			<UiCardTitle class="text-sm font-medium">Overlay Settings</UiCardTitle>
		</div>
		<UiCardContent class="px-3 pb-3 space-y-3">
			<!-- Overlay Name -->
			<div class="space-y-1.5">
				<UiLabel for="overlay-name" class="text-xs">
					Name <span class="text-destructive">*</span>
				</UiLabel>
				<UiInput
					id="overlay-name"
					v-model="overlayName"
					placeholder="My Overlay"
					maxlength="30"
					class="h-8 text-sm"
					@keydown.stop
				/>
				<p class="text-xs text-muted-foreground">
					{{ overlayName.length }}/30 characters
				</p>
			</div>

			<!-- Canvas Size (Fixed) -->
			<div class="space-y-1.5">
				<UiLabel class="text-xs">Canvas Size</UiLabel>
				<p class="text-sm text-muted-foreground">
					1920 × 1080 (Full HD)
				</p>
			</div>

			<!-- Insta Save -->
			<div class="space-y-1.5">
				<div class="flex items-center justify-between">
					<UiLabel for="insta-save" class="text-xs">
						Instant Save
					</UiLabel>
					<UiSwitch
						id="insta-save"
						:model-value="instaSave"
						@update:model-value="(val: boolean) => (instaSave = val)"
					/>
				</div>
				<p class="text-xs text-muted-foreground">
					Automatically save position and rotation changes in real-time
				</p>
			</div>
		</UiCardContent>
	</UiCard>
</template>
