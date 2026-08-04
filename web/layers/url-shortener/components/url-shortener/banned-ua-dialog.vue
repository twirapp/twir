<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { toast } from 'vue-sonner'

import Button from '@/components/ui/button/Button.vue'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'

import { useUrlShortener } from '../../composables/use-url-shortener'

const props = defineProps<{
	open: boolean
	linkId: string
	shortUrl: string
}>()

const emit = defineEmits<{
	(e: 'update:open', value: boolean): void
}>()

const urlShortener = useUrlShortener()
const {
	bannedUaPresets,
	presetPatterns,
	linkPresets,
	isLinkPresetsLoading,
	perLinkBannedUserAgents,
	isPerLinkBannedUserAgentsLoading,
} = storeToRefs(urlShortener)

const selectedPresetId = ref('')
const isApplyingPreset = ref(false)
const removingPresetId = ref<string | null>(null)

const patternInput = ref('')
const descriptionInput = ref('')
const errorMessage = ref<string | null>(null)
const isCreatingPattern = ref(false)
const deletingPatternId = ref<string | null>(null)

const appliedLinkPresets = computed(() => linkPresets.value.get(props.linkId) ?? [])
const appliedPresetIds = computed(
	() => new Set(appliedLinkPresets.value.map((item) => item.preset_id))
)
const availablePresets = computed(() =>
	bannedUaPresets.value.filter((preset) => !appliedPresetIds.value.has(preset.id))
)
const linkPatterns = computed(() => perLinkBannedUserAgents.value.get(props.linkId) ?? [])

function presetName(presetId: string) {
	return bannedUaPresets.value.find((preset) => preset.id === presetId)?.name ?? 'Unknown preset'
}

function presetPatternCount(presetId: string) {
	return presetPatterns.value.get(presetId)?.length ?? 0
}

watch(
	() => props.open,
	async (isOpen) => {
		if (!isOpen || !import.meta.client) return

		selectedPresetId.value = ''
		errorMessage.value = null
		patternInput.value = ''
		descriptionInput.value = ''

		if (!bannedUaPresets.value.length) {
			await urlShortener.fetchBannedUaPresets()
		}

		await Promise.all([
			urlShortener.fetchLinkPresets(props.linkId),
			urlShortener.fetchPerLinkBannedUserAgents(props.linkId),
			...bannedUaPresets.value.map((preset) => urlShortener.fetchPresetPatterns(preset.id)),
		])
	}
)

async function handleApplyPreset() {
	if (!selectedPresetId.value) return

	isApplyingPreset.value = true
	errorMessage.value = null

	const { error } = await urlShortener.applyPresetToLink(props.linkId, selectedPresetId.value)
	if (error) {
		errorMessage.value = error.toString()
	} else {
		toast.success('Preset applied', {
			description: `All patterns from "${presetName(selectedPresetId.value)}" now apply to this link.`,
		})
		selectedPresetId.value = ''
	}

	isApplyingPreset.value = false
}

async function handleRemovePreset(presetId: string) {
	removingPresetId.value = presetId
	errorMessage.value = null

	const { error } = await urlShortener.removePresetFromLink(props.linkId, presetId)
	if (error) {
		errorMessage.value = error.toString()
	} else {
		toast.success('Preset removed')
	}

	removingPresetId.value = null
}

function isValidRegex(pattern: string): boolean {
	try {
		return Boolean(new RegExp(pattern))
	} catch {
		return false
	}
}

async function handleCreatePattern() {
	const pattern = patternInput.value.trim()
	if (!pattern) {
		errorMessage.value = 'Enter a regex pattern to continue'
		return
	}

	if (!isValidRegex(pattern)) {
		errorMessage.value = 'Invalid regex pattern'
		return
	}

	isCreatingPattern.value = true
	errorMessage.value = null

	const { error } = await urlShortener.createPerLinkBannedUserAgent(props.linkId, {
		pattern,
		description: descriptionInput.value.trim() || null,
	})

	if (error) {
		errorMessage.value = error.toString()
		isCreatingPattern.value = false
		return
	}

	patternInput.value = ''
	descriptionInput.value = ''
	toast.success('Pattern added', {
		description: 'User agents matching this pattern will receive a 404 response.',
	})
	isCreatingPattern.value = false
}

async function handleDeletePattern(id: string) {
	deletingPatternId.value = id
	errorMessage.value = null

	const { error } = await urlShortener.deletePerLinkBannedUserAgent(props.linkId, id)
	if (error) {
		errorMessage.value = error.toString()
	} else {
		toast.success('Pattern removed')
	}

	deletingPatternId.value = null
}

function closeDialog() {
	emit('update:open', false)
}
</script>

<template>
	<Dialog :open="open" @update:open="closeDialog">
		<DialogContent class="max-w-lg max-h-[85vh] overflow-y-auto">
			<DialogHeader>
				<DialogTitle>Banned User Agents</DialogTitle>
				<DialogDescription>
					Block specific clients (e.g. Chatterino) from seeing previews for
					<span class="font-mono">{{ shortUrl }}</span> using regex patterns.
				</DialogDescription>
			</DialogHeader>

			<div class="space-y-4">
				<div class="space-y-3">
					<h3 class="text-sm font-semibold">Presets applied to this link</h3>

					<p v-if="isLinkPresetsLoading.get(linkId)" class="text-sm text-[hsl(240,11%,70%)]">
						Loading presets...
					</p>
					<p v-else-if="!appliedLinkPresets.length" class="text-xs text-[hsl(240,11%,55%)]">
						No presets applied. Apply one below or manage presets on the main page.
					</p>
					<ul v-else class="space-y-2">
						<li
							v-for="linkPreset in appliedLinkPresets"
							:key="linkPreset.id"
							class="flex items-center justify-between gap-3 rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2"
						>
							<div class="min-w-0 flex-1">
								<p class="text-sm font-medium truncate">{{ presetName(linkPreset.preset_id) }}</p>
								<p class="text-xs text-[hsl(240,11%,55%)]">
									{{ presetPatternCount(linkPreset.preset_id) }} pattern{{ presetPatternCount(linkPreset.preset_id) === 1 ? '' : 's' }}
								</p>
							</div>
							<button
								:disabled="removingPresetId === linkPreset.preset_id"
								class="flex-shrink-0 rounded-md p-1 text-[hsl(240,11%,55%)] hover:text-red-400 hover:bg-[hsl(240,11%,18%)] transition-colors disabled:opacity-50"
								title="Remove preset from link"
								@click="handleRemovePreset(linkPreset.preset_id)"
							>
								<Icon
									v-if="removingPresetId === linkPreset.preset_id"
									name="lucide:loader-2"
									class="h-4 w-4 animate-spin"
								/>
								<Icon v-else name="lucide:x" class="h-4 w-4" />
							</button>
						</li>
					</ul>

					<div v-if="availablePresets.length" class="flex gap-2">
						<Select v-model="selectedPresetId">
							<SelectTrigger
								class="flex-1 border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,20%)] transition-colors"
							>
								<SelectValue placeholder="Select a preset" />
							</SelectTrigger>
							<SelectContent>
								<SelectItem v-for="preset in availablePresets" :key="preset.id" :value="preset.id">
									{{ preset.name }}
								</SelectItem>
							</SelectContent>
						</Select>
						<Button
							type="button"
							variant="outline"
							size="sm"
							class="border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,20%)]"
							:disabled="!selectedPresetId || isApplyingPreset"
							@click="handleApplyPreset"
						>
							<Icon v-if="isApplyingPreset" name="lucide:loader-2" class="h-4 w-4 mr-2 animate-spin" />
							Apply
						</Button>
					</div>
					<p v-else-if="bannedUaPresets.length" class="text-xs text-[hsl(240,11%,55%)]">
						All presets are already applied to this link.
					</p>
					<p v-else class="text-xs text-[hsl(240,11%,55%)]">
						No presets yet. Create one in the Banned User Agent Presets panel on the main page.
					</p>
				</div>

				<div class="space-y-3 border-t border-[hsl(240,11%,18%)] pt-4">
					<h3 class="text-sm font-semibold">Per-link banned user agents</h3>

					<p
						v-if="isPerLinkBannedUserAgentsLoading.get(linkId)"
						class="text-sm text-[hsl(240,11%,70%)]"
					>
						Loading patterns...
					</p>
					<p v-else-if="!linkPatterns.length" class="text-xs text-[hsl(240,11%,55%)]">
						No per-link patterns. These apply only to this link.
					</p>
					<ul v-else class="space-y-2">
						<li
							v-for="pattern in linkPatterns"
							:key="pattern.id"
							class="flex items-start justify-between gap-3 rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2"
						>
							<div class="min-w-0 flex-1">
								<p class="font-mono text-sm text-[hsl(240,11%,90%)] break-all">
									{{ pattern.pattern }}
								</p>
								<p v-if="pattern.description" class="text-xs text-[hsl(240,11%,55%)] mt-0.5">
									{{ pattern.description }}
								</p>
							</div>
							<button
								:disabled="deletingPatternId === pattern.id"
								class="flex-shrink-0 rounded-md p-1 text-[hsl(240,11%,55%)] hover:text-red-400 hover:bg-[hsl(240,11%,18%)] transition-colors disabled:opacity-50"
								title="Delete pattern"
								@click="handleDeletePattern(pattern.id)"
							>
								<Icon
									v-if="deletingPatternId === pattern.id"
									name="lucide:loader-2"
									class="h-4 w-4 animate-spin"
								/>
								<Icon v-else name="lucide:trash-2" class="h-4 w-4" />
							</button>
						</li>
					</ul>

					<div class="space-y-3">
						<div class="flex flex-col gap-2">
							<label class="text-xs text-[hsl(240,11%,60%)]" for="link-ua-pattern">
								Regex pattern <span class="text-red-400">*</span>
							</label>
							<input
								id="link-ua-pattern"
								v-model="patternInput"
								type="text"
								maxlength="512"
								placeholder="Chatterino|TwitchLib"
								class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
							/>
						</div>
						<div class="flex flex-col gap-2">
							<label class="text-xs text-[hsl(240,11%,60%)]" for="link-ua-description">
								Description <span class="text-[hsl(240,11%,45%)]">(optional)</span>
							</label>
							<input
								id="link-ua-description"
								v-model="descriptionInput"
								type="text"
								maxlength="256"
								placeholder="Block Chatterino previews"
								class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
							/>
						</div>
						<Button
							type="button"
							variant="outline"
							size="sm"
							class="border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,20%)]"
							:disabled="isCreatingPattern"
							@click="handleCreatePattern"
						>
							<Icon v-if="isCreatingPattern" name="lucide:loader-2" class="h-4 w-4 mr-2 animate-spin" />
							<Icon v-else name="lucide:plus" class="h-4 w-4 mr-2" />
							Add pattern
						</Button>
						<p class="text-xs text-[hsl(240,11%,55%)]">
							Use a valid JavaScript regex. Matching is case-insensitive against the full
							<span class="font-mono">User-Agent</span> header.
						</p>
					</div>
				</div>

				<p v-if="errorMessage" class="text-sm text-red-400">
					{{ errorMessage }}
				</p>
			</div>
		</DialogContent>
	</Dialog>
</template>
