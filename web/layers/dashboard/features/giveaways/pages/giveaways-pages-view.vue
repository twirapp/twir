<script setup lang="ts">
import { watch } from "vue";
import { useRoute } from "vue-router";

import { useGiveaways } from "@/features/giveaways/composables/giveaways-use-giveaways.js";
import GiveawaysCurrentGiveaway from "@/features/giveaways/ui/giveaways-current-giveaway.vue";

const route = useRoute<'dashboard-giveaways-view-id'>();
const { loadParticipants } = useGiveaways();

// Следим за изменением ID гива в роуте и перезагружаем данные
watch(
	() => route.params.id,
	(giveawayId) => {
		if (giveawayId && typeof giveawayId === "string") {
			loadParticipants(giveawayId);
		}
	},
	{ immediate: true },
);
</script>

<template>
	<GiveawaysCurrentGiveaway />
</template>
