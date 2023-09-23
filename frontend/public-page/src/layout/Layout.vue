<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import Menu from './Menu.vue';
import Profile from './Profile.vue';

const router = useRouter();
const channelName = computed(() => router.currentRoute.value.params?.channelName as string);
const channelId = ref<string | undefined>();
</script>

<template>
	<div class="flex flex-col space-y-[30px] w-[100%]">
		<div class="flex justify-between flex-wrap gap-y-5">
			<div>
				<Profile @updateChannelId="(v) => channelId = v" />
			</div>
			<div>
				<Menu />
			</div>
		</div>

		<Transition>
			<router-view v-if="channelId && channelName" v-slot="{ Component }">
				<component :is="Component" :channelId="channelId" :channelName="channelName" />
			</router-view>
		</Transition>
	</div>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
