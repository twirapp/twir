<script setup lang="ts">
import type { HTMLAttributes } from 'vue'

import { Badge } from '@/components/ui/badge'
import type { Platform } from '~/gql/graphql.js'

import { getPlatformMeta } from '~/utils/platforms.js'

const props = defineProps<{
	platform: Platform
	class?: HTMLAttributes['class']
}>()

const meta = computed(() => getPlatformMeta(props.platform))
</script>

<template>
	<Badge
		variant="outline"
		class="w-fit text-[10px] uppercase tracking-wide"
		:class="[meta.badgeClass, props.class]"
	>
		{{ meta.label }}
	</Badge>
</template>
