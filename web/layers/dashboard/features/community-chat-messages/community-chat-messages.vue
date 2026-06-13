<script lang="ts" setup>
import ChatMessages from './ui/chat-messages.vue'
import Filters from './ui/filters.vue'
</script>

<template>
	<div class="flex flex-col gap-4">
		<Filters />
		<ChatMessages />
	</div>
</template>
