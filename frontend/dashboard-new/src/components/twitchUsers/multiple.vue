<script setup lang='ts'>
import { refDebounced } from '@vueuse/core';
import { NSelect } from 'naive-ui';
import { computed, ref, watch } from 'vue';
import { defineModel } from 'vue/dist/vue.js';

import { useTwitchSearchUsers } from '@/api/index.js';

const usersIds = defineModel<string[]>({ default: [] });

const userName = ref<string>('');
const userNameDebounced = refDebounced(userName, 1000);
const twitchSearch = useTwitchSearchUsers({
	ids: usersIds,
	names: userNameDebounced,
});

const options = computed(() => {
	if (!twitchSearch.data.value) return [];
	return twitchSearch.data.value.users.map((user) => ({
		label: user.login === user.displayName.toLowerCase()
			? user.displayName
			: `${user.login} (${user.displayName})`,
		value: user.id,
	}));
});

function handleSearch(query: string) {
	userName.value = query;
}
</script>

<template>
  <n-select
    v-model:value="usersIds"
    multiple
    filterable
    placeholder="Search users..."
    :options="options"
    :loading="twitchSearch.isLoading.value"
    clearable
    remote
    :clear-filter-after-select="false"
    @search="handleSearch"
  />
</template>

<style scoped lang='postcss'>

</style>
