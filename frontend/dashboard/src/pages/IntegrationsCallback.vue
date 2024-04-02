<script setup lang="ts">
import { NSpin } from 'naive-ui';
import { onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import {
	useSpotifyIntegration,
	useLastfmIntegration,
	useVKIntegration,
	useStreamlabsIntegration,
	useDonationAlertsIntegration,
	useFaceitIntegration,
	useDiscordIntegration, useValorantIntegration, useNightbotIntegration,
} from '@/api/index.js';

const router = useRouter();
const route = useRoute();

const integrationsHooks: {
	[x: string]: {
		manager: {
			usePostCode: (...args: any) => any | Promise<any>,
			useData?: () => {
				refetch: (...args: any) => any | Promise<any>
			}
		},
		closeWindow?: boolean,
	}
} = {
	'spotify': {
		manager: useSpotifyIntegration(),
		closeWindow: true,
	},
	'lastfm': {
		manager: useLastfmIntegration(),
		closeWindow: true,
	},
	'vk': {
		manager: useVKIntegration(),
		closeWindow: true,
	},
	'streamlabs': {
		manager: useStreamlabsIntegration(),
		closeWindow: true,
	},
	'donationalerts': {
		manager: useDonationAlertsIntegration(),
		closeWindow: true,
	},
	'faceit': {
		manager: useFaceitIntegration(),
		closeWindow: true,
	},
	'discord': {
		manager: useDiscordIntegration(),
		closeWindow: true,
	},
	'valorant': {
		manager: useValorantIntegration(),
		closeWindow: true,
	},
	'nightbot': {
		manager: useNightbotIntegration(),
		closeWindow: true,
	},
};

onMounted(async () => {
	const integrationName = route.params.integrationName;
	if (!integrationName || typeof integrationName !== 'string') {
		router.push({ name: 'Integrations' });
		return;
	}

	const integration = integrationsHooks[integrationName];
	const postCodeHook = integration?.manager?.usePostCode();
	const getDataHook = integration?.manager?.useData?.();

	const { code, token } = route.query;
	const incomingCode = code ?? token;

	if (typeof incomingCode !== 'string' || !postCodeHook) {
		if (integration?.closeWindow) {
			window.close();
		} else {
			router.push({ name: 'Integrations' });
		}
		return;
	}

	postCodeHook.mutateAsync({ code: incomingCode }).finally(async () => {
		if (integration?.closeWindow) {
			if (getDataHook) {
				await getDataHook.refetch({});
			}
			window.close();
		} else {
			router.push({ name: 'Integrations' });
		}
	});
});
</script>

<template>
	<div class="flex items-center justify-center w-full h-full bg-[#0f0f14]">
		<n-spin size="large" />
	</div>
</template>
