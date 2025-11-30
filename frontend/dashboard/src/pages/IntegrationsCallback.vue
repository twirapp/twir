<script setup lang="ts">
import { NSpin } from 'naive-ui'
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import {
	useFaceitIntegration,
	useLastfmIntegration,
	useStreamlabsIntegration,
	useVKIntegration,
} from '@/api/index.js'
import { useDiscordIntegration } from '@/features/integrations/composables/discord/use-discord-integration.js'

const router = useRouter()
const route = useRoute()

const discordIntegration = useDiscordIntegration()

const integrationsHooks: {
	[x: string]:
		| {
				manager: {
					usePostCode: (...args: any) => any | Promise<any>
					useData?: () => {
						refetch: (...args: any) => any | Promise<any>
					}
				}
				closeWindow?: boolean
		  }
		| {
				custom: true
				handler: (code: string) => Promise<void>
				closeWindow?: boolean
		  }
} = {
	lastfm: {
		manager: useLastfmIntegration(),
		closeWindow: true,
	},
	vk: {
		manager: useVKIntegration(),
		closeWindow: true,
	},
	streamlabs: {
		manager: useStreamlabsIntegration(),
		closeWindow: true,
	},
	faceit: {
		manager: useFaceitIntegration(),
		closeWindow: true,
	},
	discord: {
		custom: true,
		closeWindow: true,
		handler: async (code: string) => {
			await discordIntegration.connectGuild(code)
			window.opener?.postMessage('discord-connected', '*')
		},
	},
}

onMounted(async () => {
	const integrationName = route.params.integrationName
	if (!integrationName || typeof integrationName !== 'string') {
		router.push({ name: 'Integrations' })
		return
	}

	const integration = integrationsHooks[integrationName]

	const { code, token } = route.query
	const incomingCode = code ?? token

	if (typeof incomingCode !== 'string') {
		if (integration?.closeWindow) {
			window.close()
		} else {
			router.push({ name: 'Integrations' })
		}
		return
	}

	// Handle custom integrations (like Discord with GraphQL)
	if (integration && 'custom' in integration && integration.custom) {
		try {
			await integration.handler(incomingCode)
		} finally {
			if (integration.closeWindow) {
				window.close()
			} else {
				router.push({ name: 'Integrations' })
			}
		}
		return
	}

	// Handle legacy integrations
	if (integration && 'manager' in integration) {
		const postCodeHook = integration.manager.usePostCode()
		const getDataHook = integration.manager.useData?.()

		postCodeHook.mutateAsync({ code: incomingCode }).finally(async () => {
			if (integration.closeWindow) {
				if (getDataHook) {
					await getDataHook.refetch({})
				}
				window.close()
			} else {
				router.push({ name: 'Integrations' })
			}
		})
		return
	}

	// Fallback
	if (integration?.closeWindow) {
		window.close()
	} else {
		router.push({ name: 'Integrations' })
	}
})
</script>

<template>
	<div class="flex items-center justify-center w-full h-full bg-[#0f0f14]">
		<NSpin size="large" />
	</div>
</template>
