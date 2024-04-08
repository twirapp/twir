import { defineStore } from 'pinia';
import { computed } from 'vue';

import { useBadges as useBadgesApi, useAdminBadges } from '@/api/admin/badges';

export const useBadges = defineStore('admin-panel/badges', () => {
	const { data } = useBadgesApi();
	const badgesCrud = useAdminBadges();

	const badgesUpload = badgesCrud.useBadgesUpload();
	const badgesDeleter = badgesCrud.useBadgesDelete();
	const badgesUpdate = badgesCrud.useBadgesUpdate();

	const badgesAdder = badgesCrud.useBadgesUserAdd();
	const badgesRemover = badgesCrud.useBadgesUserRemove();

	const badges = computed(() => data.value?.badges ?? []);

	return {
		badges,
		badgesUpload,
		badgesDeleter,
		badgesUpdate,
		badgesAdder,
		badgesRemover,
	};
});
