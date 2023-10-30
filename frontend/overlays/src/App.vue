<script setup lang="ts">
import { RouterView } from 'vue-router';

if (import.meta.env.DEV) {
	document.body.style.background = 'black';
}
</script>

<template>
	<router-view />
</template>
