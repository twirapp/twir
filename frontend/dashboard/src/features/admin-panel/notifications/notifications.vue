<script setup lang="ts">
import NotificationsForm from './components/notifications-form.vue';
import NotificationsTable from './components/notifications-table.vue';
</script>

<template>
	<div class="w-full flex flex-wrap gap-6">
		<notifications-form />
		<notifications-table />
	</div>
</template>
