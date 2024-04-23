<script setup lang="ts">
import { useI18n } from 'vue-i18n';

import { useCommunityTableActions } from '../composables/use-community-table-actions.js';

import SearchBar from '@/components/search-bar.vue';
import { Button } from '@/components/ui/button';
import { Command, CommandList, CommandItem } from '@/components/ui/command';
import { Popover, PopoverTrigger, PopoverContent } from '@/components/ui/popover';

const { t } = useI18n();
const communityUsersActions = useCommunityTableActions();

const filters = [
	{
		key: 'watched-time',
		label: 'Watched time',
	},
	{
		key: 'messages',
		label: 'Messages',
	},
	{
		key: 'used-emotes',
		label: 'Used emotes',
	},
	{
		key: 'channel-points-spent',
		label: 'Channel points spent',
	},
];
</script>

<template>
	<div class="flex gap-2 max-sm:flex-col">
		<search-bar v-model="communityUsersActions.searchInput" />
		<Popover>
			<PopoverTrigger as-child>
				<Button variant="outline" size="sm" class="h-9">
					{{ t('community.users.resetActions') }}
				</Button>
			</PopoverTrigger>
			<PopoverContent class="w-[200px] p-0" align="end">
				<Command>
					<CommandList>
						<CommandItem
							v-for="filter of filters"
							:key="filter.key"
							:value="filter.key"
							@select="console.log"
						>
							{{ filter.label }}
						</CommandItem>
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
	</div>
</template>
