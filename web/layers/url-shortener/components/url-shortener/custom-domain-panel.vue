<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { toast } from 'vue-sonner'

import { useUrlShortener } from '../../composables/use-url-shortener'
import { UserStoreKey } from '~/stores/user'

const userStore = useAuth()
await callOnce(UserStoreKey, () => userStore.getUserDataWithoutDashboards())

const urlShortener = useUrlShortener()
const { customDomain, isCustomDomainLoading } = storeToRefs(urlShortener)

const domainInput = ref('')
const errorMessage = ref<string | null>(null)
const isCreating = ref(false)
const isVerifying = ref(false)
const isDeleting = ref(false)
const showDeleteConfirm = ref(false)

const isAuthenticated = computed(() => Boolean(userStore.userWithoutDashboards))
const hasCustomDomain = computed(() => Boolean(customDomain.value))
const statusLabel = computed(() => {
	if (!customDomain.value) return 'Not connected'
	return customDomain.value.verified ? 'Verified' : 'Pending verification'
})
const statusClasses = computed(() => {
	if (!customDomain.value) {
		return 'border-[hsl(240,11%,24%)] text-[hsl(240,11%,70%)]'
	}

	return customDomain.value.verified
		? 'border-emerald-500/40 text-emerald-300'
		: 'border-yellow-500/40 text-yellow-300'
})
const verificationTarget = computed(
	() => customDomain.value?.verification_target ?? 'short-{token}.twir.app'
)
const dnsRecordDomain = computed(
	() => customDomain.value?.domain || domainInput.value || 'links.example.com'
)
const previewDomain = computed(() => customDomain.value?.domain || 'links.example.com')

const clipboard = useClipboard()

async function refreshCustomDomain() {
	errorMessage.value = null

	if (!isAuthenticated.value) {
		customDomain.value = null
		return
	}

	const { error } = await urlShortener.fetchCustomDomain()
	if (error) {
		errorMessage.value = error
	}
}

onMounted(async () => {
	if (!import.meta.client) return
	await refreshCustomDomain()
})

watch(
	() => isAuthenticated.value,
	(value) => {
		if (!value) {
			customDomain.value = null
			return
		}
		refreshCustomDomain()
	}
)

async function handleCreate() {
	const value = domainInput.value.trim()
	if (!value) {
		errorMessage.value = 'Enter a domain to continue'
		return
	}

	isCreating.value = true
	errorMessage.value = null

	const { error } = await urlShortener.createCustomDomain(value)
	if (error) {
		errorMessage.value = error
		isCreating.value = false
		return
	}

	domainInput.value = ''
	toast.success('Custom domain saved', {
		description: 'Add the CNAME record and verify it to start using the domain.',
	})
	isCreating.value = false
}

async function handleVerify() {
	isVerifying.value = true
	errorMessage.value = null

	const { error } = await urlShortener.verifyCustomDomain()
	if (error) {
		errorMessage.value = error
		isVerifying.value = false
		return
	}

	toast.success('Domain verified', {
		description: 'New short links will now use your custom domain.',
	})
	isVerifying.value = false
}

async function handleDelete() {
	isDeleting.value = true
	errorMessage.value = null

	const { error } = await urlShortener.deleteCustomDomain()
	if (error) {
		errorMessage.value = error
		isDeleting.value = false
		return
	}

	showDeleteConfirm.value = false
	toast.success('Custom domain deleted', {
		description: 'You can connect a new domain whenever you want.',
	})
	isDeleting.value = false
}

function handleCopyTarget() {
	clipboard.copy(verificationTarget.value)
	toast.success('Copied', {
		description: 'CNAME target copied to clipboard.',
	})
}

function handleCopyPreview() {
	clipboard.copy(`https://${previewDomain.value}/your-slug`)
	toast.success('Copied', {
		description: 'Example link copied to clipboard.',
	})
}
</script>

<template>
	<div
		class="flex w-full max-w-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,9%)] rounded-2xl p-4 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
	>
		<details class="group w-full">
			<summary class="flex items-start justify-between gap-4 cursor-pointer list-none">
				<div class="space-y-1">
					<h2 class="text-lg font-semibold">Custom Domain</h2>
					<p class="text-sm text-[hsl(240,11%,60%)]">
						Use your own domain for short links and share cleaner URLs.
					</p>
					<p v-if="!isAuthenticated" class="text-xs text-[hsl(240,11%,55%)]">
						Sign in to connect and manage custom domains.
					</p>
				</div>
				<div class="flex items-center gap-2">
					<span
						class="inline-flex items-center rounded-full border px-2.5 py-1 text-xs font-semibold"
						:class="statusClasses"
					>
						{{ statusLabel }}
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
						Custom domains are available only for authorized users.
					</p>
					<UiButton class="mt-3" @click="userStore.login">Login</UiButton>
				</div>

				<div v-else class="space-y-4">
					<div
						v-if="isCustomDomainLoading"
						class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4"
					>
						<p class="text-sm text-[hsl(240,11%,70%)]">Loading domain status...</p>
					</div>

					<div
						v-else-if="hasCustomDomain"
						class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4"
					>
						<div class="flex flex-wrap items-start justify-between gap-3">
							<div>
								<p class="text-xs uppercase tracking-wide text-[hsl(240,11%,60%)]">Current domain</p>
								<p class="text-lg font-semibold">{{ customDomain?.domain }}</p>
							</div>
							<UiButton variant="outline" size="sm" @click="handleCopyPreview">
								<Icon name="lucide:copy" class="mr-2 h-4 w-4" />
								Copy example
							</UiButton>
						</div>

						<div class="mt-4 space-y-2 text-sm text-[hsl(240,11%,70%)]">
							<div class="flex flex-col gap-1">
								<span class="text-xs uppercase tracking-wide text-[hsl(240,11%,55%)]">
									CNAME target
								</span>
								<span class="rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-2 py-1 font-mono text-xs text-[hsl(240,11%,85%)]">
									{{ verificationTarget }}
								</span>
							</div>
							<p>
								Short links will look like
								<span class="font-mono text-[hsl(240,11%,85%)]">
									https://{{ previewDomain }}/your-slug
								</span>
							</p>
						</div>

						<div class="mt-4 flex flex-wrap gap-2">
							<UiButton variant="outline" size="sm" @click="handleCopyTarget">
								<Icon name="lucide:copy" class="mr-2 h-4 w-4" />
								Copy CNAME target
							</UiButton>
							<UiButton
								v-if="customDomain && !customDomain.verified"
								size="sm"
								:disabled="isVerifying"
								@click="handleVerify"
							>
								<Icon v-if="isVerifying" name="lucide:loader-2" class="mr-2 h-4 w-4 animate-spin" />
								Verify DNS
							</UiButton>
							<UiButton
								v-if="!showDeleteConfirm"
								variant="destructive"
								size="sm"
								:disabled="isDeleting"
								@click="showDeleteConfirm = true"
							>
								Delete domain
							</UiButton>
							<div v-else class="flex flex-wrap gap-2">
								<UiButton variant="outline" size="sm" @click="showDeleteConfirm = false">
									Cancel
								</UiButton>
								<UiButton
									variant="destructive"
									size="sm"
									:disabled="isDeleting"
									@click="handleDelete"
								>
									<Icon v-if="isDeleting" name="lucide:loader-2" class="mr-2 h-4 w-4 animate-spin" />
									Confirm delete
								</UiButton>
							</div>
						</div>
					</div>

					<div v-else class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4">
						<label class="text-xs uppercase tracking-wide text-[hsl(240,11%,60%)]" for="custom-domain">
							Domain
						</label>
						<div class="mt-2 flex flex-col gap-2 sm:flex-row">
							<input
								id="custom-domain"
								v-model="domainInput"
								type="text"
								placeholder="links.example.com"
								class="flex-1 rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-sm text-white placeholder-[hsl(240,11%,45%)] focus-visible:outline-none focus-visible:ring focus-visible:ring-[hsl(240,11%,30%)]"
							/>
							<UiButton :disabled="isCreating" @click="handleCreate">
								<Icon v-if="isCreating" name="lucide:loader-2" class="mr-2 h-4 w-4 animate-spin" />
								Connect domain
							</UiButton>
						</div>
						<p class="mt-2 text-xs text-[hsl(240,11%,55%)]">
							Use a subdomain like <span class="font-mono">links.yoursite.com</span> for best results.
						</p>
						<p class="mt-1 text-xs text-[hsl(240,11%,55%)]">
							Only one custom domain can be connected per account.
						</p>
					</div>

					<p v-if="errorMessage" class="text-sm text-red-400">
						{{ errorMessage }}
					</p>
				</div>

				<div class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4">
					<h3 class="text-sm font-semibold">Setup guide</h3>
					<ol class="mt-2 space-y-2 text-sm text-[hsl(240,11%,70%)]">
						<li>1. Create a CNAME record for <span class="font-mono">{{ dnsRecordDomain }}</span>.</li>
						<li>2. Point it to <span class="font-mono">{{ verificationTarget }}</span>.</li>
						<li>3. Wait for DNS to propagate, then click Verify DNS.</li>
					</ol>
					<div
						class="mt-3 flex flex-col gap-1 rounded-lg border border-[hsl(240,11%,20%)] bg-[hsl(240,11%,15%)] px-3 py-2 text-xs font-mono text-[hsl(240,11%,80%)]"
					>
						<span>CNAME {{ dnsRecordDomain }}</span>
						<span>to {{ verificationTarget }}</span>
					</div>
				</div>
			</div>
		</details>
	</div>
</template>
