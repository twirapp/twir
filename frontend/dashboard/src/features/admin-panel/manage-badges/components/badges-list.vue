<script setup lang="ts">
import { EditIcon, MoreVerticalIcon, ToggleLeftIcon, ToggleRightIcon, TrashIcon, UserIcon } from 'lucide-vue-next'
import { NCard, NTime } from 'naive-ui'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'

import BadgesPreview from './badges-preview.vue'
import { useBadgesActions } from '../composables/use-badges-actions.js'
import { useBadges } from '../composables/use-badges.js'

import { Badge as UiBadge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import DeleteConfirm from '@/components/ui/delete-confirm.vue'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'

const { t } = useI18n()

const { badges } = storeToRefs(useBadges())
const badgesActions = useBadgesActions()
</script>

<template>
	<h4 v-if="badges.length" class="scroll-m-20 text-xl font-semibold tracking-tight w-full">
		{{ t('adminPanel.manageBadges.title') }}
	</h4>

	<div class="grid grid-cols-1 gap-4 xl:grid-cols-2 w-full">
		<NCard
			v-for="badge of badges"
			:key="badge.id"
			:title="badge.name"
			size="small"
			:segmented="{ content: true }"
			:header-style="{ fontSize: '18px' }"
			bordered
		>
			<template #header-extra>
				<DropdownMenu>
					<DropdownMenuTrigger as-child>
						<Button variant="ghost" size="icon">
							<MoreVerticalIcon class="size-4" />
						</Button>
					</DropdownMenuTrigger>
					<DropdownMenuContent align="end">
						<DropdownMenuItem @click="badgesActions.applyUserSearchBadgeFilter(badge)">
							<UserIcon class="mr-2 h-4 w-4" />
							<span>{{ t('adminPanel.manageBadges.users') }}</span>
						</DropdownMenuItem>
						<DropdownMenuItem @click="badgesActions.editBadge(badge)">
							<EditIcon class="mr-2 h-4 w-4" />
							<span>{{ t('sharedButtons.edit') }}</span>
						</DropdownMenuItem>
						<DropdownMenuSeparator />
						<DropdownMenuItem @click="badgesActions.toggleBadgeEnabled(badge)">
							<template v-if="badge.enabled">
								<ToggleRightIcon class="mr-2 h-4 w-4" />
								<span>{{ t('sharedTexts.enabled') }}</span>
							</template>
							<template v-else>
								<ToggleLeftIcon class="mr-2 h-4 w-4" />
								<span>{{ t('sharedTexts.disabled') }}</span>
							</template>
						</DropdownMenuItem>
						<DropdownMenuItem @click="badgesActions.showModalDeleteBadge(badge)">
							<TrashIcon class="mr-2 h-4 w-4" />
							<span>{{ t('sharedButtons.delete') }}</span>
						</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
			</template>
			<div class="flex flex-col gap-3">
				<BadgesPreview :image="badge.fileUrl" />
				<div class="flex flex-wrap gap-2">
					<UiBadge :variant="badge.enabled ? 'secondary' : 'destructive'">
						<template v-if="badge.enabled">
							{{ t('sharedTexts.enabled') }}
						</template>
						<template v-else>
							{{ t('sharedTexts.disabled') }}
						</template>
					</UiBadge>
					<UiBadge variant="secondary">
						{{ t('adminPanel.manageBadges.usesCount', { count: badge.users?.length ?? 0 }) }}
					</UiBadge>
					<UiBadge variant="secondary">
						{{ t('adminPanel.manageBadges.badgeSlot', { slot: badge.ffzSlot }) }}
					</UiBadge>
					<UiBadge variant="secondary">
						<span>
							Created <NTime :to="new Date(badge.createdAt)" format="dd.MM.yyyy HH:mm:ss" type="datetime" />
						</span>
					</UiBadge>
				</div>
			</div>
		</NCard>
	</div>

	<DeleteConfirm
		v-model:open="badgesActions.isShowModalDelete"
		@confirm="badgesActions.deleteBadge()"
	/>
</template>
