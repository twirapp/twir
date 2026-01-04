<script setup lang="ts">
import { ref, watch } from 'vue'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardTitle } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'

interface Props {
	overlayName: string
	instaSave: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
	'update:overlayName': [name: string]
	'update:instaSave': [enabled: boolean]
}>()

const localName = ref(props.overlayName)
const localInstaSave = ref(props.instaSave)

watch(() => props.overlayName, (newVal) => {
	localName.value = newVal
})

watch(() => props.instaSave, (newVal) => {
	localInstaSave.value = newVal
})

watch(localName, (newVal) => {
	emit('update:overlayName', newVal)
})

watch(localInstaSave, (newVal) => {
	emit('update:instaSave', newVal)
})
</script>

<template>
	<Card class="border-0 shadow-none p-0">
		<div class="border-b p-2">
			<CardTitle class="text-sm font-medium">Overlay Settings</CardTitle>
		</div>
		<CardContent class="px-3 pb-3 space-y-3">
			<!-- Overlay Name -->
			<div class="space-y-1.5">
				<Label for="overlay-name" class="text-xs">
					Name <span class="text-destructive">*</span>
				</Label>
				<Input
					id="overlay-name"
					v-model="localName"
					placeholder="My Overlay"
					maxlength="30"
					class="h-8 text-sm"
					@keydown.stop
				/>
				<p class="text-xs text-muted-foreground">
					{{ localName.length }}/30 characters
				</p>
			</div>

			<!-- Canvas Size (Fixed) -->
			<div class="space-y-1.5">
				<Label class="text-xs">Canvas Size</Label>
				<p class="text-sm text-muted-foreground">
					1920 × 1080 (Full HD)
				</p>
			</div>

			<!-- Insta Save -->
			<div class="space-y-1.5">
				<div class="flex items-center justify-between">
					<Label for="insta-save" class="text-xs">
						Instant Save
					</Label>
					<Switch
						id="insta-save"
						v-model:checked="localInstaSave"
					/>
				</div>
				<p class="text-xs text-muted-foreground">
					Automatically save position and rotation changes in real-time
				</p>
			</div>
		</CardContent>
	</Card>
</template>
