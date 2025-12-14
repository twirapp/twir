<script setup lang="ts">
import { useCommunityEmotesDetailsName } from '@/features/community-emotes-statistic/composables/use-community-emotes-details'
import CommunityEmotesDetails from '@/features/community-emotes-statistic/ui/community-emotes-details.vue'

const props = defineProps<{ emoteName: string }>()

const nameStore = useCommunityEmotesDetailsName()
function setDetailsName() {
	nameStore.setEmoteName(props.emoteName)
}
</script>

<template>
	<div class="flex flex-col gap-0.5 wrap-break-word max-w-[450px]">
		<h1 class="text-xl">
			{{ emoteName }}
		</h1>
		<CommunityEmotesDetails @select="setDetailsName" />
	</div>
</template>
