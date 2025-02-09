<script setup lang="ts">
import {
	useCommunityUsersSorting,
} from '~/features/community-users/composables/use-community-users-sorting'

const store = useCommunityUsersSorting()
</script>

<template>
	<SearchBar v-model:model-value="store.searchInput" placeholder="Search users..." />
</template>
