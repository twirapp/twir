<script setup lang="ts">
import { CheckIcon } from 'lucide-vue-next';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUsersTable } from '../composables/use-users-table';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
	Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator,
} from '@/components/ui/command';
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover';
import { cn } from '@/lib/utils';

const { t } = useI18n();
const usersTable = useUsersTable();

const { selectFilters } = storeToRefs(usersTable);
const selectedValuesCount = computed(() => Object.values(selectFilters.value).flat().filter(Boolean).length);

function clearFilters() {
	Object.keys(selectFilters.value).forEach((key) => {
		selectFilters.value[key] = undefined;
	});
}

function setFilterValue(key: string) {
	if (selectFilters.value[key]) {
		selectFilters.value[key] = undefined;
		return;
	}
	selectFilters.value[key] = true;
}
</script>

<template>
	<Popover>
		<PopoverTrigger as-child>
			<Button variant="outline" size="sm" class="h-8 border-dashed">
				<PlusCircledIcon class="mr-2 h-4 w-4" />
				filters

				<template v-if="selectedValuesCount">
					<Separator orientation="vertical" class="mx-2 h-4" />
					<Badge
						variant="secondary"
						class="rounded-sm px-1 font-normal lg:hidden"
					>
						{{ selectedValuesCount }}
					</Badge>
					<div class="hidden space-x-1 lg:flex">
						<Badge
							v-if="selectedValuesCount"
							variant="secondary"
							class="rounded-sm px-1 font-normal"
						>
							{{ selectedValuesCount }} selected
						</Badge>

						<template v-else>
							<Badge
								v-for="filterKey in Object.keys(selectFilters)"
								:key="filterKey"
								variant="secondary"
								class="rounded-sm px-1 font-normal"
							>
								{{ filterKey }}
							</Badge>
						</template>
					</div>
				</template>
			</Button>
		</PopoverTrigger>
		<PopoverContent class="w-[200px] p-0" align="start">
			<Command>
				<CommandInput placeholder="qweqweqweqwe" />
				<CommandList>
					<CommandEmpty>No results found.</CommandEmpty>
					<CommandGroup>
						<CommandItem
							v-for="filterKey in Object.keys(selectFilters)"
							:key="filterKey"
							:value="filterKey"
							@select="setFilterValue(filterKey)"
						>
							<div
								:class="cn(
									'mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary',
									selectFilters[filterKey]
										? 'bg-primary text-primary-foreground'
										: 'opacity-50 [&_svg]:invisible',
								)"
							>
								<CheckIcon :class="cn('h-4 w-4')" />
							</div>
							<span>{{ filterKey }}</span>
						</CommandItem>
					</CommandGroup>

					<template v-if="selectedValuesCount">
						<CommandSeparator />
						<CommandGroup>
							<CommandItem
								:value="{ label: 'Clear filters' }"
								class="justify-center text-center"
								@select="clearFilters"
							>
								Clear filters
							</CommandItem>
						</CommandGroup>
					</template>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
