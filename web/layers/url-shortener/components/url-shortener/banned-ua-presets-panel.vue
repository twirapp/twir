<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { toast } from 'vue-sonner'

import ActionConfirm from '@/components/ui/action-confirm/ActionConfirm.vue'

import { useUrlShortener } from '../../composables/use-url-shortener'
import { UserStoreKey } from '~/stores/user'

const userStore = useAuth()
await callOnce(UserStoreKey, () => userStore.getUserDataWithoutDashboards())

const urlShortener = useUrlShortener()
const { bannedUaPresets, isBannedUaPresetsLoading, presetPatterns, isPresetPatternsLoading } =
	storeToRefs(urlShortener)

const presetNameInput = ref('')
const presetDescriptionInput = ref('')
const errorMessage = ref<string | null>(null)
const isCreatingPreset = ref(false)

const editingPresetId = ref<string | null>(null)
const editingName = ref('')
const editingDescription = ref('')
const isSavingEdit = ref(false)

const presetPendingDeleteId = ref<string | null>(null)

const patternInputs = ref<Map<string, string>>(new Map())
const patternDescriptionInputs = ref<Map<string, string>>(new Map())
const patternErrors = ref<Map<string, string | null>>(new Map())
const creatingPatternPresetId = ref<string | null>(null)
const deletingPatternId = ref<string | null>(null)

const isAuthenticated = computed(() => Boolean(userStore.userWithoutDashboards))

const presetPendingDelete = computed(
	() => bannedUaPresets.value.find((preset) => preset.id === presetPendingDeleteId.value) ?? null
)

const isDeleteConfirmOpen = computed({
	get: () => presetPendingDeleteId.value !== null,
	set: (value: boolean) => {
		if (!value) {
			presetPendingDeleteId.value = null
		}
	},
})

async function refresh() {
	errorMessage.value = null

	if (!isAuthenticated.value) {
		bannedUaPresets.value = []
		presetPatterns.value = new Map()
		return
	}

	const { error } = await urlShortener.fetchBannedUaPresets()
	if (error) {
		errorMessage.value = error.toString()
		return
	}

	await Promise.all(
		bannedUaPresets.value.map((preset) => urlShortener.fetchPresetPatterns(preset.id))
	)
}

await refresh()

watch(
	() => isAuthenticated.value,
	(value) => {
		if (!value) {
			bannedUaPresets.value = []
			presetPatterns.value = new Map()
			return
		}
		refresh()
	}
)

async function handleCreatePreset() {
	const name = presetNameInput.value.trim()
	if (!name) {
		errorMessage.value = 'Enter a preset name to continue'
		return
	}

	isCreatingPreset.value = true
	errorMessage.value = null

	const { error } = await urlShortener.createBannedUaPreset({
		name,
		description: presetDescriptionInput.value.trim() || null,
	})

	if (error) {
		errorMessage.value = error.toString()
		isCreatingPreset.value = false
		return
	}

	presetNameInput.value = ''
	presetDescriptionInput.value = ''
	toast.success('Preset created', {
		description: 'Add regex patterns to it, then apply it to your links.',
	})
	isCreatingPreset.value = false
}

function startEdit(presetId: string) {
	const preset = bannedUaPresets.value.find((item) => item.id === presetId)
	if (!preset) return

	editingPresetId.value = presetId
	editingName.value = preset.name
	editingDescription.value = preset.description ?? ''
}

function cancelEdit() {
	editingPresetId.value = null
	editingName.value = ''
	editingDescription.value = ''
}

async function handleSaveEdit(presetId: string) {
	const name = editingName.value.trim()
	if (!name) {
		errorMessage.value = 'Preset name cannot be empty'
		return
	}

	isSavingEdit.value = true
	errorMessage.value = null

	const { error } = await urlShortener.updateBannedUaPreset(presetId, {
		name,
		description: editingDescription.value.trim() || null,
	})

	if (error) {
		errorMessage.value = error.toString()
		isSavingEdit.value = false
		return
	}

	toast.success('Preset updated')
	cancelEdit()
	isSavingEdit.value = false
}

async function handleDeletePreset() {
	const presetId = presetPendingDeleteId.value
	if (!presetId) return

	errorMessage.value = null

	const { error } = await urlShortener.deleteBannedUaPreset(presetId)
	if (error) {
		errorMessage.value = error.toString()
	} else {
		toast.success('Preset deleted')
	}

	if (editingPresetId.value === presetId) {
		cancelEdit()
	}

	presetPendingDeleteId.value = null
}

function setPatternInput(presetId: string, value: string) {
	patternInputs.value.set(presetId, value)
}

function setPatternDescriptionInput(presetId: string, value: string) {
	patternDescriptionInputs.value.set(presetId, value)
}

function isValidRegex(pattern: string): boolean {
	try {
		return Boolean(new RegExp(pattern))
	} catch {
		return false
	}
}

async function handleCreatePattern(presetId: string) {
	const pattern = (patternInputs.value.get(presetId) ?? '').trim()
	if (!pattern) {
		patternErrors.value.set(presetId, 'Enter a regex pattern to continue')
		return
	}

	if (!isValidRegex(pattern)) {
		patternErrors.value.set(presetId, 'Invalid regex pattern')
		return
	}

	creatingPatternPresetId.value = presetId
	patternErrors.value.set(presetId, null)

	const { error } = await urlShortener.createPresetPattern(presetId, {
		pattern,
		description: (patternDescriptionInputs.value.get(presetId) ?? '').trim() || null,
	})

	if (error) {
		patternErrors.value.set(presetId, error.toString())
		creatingPatternPresetId.value = null
		return
	}

	patternInputs.value.set(presetId, '')
	patternDescriptionInputs.value.set(presetId, '')
	toast.success('Pattern added', {
		description: 'User agents matching this pattern will receive a 404 response.',
	})
	creatingPatternPresetId.value = null
}

async function handleDeletePattern(presetId: string, patternId: string) {
	deletingPatternId.value = patternId

	const { error } = await urlShortener.deletePresetPattern(presetId, patternId)
	if (error) {
		patternErrors.value.set(presetId, error.toString())
	} else {
		toast.success('Pattern removed')
	}

	deletingPatternId.value = null
}
</script>

<template>
	<div
		class="flex w-full max-w-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,9%)] rounded-2xl p-4 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
	>
		<details class="group w-full">
			<summary class="flex items-start justify-between gap-4 cursor-pointer list-none">
				<div class="space-y-1">
					<h2 class="text-lg font-semibold">Banned User Agent Presets</h2>
					<p class="text-sm text-[hsl(240,11%,60%)]">
						Group regex patterns into named presets and apply them to your short links.
					</p>
					<p v-if="!isAuthenticated" class="text-xs text-[hsl(240,11%,55%)]">
						Sign in to manage banned user agent presets.
					</p>
				</div>
				<div class="flex items-center gap-2">
					<span
						v-if="isAuthenticated"
						class="inline-flex items-center rounded-full border px-2.5 py-1 text-xs font-semibold border-[hsl(240,11%,24%)] text-[hsl(240,11%,70%)]"
					>
						{{ bannedUaPresets.length }} preset{{ bannedUaPresets.length === 1 ? '' : 's' }}
					</span>
					<Icon
						name="lucide:chevron-down"
						class="h-4 w-4 text-[hsl(240,11%,55%)] transition-transform group-open:rotate-180"
					/>
				</div>
			</summary>

			<div class="mt-4 flex flex-col gap-4">
				<div
					v-if="!isAuthenticated"
					class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4"
				>
					<p class="text-sm text-[hsl(240,11%,70%)]">
						Banned user agent presets are available only for authorized users.
					</p>
					<UiButton class="mt-3" @click="userStore.login">Login</UiButton>
				</div>

				<div v-else class="space-y-4">
					<div
						v-if="isBannedUaPresetsLoading"
						class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4"
					>
						<p class="text-sm text-[hsl(240,11%,70%)]">Loading presets...</p>
					</div>

					<template v-else>
						<div
							v-for="preset in bannedUaPresets"
							:key="preset.id"
							class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4 space-y-3"
						>
							<div v-if="editingPresetId === preset.id" class="space-y-3">
								<div class="flex flex-col gap-2">
									<label class="text-xs text-[hsl(240,11%,60%)]" :for="`preset-edit-name-${preset.id}`">
										Name <span class="text-red-400">*</span>
									</label>
									<input
										:id="`preset-edit-name-${preset.id}`"
										v-model="editingName"
										type="text"
										maxlength="100"
										class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
									/>
								</div>
								<div class="flex flex-col gap-2">
									<label class="text-xs text-[hsl(240,11%,60%)]" :for="`preset-edit-description-${preset.id}`">
										Description <span class="text-[hsl(240,11%,45%)]">(optional)</span>
									</label>
									<input
										:id="`preset-edit-description-${preset.id}`"
										v-model="editingDescription"
										type="text"
										maxlength="256"
										class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
									/>
								</div>
								<div class="flex flex-wrap gap-2">
									<UiButton variant="outline" size="sm" :disabled="isSavingEdit" @click="cancelEdit">
										Cancel
									</UiButton>
									<UiButton size="sm" :disabled="isSavingEdit" @click="handleSaveEdit(preset.id)">
										<Icon
											v-if="isSavingEdit"
											name="lucide:loader-2"
											class="mr-2 h-4 w-4 animate-spin"
										/>
										Save
									</UiButton>
								</div>
							</div>

							<div v-else class="flex items-start justify-between gap-3">
								<div class="min-w-0 flex-1">
									<p class="font-semibold truncate">{{ preset.name }}</p>
									<p v-if="preset.description" class="text-xs text-[hsl(240,11%,55%)] mt-0.5">
										{{ preset.description }}
									</p>
								</div>
								<div class="flex items-center gap-1">
									<button
										class="flex-shrink-0 rounded-md p-1 text-[hsl(240,11%,55%)] hover:text-white hover:bg-[hsl(240,11%,18%)] transition-colors"
										title="Edit preset"
										@click="startEdit(preset.id)"
									>
										<Icon name="lucide:pencil" class="h-4 w-4" />
									</button>
									<button
										class="flex-shrink-0 rounded-md p-1 text-[hsl(240,11%,55%)] hover:text-red-400 hover:bg-[hsl(240,11%,18%)] transition-colors"
										title="Delete preset"
										@click="presetPendingDeleteId = preset.id"
									>
										<Icon name="lucide:trash-2" class="h-4 w-4" />
									</button>
								</div>
							</div>

							<div class="space-y-2">
								<p class="text-xs uppercase tracking-wide text-[hsl(240,11%,60%)]">Patterns</p>
								<p
									v-if="isPresetPatternsLoading.get(preset.id)"
									class="text-sm text-[hsl(240,11%,70%)]"
								>
									Loading patterns...
								</p>
								<p
									v-else-if="!(presetPatterns.get(preset.id) ?? []).length"
									class="text-xs text-[hsl(240,11%,55%)]"
								>
									No patterns yet. Add one below.
								</p>
								<ul v-else class="space-y-2">
									<li
										v-for="pattern in presetPatterns.get(preset.id) ?? []"
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
											@click="handleDeletePattern(preset.id, pattern.id)"
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
							</div>

							<div class="space-y-3 rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] p-3">
								<div class="flex flex-col gap-2">
									<label class="text-xs text-[hsl(240,11%,60%)]" :for="`preset-pattern-${preset.id}`">
										Regex pattern <span class="text-red-400">*</span>
									</label>
									<input
										:id="`preset-pattern-${preset.id}`"
										:value="patternInputs.get(preset.id) ?? ''"
										type="text"
										maxlength="512"
										placeholder="Chatterino|TwitchLib"
										class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
										@input="setPatternInput(preset.id, ($event.target as HTMLInputElement).value)"
									/>
								</div>
								<div class="flex flex-col gap-2">
									<label class="text-xs text-[hsl(240,11%,60%)]" :for="`preset-pattern-description-${preset.id}`">
										Description <span class="text-[hsl(240,11%,45%)]">(optional)</span>
									</label>
									<input
										:id="`preset-pattern-description-${preset.id}`"
										:value="patternDescriptionInputs.get(preset.id) ?? ''"
										type="text"
										maxlength="256"
										placeholder="Block Chatterino previews"
										class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
										@input="setPatternDescriptionInput(preset.id, ($event.target as HTMLInputElement).value)"
									/>
								</div>
								<UiButton
									size="sm"
									:disabled="creatingPatternPresetId === preset.id"
									@click="handleCreatePattern(preset.id)"
								>
									<Icon
										v-if="creatingPatternPresetId === preset.id"
										name="lucide:loader-2"
										class="mr-2 h-4 w-4 animate-spin"
									/>
									<Icon v-else name="lucide:plus" class="mr-2 h-4 w-4" />
									Add pattern
								</UiButton>
								<p v-if="patternErrors.get(preset.id)" class="text-sm text-red-400">
									{{ patternErrors.get(preset.id) }}
								</p>
								<p class="text-xs text-[hsl(240,11%,55%)]">
									Use a valid JavaScript regex. Matching is case-insensitive against the full
									<span class="font-mono">User-Agent</span> header.
								</p>
							</div>
						</div>

						<div class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4 space-y-3">
							<p class="text-xs uppercase tracking-wide text-[hsl(240,11%,60%)]">Create preset</p>
							<div class="flex flex-col gap-2">
								<label class="text-xs text-[hsl(240,11%,60%)]" for="preset-name">
									Name <span class="text-red-400">*</span>
								</label>
								<input
									id="preset-name"
									v-model="presetNameInput"
									type="text"
									maxlength="100"
									placeholder="Chat clients"
									class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
								/>
							</div>
							<div class="flex flex-col gap-2">
								<label class="text-xs text-[hsl(240,11%,60%)]" for="preset-description">
									Description <span class="text-[hsl(240,11%,45%)]">(optional)</span>
								</label>
								<input
									id="preset-description"
									v-model="presetDescriptionInput"
									type="text"
									maxlength="256"
									placeholder="Block known chat clients from link previews"
									class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
								/>
							</div>
							<UiButton :disabled="isCreatingPreset" @click="handleCreatePreset">
								<Icon v-if="isCreatingPreset" name="lucide:loader-2" class="mr-2 h-4 w-4 animate-spin" />
								<Icon v-else name="lucide:plus" class="mr-2 h-4 w-4" />
								Create preset
							</UiButton>
						</div>
					</template>

					<p v-if="errorMessage" class="text-sm text-red-400">
						{{ errorMessage }}
					</p>
				</div>
			</div>
		</details>

		<ClientOnly>
			<ActionConfirm
				v-model:open="isDeleteConfirmOpen"
				:confirm-text="`Delete preset '${presetPendingDelete?.name ?? ''}' and all its patterns? Links using it will lose this protection.`"
				@confirm="handleDeletePreset"
				@cancel="presetPendingDeleteId = null"
			/>
		</ClientOnly>
	</div>
</template>
