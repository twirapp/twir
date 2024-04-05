<script setup lang="ts">
import { SearchIcon } from 'lucide-vue-next';
import { Check, ChevronsUpDown } from 'lucide-vue-next';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUsersTable } from '../composables/use-users-table';

import { Button } from '@/components/ui/button';
import {
  Command,
  CommandGroup,
  CommandItem,
  CommandList,
} from '@/components/ui/command';
import { Input } from '@/components/ui/input';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { cn } from '@/lib/utils';

const { t } = useI18n();
const usersTable = useUsersTable();
const open = ref(false);
</script>

<template>
	<div class="flex gap-2 max-sm:w-full">
		<div class="relative w-full items-center">
			<Input v-model="usersTable.searchInput" type="text" :placeholder="t('sharedTexts.searchPlaceholder')" class="pl-10" />
			<span class="absolute start-2 inset-y-0 flex items-center justify-center px-2">
				<SearchIcon class="size-4 text-muted-foreground" />
			</span>
		</div>
		<Popover v-model:open="open">
			<PopoverTrigger as-child>
				<Button
					variant="outline"
					role="combobox"
					:aria-expanded="open"
					class="w-[200px] justify-between"
				>
					Select filters...
					<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
				</Button>
			</PopoverTrigger>
			<PopoverContent class="w-[200px] p-0">
				<Command multiple>
					<CommandList>
						<CommandGroup>
							<CommandItem
								v-for="filter in usersTable.searchFilters"
								:key="filter.value"
								:value="filter.value"
								@select="(ev) => {
									const value = ev.detail.value;
									if (typeof value === 'string') {
										if (usersTable.selectedFilters[value]) {
											usersTable.selectedFilters[value] = undefined;
										} else {
											usersTable.selectedFilters[value] = true;
										}
									}
								}"
							>
								{{ filter.label }}
								<Check
									:class="cn(
										'ml-auto h-4 w-4',
										usersTable.selectedFilters[filter.value] ? 'opacity-100' : 'opacity-0',
									)"
								/>
							</CommandItem>
						</CommandGroup>
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
	</div>
</template>
