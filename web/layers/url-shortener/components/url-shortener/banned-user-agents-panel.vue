<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { toast } from 'vue-sonner'

import { useUrlShortener } from '../../composables/use-url-shortener'
import { UserStoreKey } from '~/stores/user'

const userStore = useAuth()
await callOnce(UserStoreKey, () => userStore.getUserDataWithoutDashboards())

const urlShortener = useUrlShortener()
const { globalBannedUserAgents, isGlobalBannedUserAgentsLoading } = storeToRefs(urlShortener)

const patternInput = ref('')
const descriptionInput = ref('')
const errorMessage = ref<string | null>(null)
const isCreating = ref(false)
const deletingId = ref<string | null>(null)

const isAuthenticated = computed(() => Boolean(userStore.userWithoutDashboards))

async function refresh() {
	errorMessage.value = null
	if (!isAuthenticated.value) return
	const { error } = await urlShortener.fetchGlobalBannedUserAgents()
	if (error) {
		errorMessage.value = error.toString()
	}
}

await refresh()

watch(
	() => isAuthenticated.value,
	(value) => {
		if (!value) {
			globalBannedUserAgents.value = []
			return
		}
		refresh()
	}
)

async function handleCreate() {
	const pattern = patternInput.value.trim()
	if (!pattern) {
		errorMessage.value = 'Enter a regex pattern to continue'
		return
	}

	try {
		new RegExp(pattern)
	} catch {
		errorMessage.value = 'Invalid regex pattern'
		return
	}

	isCreating.value = true
	errorMessage.value = null

	const { error } = await urlShortener.createGlobalBannedUserAgent({
		pattern,
		description: descriptionInput.value.trim() || null,
	})

	if (error) {
		errorMessage.value = error.toString()
		isCreating.value = false
		return
	}

	patternInput.value = ''
	descriptionInput.value = ''
	toast.success('Pattern added', {
		description: 'User agents matching this pattern will receive a 404 response.',
	})
	isCreating.value = false
}

async function handleDelete(id: string) {
	deletingId.value = id
	errorMessage.value = null

	const { error } = await urlShortener.deleteGlobalBannedUserAgent(id)
	if (error) {
		errorMessage.value = error.toString()
	} else {
		toast.success('Pattern removed')
	}

	deletingId.value = null
}
</script>

<template>
	<div
		class="flex w-full max-w-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,9%)] rounded-2xl p-4 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
	>
		<details class="group w-full">
			<summary class="flex items-start justify-between gap-4 cursor-pointer list-none">
				<div class="space-y-1">
					<h2 class="text-lg font-semibold">Banned User Agents</h2>
					<p class="text-sm text-[hsl(240,11%,60%)]">
						Block specific clients (e.g. Chatterino) from seeing link previews using regex patterns.
					</p>
					<p v-if="!isAuthenticated" class="text-xs text-[hsl(240,11%,55%)]">
						Sign in to manage banned user agent patterns.
					</p>
				</div>
				<div class="flex items-center gap-2">
					<span
						v-if="isAuthenticated"
						class="inline-flex items-center rounded-full border px-2.5 py-1 text-xs font-semibold border-[hsl(240,11%,24%)] text-[hsl(240,11%,70%)]"
					>
						{{ globalBannedUserAgents.length }} pattern{{ globalBannedUserAgents.length === 1 ? '' : 's' }}
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
						Banned user agent patterns are available only for authorized users.
					</p>
					<UiButton class="mt-3" @click="userStore.login">Login</UiButton>
				</div>

				<div v-else class="space-y-4">
					<div
						v-if="isGlobalBannedUserAgentsLoading"
						class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4"
					>
						<p class="text-sm text-[hsl(240,11%,70%)]">Loading patterns...</p>
					</div>

					<div
						v-else-if="globalBannedUserAgents.length > 0"
						class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4 space-y-2"
					>
						<p class="text-xs uppercase tracking-wide text-[hsl(240,11%,60%)]">Active patterns</p>
						<ul class="space-y-2">
							<li
								v-for="agent in globalBannedUserAgents"
								:key="agent.id"
								class="flex items-start justify-between gap-3 rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2"
							>
								<div class="min-w-0 flex-1">
									<p class="font-mono text-sm text-[hsl(240,11%,90%)] break-all">{{ agent.pattern }}</p>
									<p v-if="agent.description" class="text-xs text-[hsl(240,11%,55%)] mt-0.5">
										{{ agent.description }}
									</p>
								</div>
								<button
									:disabled="deletingId === agent.id"
									class="flex-shrink-0 rounded-md p-1 text-[hsl(240,11%,55%)] hover:text-red-400 hover:bg-[hsl(240,11%,18%)] transition-colors disabled:opacity-50"
									@click="handleDelete(agent.id)"
								>
									<Icon
										v-if="deletingId === agent.id"
										name="lucide:loader-2"
										class="h-4 w-4 animate-spin"
									/>
									<Icon v-else name="lucide:trash-2" class="h-4 w-4" />
								</button>
							</li>
						</ul>
					</div>

					<!-- Add new pattern -->
					<div class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4 space-y-3">
						<p class="text-xs uppercase tracking-wide text-[hsl(240,11%,60%)]">Add pattern</p>
						<div class="flex flex-col gap-2">
							<label class="text-xs text-[hsl(240,11%,60%)]" for="ua-pattern">
								Regex pattern <span class="text-red-400">*</span>
							</label>
							<input
								id="ua-pattern"
								v-model="patternInput"
								type="text"
								placeholder="Chatterino|TwitchLib"
								class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
							/>
						</div>
						<div class="flex flex-col gap-2">
							<label class="text-xs text-[hsl(240,11%,60%)]" for="ua-description">
								Description <span class="text-[hsl(240,11%,45%)]">(optional)</span>
							</label>
							<input
								id="ua-description"
								v-model="descriptionInput"
								type="text"
								placeholder="Block Chatterino previews"
								class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
							/>
						</div>
						<UiButton :disabled="isCreating" @click="handleCreate">
							<Icon v-if="isCreating" name="lucide:loader-2" class="mr-2 h-4 w-4 animate-spin" />
							<Icon v-else name="lucide:plus" class="mr-2 h-4 w-4" />
							Add pattern
						</UiButton>
						<p class="text-xs text-[hsl(240,11%,55%)]">
							Use a valid JavaScript regex. Matching is case-insensitive against the full
							<span class="font-mono">User-Agent</span> header.
						</p>
					</div>

					<p v-if="errorMessage" class="text-sm text-red-400">
						{{ errorMessage }}
					</p>
				</div>
			</div>
		</details>
	</div>
</template>
