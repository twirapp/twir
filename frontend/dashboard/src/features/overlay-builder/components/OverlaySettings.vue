<script setup lang="ts">
import { ref, watch } from 'vue'
import { ExternalLink, Plus, Save } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useMessage } from 'naive-ui'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from '@/components/ui/tooltip'

interface Props {
	overlayId?: string
	overlayName: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
	'update:overlayName': [name: string]
	save: []
	addLayer: []
}>()

const { t } = useI18n()
const message = useMessage()

const localName = ref(props.overlayName)

watch(() => props.overlayName, (newVal) => {
	localName.value = newVal
})

watch(localName, (newVal) => {
	emit('update:overlayName', newVal)
})

function handleSave() {
	if (!localName.value || localName.value.trim().length === 0) {
		message.error('Overlay name is required')
		return
	}

	if (localName.value.length > 30) {
		message.error('Overlay name must be less than 30 characters')
		return
	}

	emit('save')
}

function copyOverlayLink() {
	if (!props.overlayId) return

	const baseUrl = window.location.origin
	const overlayUrl = `${baseUrl}/overlays/${props.overlayId}`

	navigator.clipboard.writeText(overlayUrl).then(() => {
		message.success(t('sharedTexts.copied') || 'Link copied to clipboard!')
	}).catch(() => {
		message.error('Failed to copy link')
	})
}

// Keyboard shortcut for save
function handleKeydown(e: KeyboardEvent) {
	if ((e.ctrlKey || e.metaKey) && e.key === 's') {
		e.preventDefault()
		handleSave()
	}
}

if (typeof window !== 'undefined') {
	window.addEventListener('keydown', handleKeydown)
}
</script>

<template>
	<Card class="border-0 shadow-none">
		<CardHeader class="px-3 py-2">
			<CardTitle class="text-sm font-medium">Overlay Settings</CardTitle>
		</CardHeader>
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
				/>
				<p class="text-xs text-muted-foreground">
					{{ localName.length }}/30 characters
				</p>
			</div>

			<!-- Copy Link -->
			<div v-if="overlayId" class="flex items-center gap-2">
				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger as-child>
							<Button
								variant="outline"
								size="sm"
								class="flex-1 h-8 text-xs"
								@click="copyOverlayLink"
							>
								<ExternalLink class="h-3 w-3 mr-1.5" />
								Copy Link
							</Button>
						</TooltipTrigger>
						<TooltipContent>
							<p>{{ t('overlaysRegistry.copyLink') || 'Copy Overlay Link' }}</p>
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>
			</div>

			<!-- Save Button -->
			<Button
				class="w-full h-9"
				@click="handleSave"
			>
				<Save class="h-4 w-4 mr-2" />
				{{ t('sharedButtons.save') || 'Save' }}
			</Button>

			<!-- Add Layer Button -->
			<Button
				variant="outline"
				class="w-full h-9"
				@click="emit('addLayer')"
			>
				<Plus class="h-4 w-4 mr-2" />
				Add Layer
			</Button>
		</CardContent>
	</Card>
</template>
