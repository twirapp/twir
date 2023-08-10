<script setup lang="ts">
import { useStore } from '@nanostores/vue';

import { authLinkStore, profileStore } from '../../../stores/auth.js';

const profile = useStore(profileStore);
const authLink = useStore(authLinkStore);
</script>

<template>
	<a
		class="inline-flex px-7 py-4 font-semibold text-white rounded-lg bg-[#5D58F5] hover:bg-[#6964FF] transition-[background,box-shadow] text-lg focus-visible:outline-none focus-visible:ring-[#5D58F5]/50 focus-visible:ring-4"
		:href="!profile ? authLink : '/dashboard'"
	>
		<span v-if="!profile">Start for free</span>
		<span v-else>Dashboard</span>
	</a>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.9s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
