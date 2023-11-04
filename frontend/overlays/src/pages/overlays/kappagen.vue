<script setup lang="ts">
import { useRoute } from 'vue-router';
import '@twir/frontend-chat/style.css';

const route = useRoute();
const apiKey = route.params.apiKey as string;
</script>

<template>
	{{ apiKey }}
</template>
