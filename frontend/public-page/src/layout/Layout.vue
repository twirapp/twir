<script setup lang="ts">
import Menu from './Menu.vue';
import Profile from './Profile.vue';
</script>

<template>
	<div class="flex flex-col space-y-[30px] w-[100%]">
		<div class="flex justify-between flex-wrap gap-y-5">
			<div>
				<Profile />
			</div>
			<div>
				<Menu />
			</div>
		</div>

		<Transition>
			<router-view v-slot="{ Component }">
				<component :is="Component" />
			</router-view>
		</Transition>
	</div>
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
