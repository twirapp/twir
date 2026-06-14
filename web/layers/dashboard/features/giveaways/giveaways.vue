<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";

import type { PageLayoutTab } from "@/layout/page-layout.vue";

import GiveawaysPagesHistory from "@/features/giveaways/pages/giveaways-pages-history.vue";
import GiveawaysPagesList from "@/features/giveaways/pages/giveaways-pages-list.vue";
import GiveawaysPagesView from "@/features/giveaways/pages/giveaways-pages-view.vue";
import GiveawaysSettingsDialog from "@/features/giveaways/ui/giveaways-settings-dialog.vue";
import PageLayout from "@/layout/page-layout.vue";

const { t } = useI18n();
const route = useRoute();

// Check if we're on the view page
const isViewPage = route.name === "giveaways-view";

// Define tabs for the main giveaways page
const tabs: PageLayoutTab[] = [
	{
		name: "list",
		component: GiveawaysPagesList,
		title: t("giveaways.pages.list.title"),
	},
	{
		name: "history",
		component: GiveawaysPagesHistory,
		title: t("giveaways.pages.history.title"),
	},
];
</script>

<template>
	<!-- Show the giveaway view page if we're on that route -->
	<div v-if="isViewPage">
		<GiveawaysPagesView />
	</div>

	<!-- Otherwise show the main giveaways page with tabs -->
	<PageLayout v-else :tabs="tabs" active-tab="list">
		<template #action>
			<GiveawaysSettingsDialog />
		</template>

		<template #title>
			{{ t("giveaways.title") }}
		</template>
	</PageLayout>
</template>
