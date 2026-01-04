<script setup lang="ts">
import { ref, watch } from 'vue'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardTitle } from '@/components/ui/card'

interface Props {
	overlayName: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
	'update:overlayName': [name: string]
}>()

const localName = ref(props.overlayName)

watch(() => props.overlayName, (newVal) => {
	localName.value = newVal
})

watch(localName, (newVal) => {
	emit('update:overlayName', newVal)
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
		</CardContent>
	</Card>
</template>
