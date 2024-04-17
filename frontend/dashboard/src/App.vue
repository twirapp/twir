<script setup lang='ts'>
import { NSpin } from 'naive-ui';
import { ref } from 'vue';
import { RouterView, useRouter } from 'vue-router';
import { provideClient } from '@urql/vue';
import { urqlClient } from './plugins/urql.js';

const isRouterReady = ref(false);
const router = useRouter();
router.isReady().finally(() => isRouterReady.value = true);

provideClient(urqlClient);
</script>

<template>
	<div v-if="!isRouterReady" class="flex justify-center items-center h-full">
		<n-spin size="large" />
	</div>
	<router-view v-else />
</template>
