<script setup lang="ts">
import UiButton from '~~/layers/landing/components/landing-ui-button.vue'

definePageMeta({
	title: 'Login',
	layout: 'clean',
})

const url = new URL(window.location.href)
const code = url.searchParams.get('code')
const state = url.searchParams.get('state')
const error = ref(url.searchParams.get('error'))
const loading = ref(true)

onMounted(async () => {
	if (import.meta.server) return

	if (error.value) {
		return
	}

	if (!code || !state) {
		error.value = `[youtube] Something unexpected happened, because authorization code wasn't provided. Please try to log in again`
		return
	}

	try {
		const res = await $fetch<{ data: { redirect_to: string } }>(
			`${window.location.origin}/api/auth/youtube/code`,
			{
				method: 'POST',
				body: { code, state },
				credentials: 'include',
			}
		)

		window.location.replace(res.data.redirect_to)
	} catch (requestError) {
		console.error(requestError)
		error.value = 'Internal error happened, please contact devs in discord'
	} finally {
		loading.value = false
	}
})
</script>

<template>
	<span
		class="purple-gradient pointer-events-none absolute -top-[220px] right-0 left-0 -z-20 mx-auto h-[482px] rounded-full content-['']"
	></span>
	<div class="flex h-screen items-center justify-center px-3">
		<div
			v-if="error"
			class="flex flex-col items-center gap-2"
		>
			<span class="text-center font-medium text-red-400">{{ error }}</span>
			<NuxtLink
				v-slot="{ navigate, href }"
				to="/"
				custom
			>
				<UiButton
					:href="href!"
					variant="primary"
					role="link"
					@click="navigate"
				>
					Back to home
				</UiButton>
			</NuxtLink>
		</div>

		<div
			v-else-if="loading"
			role="status"
		>
			<Icon name="lucide:loader-circle" class="animate-spin size-12 text-[#FF0000]" />
			<span class="sr-only">Loading...</span>
		</div>
	</div>
</template>
