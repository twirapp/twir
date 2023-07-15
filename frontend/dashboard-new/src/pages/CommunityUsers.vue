<script setup lang='ts'>
import { useCommunityUsers, UsersOrder, UsersSortBy } from '@/api/index.js';

const communityManager = useCommunityUsers();
const users = communityManager.getAll({
	page: 1,
	limit: 100,
	order: UsersOrder.Desc,
	sortBy: UsersSortBy.Watched,
});
</script>

<template>
  {{ users.data }}
</template>
