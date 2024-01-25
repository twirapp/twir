<script lang="ts" setup>
import { storeToRefs } from 'pinia';

import TableRowsSkeleton from '@/components/TableRowsSkeleton.vue';
import {
	Table,
	TableBody,
	TableHead,
	TableHeader,
	TableRow,
} from '@/components/ui/table';
import { useCommands } from '@/composables/use-commands';
import { useStreamerProfile } from '@/composables/use-streamer-profile';
import TableRowsCommands from '@/pages/commands/TableRowsCommands.vue';

const { profile } = storeToRefs(useStreamerProfile());

const { data: commands, isLoading: isCommandsLoading } = useCommands(profile.value?.id);
</script>

<template>
	<div class="rounded-md border">
		<Table>
			<TableHeader>
				<TableRow>
					<TableHead class="w-[150px] max-w-[150px] min-w-[150px]">
						Names
					</TableHead>
					<TableHead class="w-full">
						Description
					</TableHead>
					<TableHead class="w-[100px]">
						Permissions
					</TableHead>
					<TableHead class="text-right w-[100px]">
						Cooldown
					</TableHead>
				</TableRow>
			</TableHeader>
			<Transition name="table-rows" appear mode="out-in">
				<TableBody v-if="isCommandsLoading || !commands?.commands">
					<table-rows-skeleton :rows="20" :colspan="4" />
				</TableBody>
				<TableBody v-else>
					<table-rows-commands :commands="commands.commands" />
				</TableBody>
			</Transition>
		</Table>
	</div>
</template>

<style scoped>
.table-rows-enter-active,
.table-rows-leave-active {
	transition: opacity 0.5s ease;
}

.table-rows-enter-from,
.table-rows-leave-to {
	opacity: 0;
}
</style>
