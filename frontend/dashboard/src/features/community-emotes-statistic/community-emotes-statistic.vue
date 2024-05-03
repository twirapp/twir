<script setup lang="ts">
import { storeToRefs } from 'pinia'

import {
	useCommunityEmotesStatisticTable
} from '@/features/community-emotes-statistic/composables/use-community-emotes-statistic-table'

const { stats } = storeToRefs(useCommunityEmotesStatisticTable())
</script>

<template>
	<div>
		<h1>stats</h1>
		{{ stats }}
	</div>
</template>
