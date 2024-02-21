<script setup lang="ts">
import { IconDiamond, IconSword, IconUserDollar, IconVideo } from '@tabler/icons-vue';
import type {
	Command_Permission,
} from '@twir/api/messages/commands_unprotected/commands_unprotected';
import type { FunctionalComponent } from 'vue';

import { Tooltip, TooltipTrigger, TooltipContent } from '@/components/ui/tooltip';

defineProps<{
	permissions: Command_Permission[]
}>();

const permissionsIconsMapping: Record<string, FunctionalComponent> = {
	'BROADCASTER': IconVideo,
	'MODERATOR': IconSword,
	'SUBSCRIBER': IconUserDollar,
	'VIP': IconDiamond,
};
</script>

<template>
	<div class="flex gap-2 items-center justify-center">
		<Tooltip v-for="perm of permissions" :key="perm.name">
			<TooltipTrigger>
				<component :is="permissionsIconsMapping[perm.type]" />
			</TooltipTrigger>
			<TooltipContent>
				{{ perm.name }}
			</TooltipContent>
		</Tooltip>
	</div>
</template>

<style scoped>

</style>
