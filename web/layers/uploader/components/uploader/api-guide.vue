<script setup lang="ts">
import { toast } from 'vue-sonner'

import UploaderCodeBlock from './code-block.vue'

const { t } = useI18n()
const clipboard = useClipboard()

// The landing layout (header-profile) already fetches authenticatedUser via
// callOnce on every landing page; for anonymous visitors the query errors
// gracefully and userWithoutDashboards stays undefined.
const userStore = useAuth()

// The token must never appear in SSR HTML of this public page: the gate stays
// closed during SSR and hydration, and flips only after client mount.
const isClientMounted = ref(false)
onMounted(() => {
	isClientMounted.value = true
})

const uploaderApiKey = computed(() => {
	if (!isClientMounted.value) return null
	return userStore.userWithoutDashboards?.apiKey ?? null
})

// Streamer safety: the token renders masked by default; the reveal toggle
// controls both visible spots at once.
const isTokenRevealed = ref(false)

function maskToken(token: string) {
	if (token.length <= 8) return '••••••••'
	return `${token.slice(0, 4)}••••••••${token.slice(-4)}`
}

const visibleApiKey = computed(() => {
	if (!uploaderApiKey.value) return null
	return isTokenRevealed.value ? uploaderApiKey.value : maskToken(uploaderApiKey.value)
})

const requestUrl = useRequestURL()
const origin = computed(() => requestUrl.origin)

const uploadEndpoint = computed(() => `${origin.value}/api/v1/uploader/files`)

const curlSnippet = computed(
	() => `curl "${uploadEndpoint.value}" \\
  -F "file=@/path/to/image.png"`
)

const curlResponseSnippet = computed(() =>
	JSON.stringify(
		{
			data: {
				id: 'AbC12345',
				name: 'image.png',
				type: 'image/png',
				ext: '.png',
				size: 1536221,
				link: `${origin.value}/u/AbC12345`,
				delete_link: `${origin.value}/api/v1/uploader/files/delete?key=DELETE_KEY&id=AbC12345`,
				created_at: '2026-08-08T12:00:00Z',
				expires_at: '2026-09-07T12:00:00Z',
			},
		},
		null,
		2
	)
)

const sharexSnippet = computed(() =>
	JSON.stringify(
		{
			Version: '16.0.1',
			Name: 'Twir',
			DestinationType: 'ImageUploader',
			RequestMethod: 'POST',
			RequestURL: uploadEndpoint.value,
			Body: 'MultipartFormData',
			FileFormName: 'file',
			URL: '{json:data.link}',
			DeletionURL: '{json:data.delete_link}',
		},
		null,
		2
	)
)

const chatterinoSnippet = computed(
	() => `Request URL: ${uploadEndpoint.value}
Form field: file
Image URL (JSON path): data.link
Deletion URL (JSON path): data.delete_link${
		visibleApiKey.value ? `\nExtra headers: api-key: ${visibleApiKey.value}` : ''
	}`
)

// Chatterino import format (src/util/ImageUploader.cpp exportSettings):
// native {property} dot-notation against the { "data": { ... } } envelope.
// Chatterino converts the Headers object into extra headers on import.
function buildChatterinoConfig(token: string | null) {
	return JSON.stringify(
		{
			Version: '1.0.0',
			Name: 'Twir Image Uploader',
			RequestMethod: 'POST',
			RequestURL: uploadEndpoint.value,
			Body: 'MultipartFormData',
			FileFormName: 'file',
			URL: '{data.link}',
			DeletionURL: '{data.delete_link}',
			...(token ? { Headers: { 'api-key': token } } : {}),
		},
		null,
		2
	)
}

// What the copy button writes: always the real token. Never render this one.
const chatterinoConfigSnippet = computed(() => buildChatterinoConfig(uploaderApiKey.value))
// What the code block shows: masked until revealed.
const chatterinoConfigDisplaySnippet = computed(() => buildChatterinoConfig(visibleApiKey.value))

const chatterinoConfigCopied = ref(false)
let chatterinoCopiedTimer: ReturnType<typeof setTimeout> | undefined

function copyChatterinoConfig() {
	clipboard.copy(chatterinoConfigSnippet.value)
	chatterinoConfigCopied.value = true
	clearTimeout(chatterinoCopiedTimer)
	chatterinoCopiedTimer = setTimeout(() => {
		chatterinoConfigCopied.value = false
	}, 2000)
	toast.success(t('uploader.copied'), {
		description: uploaderApiKey.value
			? t('uploader.guide.copyJsonSuccessAuthed')
			: t('uploader.guide.copyJsonSuccess'),
		duration: 2000,
	})
}

onUnmounted(() => {
	clearTimeout(chatterinoCopiedTimer)
})

type GuideTab = 'curl' | 'sharex' | 'chatterino'

const tabs: { id: GuideTab; labelKey: string }[] = [
	{ id: 'curl', labelKey: 'uploader.guide.tabs.curl' },
	{ id: 'sharex', labelKey: 'uploader.guide.tabs.sharex' },
	{ id: 'chatterino', labelKey: 'uploader.guide.tabs.chatterino' },
]

const activeTab = ref<GuideTab>('curl')
</script>

<template>
	<div
		class="flex flex-col w-full max-w-xl gap-3 rounded-2xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,9%)] p-4 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
	>
		<div class="flex items-center gap-2.5">
			<div class="flex rounded-lg border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,10%)] p-2">
				<Icon name="lucide:code-2" class="w-4 h-4 text-[hsl(240,11%,80%)]" />
			</div>
			<div>
				<h3 class="text-sm font-semibold text-[hsl(240,11%,90%)]">
					{{ t('uploader.guide.title') }}
				</h3>
				<p class="text-xs text-[hsl(240,11%,55%)]">
					{{ t('uploader.guide.description') }}
				</p>
			</div>
		</div>

		<div class="flex flex-wrap gap-2">
			<button
				v-for="tab in tabs"
				:key="tab.id"
				type="button"
				class="flex items-center rounded-lg border px-2.5 py-1.5 text-xs font-semibold transition-colors"
				:class="
					activeTab === tab.id
						? 'border-[hsl(240,11%,45%)] bg-[hsl(240,11%,30%)] text-white'
						: 'border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] text-[hsl(240,11%,60%)] hover:bg-[hsl(240,11%,20%)]'
				"
				@click="activeTab = tab.id"
			>
				{{ t(tab.labelKey) }}
			</button>
		</div>

		<div v-if="activeTab === 'curl'" class="flex flex-col gap-2">
			<p class="text-xs text-[hsl(240,11%,65%)]">
				{{ t('uploader.guide.curlDescription') }}
			</p>
			<UploaderCodeBlock :code="curlSnippet" />
			<p class="text-xs font-medium text-[hsl(240,11%,65%)]">
				{{ t('uploader.guide.curlResponse') }}
			</p>
			<UploaderCodeBlock :code="curlResponseSnippet" />
		</div>

		<div v-else-if="activeTab === 'sharex'" class="flex flex-col gap-2">
			<p class="text-xs text-[hsl(240,11%,65%)]">
				{{ t('uploader.guide.sharexDescription') }}
			</p>
			<UploaderCodeBlock :code="sharexSnippet" />
		</div>

		<div v-else class="flex flex-col gap-2">
			<p class="text-xs text-[hsl(240,11%,65%)]">
				{{ t('uploader.guide.chatterinoDescription') }}
			</p>
			<div
				v-if="uploaderApiKey"
				class="flex items-start gap-1.5 text-xs text-yellow-300/80"
			>
				<Icon name="lucide:shield-alert" class="h-3.5 w-3.5 flex-none mt-px" />
				<span class="flex-1">{{ t('uploader.guide.copyJsonAuthedNote') }}</span>
				<button
					type="button"
					class="flex-none flex items-center justify-center rounded-lg border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] p-1.5 text-[hsl(240,11%,80%)] hover:border-[hsl(240,11%,40%)] hover:bg-[hsl(240,11%,25%)] transition-colors"
					:aria-label="t(isTokenRevealed ? 'uploader.guide.hideToken' : 'uploader.guide.showToken')"
					:title="t(isTokenRevealed ? 'uploader.guide.hideToken' : 'uploader.guide.showToken')"
					@click="isTokenRevealed = !isTokenRevealed"
				>
					<Icon :name="isTokenRevealed ? 'lucide:eye-off' : 'lucide:eye'" class="h-3.5 w-3.5" />
				</button>
			</div>
			<p v-else class="flex items-start gap-1.5 text-xs text-[hsl(240,11%,55%)]">
				<Icon name="lucide:info" class="h-3.5 w-3.5 flex-none mt-px" />
				{{ t('uploader.guide.copyJsonLoginHint') }}
			</p>
			<button
				type="button"
				class="flex items-center justify-center gap-2 py-1.5 px-3 rounded-lg text-sm font-semibold border border-[hsl(240,11%,30%)] hover:border-[hsl(240,11%,45%)] bg-[hsl(240,11%,25%)] hover:bg-[hsl(240,11%,35%)] text-[hsl(240,11%,90%)] transition-colors"
				@click="copyChatterinoConfig"
			>
				<Icon :name="chatterinoConfigCopied ? 'lucide:check' : 'lucide:copy'" class="w-4 h-4" />
				{{ t('uploader.guide.copyJson') }}
			</button>
			<p class="text-xs text-[hsl(240,11%,55%)]">
				{{ t('uploader.guide.copyJsonHint') }}
			</p>
			<UploaderCodeBlock :code="chatterinoConfigDisplaySnippet" :copy-text="chatterinoConfigSnippet" />
			<UploaderCodeBlock :code="chatterinoSnippet" />
		</div>
	</div>
</template>
