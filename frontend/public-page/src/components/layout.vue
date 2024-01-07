<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import Navigation from './navigation.vue';
import Profile from './profile.vue';

const router = useRouter();
const channelName = computed(() => router.currentRoute.value.params?.channelName as string);
const channelId = ref<string | undefined>();
</script>

<template>
	<div class="flex flex-col space-y-[30px] w-[100%]">
		<div class="flex justify-between flex-wrap gap-y-5">
			<profile @update-channel-id="(v) => channelId = v" />
			<navigation />
		</div>

		<router-view v-if="channelId && channelName" v-slot="{ Component }">
			<transition appear mode="out-in">
				<component :is="Component" :channelId="channelId" :channelName="channelName" />
			</transition>
		</router-view>
	</div>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.2s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
