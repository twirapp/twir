<script setup lang="ts">
import { useProfile, useUser } from '@/api';

const { data: profile } = useProfile();

const manager = useUser();
const regenerationKeyApi = manager.useRegenerateApiKey();
const updateApi = manager.useUpdate();
</script>

<template>
	<div>
		<h1>Profile</h1>
		{{ profile }}
	</div>
</template>

<style scoped>

</style>
