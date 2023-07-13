<script setup lang='ts'>
import { onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import {
	useSpotifyIntegration,
	useLastfmIntegration,
	useVKIntegration,
	useStreamlabsIntegration, useDonationAlertsIntegration, useFaceitIntegration, createIntegrationOauth,
} from '@/api/index.js';

const router = useRouter();
const route = useRoute();

const integrationsHooks: { [x: string]: ReturnType<typeof createIntegrationOauth>} = {
	'spotify': useSpotifyIntegration(),
	'lastfm': useLastfmIntegration(),
	'vk': useVKIntegration(),
	'streamlabs': useStreamlabsIntegration(),
	'donationalerts': useDonationAlertsIntegration(),
	'faceit': useFaceitIntegration(),
};

onMounted(() => {
	const integrationName = route.params.integrationName;
	if (!integrationName || typeof integrationName !== 'string') {
		router.push({ name: 'Integrations' });
		return;
	}

	const integrationHook = integrationsHooks[integrationName];
	const postCodeHook = integrationHook?.usePostCode();

	const { code, token } = route.query;
	const incomingCode = code ?? token;

	if (typeof incomingCode !== 'string' || !postCodeHook) {
		router.push({ name: 'Integrations' });
		return;
	}

	postCodeHook.mutateAsync({ code: incomingCode }).finally(() => {
		router.push({ name: 'Integrations' });
	});
});
</script>

<template>
  <div></div>
</template>

<style scoped lang='postcss'>

</style>
