<script setup lang="ts">
import NotificationsForm from './components/notifications-form.vue'
import NotificationsList from './components/notifications-list.vue'
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<NotificationsForm />
		<NotificationsList />
	</div>
</template>
