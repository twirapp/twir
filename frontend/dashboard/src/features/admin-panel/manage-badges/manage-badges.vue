<script setup lang="ts">
import BadgesForm from './components/badges-form.vue';
import BadgesList from './components/badges-list.vue';
</script>

<template>
	<div class="w-full flex flex-wrap gap-6">
		<badges-form />
		<badges-list />
	</div>
</template>
