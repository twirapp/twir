<script setup lang="ts">
import { AwardIcon, CheckIcon } from 'lucide-vue-next'
import { computed } from 'vue'

import { useBadges } from '../../manage-badges/composables/use-badges'

import { Button } from '@/components/ui/button'
import {
	Command,
	CommandGroup,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import { cn } from '@/lib/utils'

const props = defineProps<{ userId: string }>()

const { badges, badgesRemoveUser, badgesAddUser } = useBadges()

const userBadgesIds = computed(() => {
	return badges.value
		.filter((badge) => badge.users?.some((userId) => userId === props.userId))
		.map((badge) => badge.id)
})

async function toggleBadge(badgeId: string) {
	if (userBadgesIds.value.includes(badgeId)) {
		await badgesRemoveUser.executeMutation({ id: badgeId, userId: props.userId })
	} else {
		await badgesAddUser.executeMutation({ id: badgeId, userId: props.userId })
	}
}
</script>

<template>
	<Popover v-if="badges.length">
		<PopoverTrigger as-child>
			<Button variant="secondary" size="sm" class="h-10">
				<AwardIcon class="h-4 w-4" />
			</Button>
		</PopoverTrigger>
		<PopoverContent class="w-[200px] p-0" align="end">
			<Command>
				<CommandList>
					<CommandGroup>
						<CommandItem
							v-for="badge of badges"
							:key="badge.id"
							:value="badge.id"
							@select="toggleBadge(badge.id)"
						>
							<div
								:class="cn(
									'mr-2 flex h-4 w-4 items-center justify-center rounded-sm border border-primary',
									userBadgesIds.includes(badge.id)
										? 'bg-primary text-primary-foreground'
										: 'opacity-50 [&_svg]:invisible',
								)"
							>
								<CheckIcon :class="cn('h-4 w-4')" />
							</div>
							<img :src="badge.fileUrl" class="h-5 w-5 mr-2">
							<span>{{ badge.name }}</span>
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
