<script setup lang='ts'>
import { NSpin } from 'naive-ui';
import { provide, ref } from 'vue';
import { RouterView, useRouter } from 'vue-router';
import { useUrqlClient } from './plugins/urql';

const isRouterReady = ref(false);
const router = useRouter();
router.isReady().finally(() => isRouterReady.value = true);

const { urqlClient, createClient } = useUrqlClient();

if (!urqlClient.value) {
  createClient();
}

provide('$urql', urqlClient);
</script>

<template>
	<div v-if="!isRouterReady" class="flex justify-center items-center h-full">
		<n-spin size="large" />
	</div>
	<router-view v-else />
</template>
