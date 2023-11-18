<script setup lang='ts'>
import { NSpin } from 'naive-ui';
import { ref } from 'vue';
import { RouterView, useRouter } from 'vue-router';

const isRouterReady = ref(false);
const router = useRouter();

router.isReady().finally(() => isRouterReady.value = true);
</script>

<template>
	<div v-if="!isRouterReady" class="app-loader">
		<n-spin size="large" />
	</div>
	<router-view v-else />
</template>

<style scoped>
.app-loader {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100%;
}
</style>
