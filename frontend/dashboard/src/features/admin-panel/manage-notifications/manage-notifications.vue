<script setup lang="ts">
import NotificationsForm from './components/notifications-form.vue';
import NotificationsList from './components/notifications-list.vue';
</script>

<template>
	<div class="w-full flex flex-wrap gap-6">
		<notifications-form />
		<notifications-list />
	</div>
</template>
