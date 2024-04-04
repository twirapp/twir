<script setup lang="ts">
import NotificationForm from './notification-form.vue';
import NotificationsList from './notifications-list.vue';
</script>

<template>
	<div class="w-full flex flex-wrap gap-6">
		<notification-form />
		<notifications-list />
	</div>
</template>
