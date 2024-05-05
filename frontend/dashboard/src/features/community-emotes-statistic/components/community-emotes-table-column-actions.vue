<script setup lang="ts">
import CommunityEmotesDetails
	from '@/features/community-emotes-statistic/components/community-emotes-details.vue'
import {
	useCommunityEmotesDetailsName,
} from '@/features/community-emotes-statistic/composables/use-community-emotes-details'

const props = defineProps<{
	emoteName: string
}>()

const nameStore = useCommunityEmotesDetailsName()
function setDetailsName() {
	nameStore.setEmoteName(props.emoteName)
}
</script>

<template>
	<div class="flex gap-2">
		<CommunityEmotesDetails @select="setDetailsName" />
	</div>
</template>
