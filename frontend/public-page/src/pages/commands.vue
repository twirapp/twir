<script setup lang="ts">
import { IconVideo, IconSword, IconDiamond, IconUserDollar } from '@tabler/icons-vue';
import type { Command } from '@twir/grpc/generated/api/api/commands_unprotected';
import { type FunctionalComponent, computed } from 'vue';

import { useCommands } from '@/api/commands.js';

const props = defineProps<{
	channelId: string
	channelName: string
}>();

const { data: commands } = useCommands(props.channelId);

const commandsGroups = computed<Record<string, Command[]>>(() => {
	if (!commands?.value?.commands.length) return {};

	const cmds = commands.value.commands;

	const groups: Record<string, Command[]> = {
		'empty': cmds.filter(c => c.module === 'CUSTOM' && !c.group),
	};

	// custom commands with groups
	for (const cmd of cmds.filter(c => c.group && c.module === 'CUSTOM')) {
		if (!groups[cmd.group!]) groups[cmd.group!] = [];
		groups[cmd.group!].push(cmd);
	}

	// default commands
	for (const cmd of cmds.filter(c => c.module !== 'CUSTOM')) {
		if (!groups[cmd.module!]) groups[cmd.module!] = [];
		groups[cmd.module!].push(cmd);
	}

	return groups;
});

const permissionsMapping: Record<string, FunctionalComponent> = {
	'BROADCASTER': IconVideo,
	'MODERATOR': IconSword,
	'SUBSCRIBER': IconUserDollar,
	'VIP': IconDiamond,
};
</script>

<template>
	<div class="overflow-auto overflow-y-hidden rounded-lg border-gray-200 shadow-lg">
		<table class="w-full border-collapse text-left text-sm text-slate-200 relative">
			<thead class="bg-neutral-700 text-slate-200">
				<tr>
					<th scope="col" class="px-6 py-4 font-medium">
						Name
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Description/Response
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Permissions
					</th>
					<th scope="col" class="px-6 py-4 font-medium">
						Cooldown
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-neutral-600 border-t border-neutral-600 bg-neutral-700">
				<template v-for="([group, cmds]) in Object.entries(commandsGroups)" :key="group">
					<tr v-if="group !== 'empty'" class="bg-[#24222296] text-center">
						<th colspan="4" class="px-6 py-4">
							{{ group }}
						</th>
					</tr>

					<tr v-for="command of cmds" :key="command.name" class="hover:bg-neutral-600">
						<th class="px-6 py-4">
							<div class="flex flex-wrap">
								{{ [command.name, ...command.aliases].join(", ") }}
							</div>
						</th>
						<th class="px-6 py-4">
							<p v-if="command.description">
								{{ command.description }}
							</p>
							<div v-else class="flex flex-col">
								<p v-for="(response, index) of command.responses" :key="index">
									{{ response }}
								</p>
							</div>
						</th>
						<th class="px-6 py-4">
							<div class="flex flex-wrap gap-2">
								<template v-for="perm of command.permissions" :key="perm.type">
									<div class="group flex relative">
										<span>
											<component
												:is="permissionsMapping[perm.type]"
												v-if="permissionsMapping[perm.type] !== undefined"
											/>
											<template v-else>
												{{ perm.name }}
											</template>
										</span>
										<span
											class="z-10 group-hover:opacity-100 transition-opacity bg-[#1f1f21] px-1 text-sm text-gray-100 rounded-md absolute left-1/2
    -translate-x-1/2 translate-y-full opacity-0 mx-auto"
										>{{ perm.name }}</span>
									</div>
								</template>
							</div>
						</th>
						<th class="px-6 py-4">
							{{ command.cooldown }} ({{ command.cooldownType === 'GLOBAL' ? 'global' : 'per user' }})
						</th>
					</tr>
				</template>
			</tbody>
		</table>
	</div>
</template>
