<script setup lang="ts">
import { useVirtualizer } from "@tanstack/vue-virtual";
import { debouncedRef } from "@vueuse/core";
import { InfoIcon } from "lucide-vue-next";
import { computed, ref, useTemplateRef } from "vue";
import { useI18n } from "vue-i18n";

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { useGiveaways } from "@/features/giveaways/composables/giveaways-use-giveaways.ts";

const { t } = useI18n();

const { participants, currentGiveaway } = useGiveaways();
const isOnlineChatterGiveaway = computed(() => currentGiveaway.value?.type === "ONLINE_CHATTERS");

const searchTerm = ref("");
const debouncerSearchTerm = debouncedRef(searchTerm, 200);
const filteredParticipants = computed(() => {
	return participants.value.filter((p) =>
		p.displayName.toLowerCase().includes(debouncerSearchTerm.value.toLowerCase()),
	);
});

const participantsRef = useTemplateRef("participantsRef");
const totalParticipants = computed(() => filteredParticipants.value.length);

const rowVirtualizer = useVirtualizer({
	get count() {
		return totalParticipants.value;
	},
	getScrollElement: () => participantsRef.value,
	estimateSize: () => 42,
	overscan: 5,
});
const virtualRows = computed(() => rowVirtualizer.value.getVirtualItems());
const totalSize = computed(() => rowVirtualizer.value.getTotalSize());
</script>

<template>
	<div class="flex-1 flex flex-col h-full min-h-0">
		<div class="p-2 border-b space-y-2">
			<Alert v-if="isOnlineChatterGiveaway" variant="default">
				<InfoIcon class="h-4 w-4" />
				<AlertTitle>{{ t("giveaways.onlineChatterNotice.title") }}</AlertTitle>
				<AlertDescription>
					{{ t("giveaways.onlineChatterNotice.description") }}
				</AlertDescription>
			</Alert>
			<Input v-model="searchTerm" :placeholder="t('sharedTexts.searchPlaceholder')" class="h-10" />
		</div>
		<div ref="participantsRef" class="overflow-y-auto h-full flex-1">
			<div
				:style="{
					height: `${totalSize}px`,
					width: '100%',
					position: 'relative',
				}"
				class="flex-1 overflow-auto"
			>
				<div
					v-for="virtualRow in virtualRows"
					:key="virtualRow.index"
					class="border-b border-border pl-2 flex items-center justify-between pr-2"
					:style="{
						position: 'absolute',
						top: 0,
						left: 0,
						width: '100%',
						height: `${virtualRow.size}px`,
						transform: `translateY(${virtualRow.start}px)`,
					}"
				>
					<span>{{ filteredParticipants[virtualRow.index].displayName }}</span>
					<Badge v-if="filteredParticipants[virtualRow.index].isWinner" variant="success">
						{{ t("giveaways.currentGiveaway.tabs.winners") }}
					</Badge>
				</div>
			</div>
		</div>
	</div>
</template>
