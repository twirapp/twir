<script setup lang="ts">
import {
	IconDiamond,
	IconSword,
	IconUser,
	IconUserDollar,
	IconVideo,
	IconWorld,
} from '@tabler/icons-vue';
import type { Command } from '@twir/api/messages/commands/commands';
import { computed, type FunctionalComponent } from 'vue';

import { TableCell, TableRow } from '@/components/ui/table';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';

const props = defineProps<{
	commands: Command[]
}>();

const mappedCommands = computed(() => {
	const cmds = props.commands.map(c => ({
		...c,
		cooldown: c.cooldown.toString(),
	}));

	if (!cmds) {
		return [];
	}

	return cmds.sort((a, b) => a.name.localeCompare(b.name));
});

const permissionsIconsMapping: Record<string, FunctionalComponent> = {
	'BROADCASTER': IconVideo,
	'MODERATOR': IconSword,
	'SUBSCRIBER': IconUserDollar,
	'VIP': IconDiamond,
};
</script>

<template>
	<TableRow
		v-for="command of mappedCommands"
		:key="command.name"
	>
		<TableCell class="font-medium">
			<div class="flex flex-wrap gap-x-2">
				<span>!{{ command.name }}</span>
				<span v-for="(aliase, aliaseIndex) of command.aliases" :key="aliaseIndex">
					!{{ aliase }}
				</span>
			</div>
		</TableCell>
		<TableCell>
			<span v-if="command.description">
				{{ command.description }}
			</span>
			<div v-else class="flex flex-col gap-0.5">
				<span v-for="(response, responseIndex) of command.responses" :key="responseIndex">
					{{ response }}
				</span>
			</div>
		</TableCell>
		<TableCell>
			<Tooltip
				v-for="(permission, permissionIndex) of command.permissions"
				:key="permissionIndex"
			>
				<TooltipTrigger>
					<component :is="permissionsIconsMapping[permission.type]" />
				</TooltipTrigger>
				<TooltipContent>
					{{ permission.name }}
				</TooltipContent>
			</Tooltip>
		</TableCell>
		<TableCell class="text-right">
			<Tooltip>
				<TooltipTrigger>
					<div class="flex items-center justify-end gap-0.5">
						<span>{{ command.cooldown }}</span>
						<IconWorld v-if="command.cooldownType === 'GLOBAL'" :class="$style.cooldownIcon" />
						<IconUser v-else :class="$style.cooldownIcon" />
					</div>
				</TooltipTrigger>
				<TooltipContent>
					{{ command.cooldownType }}
				</TooltipContent>
			</Tooltip>
		</TableCell>
	</TableRow>
</template>

<style module>
.cooldownIcon {
	height: 18px;
	width: 18px;
}
</style>
