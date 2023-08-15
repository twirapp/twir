<script setup lang="ts">
import { useStore } from '@nanostores/vue';

import { authLinkStore, profileStore } from '../../stores/auth';

const profile = useStore(profileStore);
const authLink = useStore(authLinkStore);
</script>

<template>
	<a
		class="inline-flex justify-center items-center py-3 px-7 xs:py-4 font-semibold text-white rounded-lg bg-[#5D58F5] hover:bg-[#6964FF] transition-[background,box-shadow] text-base sm:text-lg focus-visible:outline-none focus-visible:ring-[#5D58F5]/50 focus-visible:ring-4 text-center whitespace-nowrap"
		:href="!profile ? authLink : '/dashboard'"
	>
		<template v-if="!profile">Start for free</template>
		<template v-else>Dashboard</template>
	</a>
</template>
