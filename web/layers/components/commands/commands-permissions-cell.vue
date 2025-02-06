<script lang="ts" setup>
import type { Command } from '~/lib/commands'

import { Icon } from '#components'

defineProps<{
	permissions: Command['permissions']
}>()

const permissionsIconsMapping: Record<string, VNode> = {
	BROADCASTER: h(Icon, { name: 'lucide:video' }),
	MODERATOR: h(Icon, { name: 'lucide:sword' }),
	SUBSCRIBER: h(Icon, { name: 'lucide:badge-dollar-sign' }),
	VIP: h(Icon, { name: 'lucide:gem' }),
}
</script>

<template>
	<div class="flex gap-2 items-center justify-center">
		<UiTooltip v-for="perm of permissions" :key="perm.name">
			<UiTooltipTrigger>
				<component :is="permissionsIconsMapping[perm.type]" />
			</UiTooltipTrigger>
			<UiTooltipContent>
				{{ perm.name }}
			</UiTooltipContent>
		</UiTooltip>
	</div>
</template>
