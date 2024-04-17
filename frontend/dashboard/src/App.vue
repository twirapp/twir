<script setup lang='ts'>
import { NSpin } from 'naive-ui';
import { ref, watch } from 'vue';
import { RouterView, useRouter } from 'vue-router';
import { useUrqlClient } from './plugins/urql';
import { provideClient } from '@urql/vue';

const isRouterReady = ref(false);
const router = useRouter();
router.isReady().finally(() => isRouterReady.value = true);

const { urqlClient, createClient } = useUrqlClient();

watch(() => urqlClient.value, () => {
	console.log("Providing new urql client")

	// create and provide a new client is client.value is not set
	if (!urqlClient.value) {
		const newClient = createClient();
		urqlClient.value = newClient;
		provideClient(newClient);
		return;
	}

	// otherwise, provide the newly set client
	provideClient(urqlClient.value);
}, { immediate: true })
</script>

<template>
	<div v-if="!isRouterReady" class="flex justify-center items-center h-full">
		<n-spin size="large" />
	</div>
	<router-view v-else />
</template>
