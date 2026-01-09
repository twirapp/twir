<script setup lang="ts">
import { EditIcon, MoreVerticalIcon, ToggleLeftIcon, ToggleRightIcon, TrashIcon, UserIcon } from 'lucide-vue-next'


import BadgesPreview from './badges-preview.vue'
import { useBadgesActions } from '../composables/use-badges-actions.js'
import { useBadges } from '../composables/use-badges.js'






const { t } = useI18n()

const { badges } = useBadges()
const badgesActions = useBadgesActions()
</script>

<template>
	<h4 v-if="badges.length" class="scroll-m-20 text-xl font-semibold tracking-tight w-full">
		{{ t('adminPanel.manageBadges.title') }}
	</h4>

	<div class="grid grid-cols-1 gap-4 xl:grid-cols-2 w-full">
		<UiCard
			v-for="badge of badges"
			:key="badge.id"
		>
			<UiCardHeader class="flex-row justify-between items-center border-b p-4">
				<UiCardTitle>
					{{ badge.name }}
				</UiCardTitle>
				<UiDropdownMenu>
					<UiDropdownMenuTrigger as-child>
						<UiButton variant="ghost" size="icon">
							<MoreVerticalIcon class="size-4" />
						</UiButton>
					</UiDropdownMenuTrigger>
					<UiDropdownMenuContent align="end">
						<UiDropdownMenuItem @click="badgesActions.applyUserSearchBadgeFilter(badge)">
							<UserIcon class="mr-2 h-4 w-4" />
							<span>{{ t('adminPanel.manageBadges.users') }}</span>
						</UiDropdownMenuItem>
						<UiDropdownMenuItem @click="badgesActions.editBadge(badge)">
							<EditIcon class="mr-2 h-4 w-4" />
							<span>{{ t('sharedButtons.edit') }}</span>
						</UiDropdownMenuItem>
						<UiDropdownMenuSeparator />
						<UiDropdownMenuItem @click="badgesActions.toggleBadgeEnabled(badge)">
							<template v-if="badge.enabled">
								<ToggleRightIcon class="mr-2 h-4 w-4" />
								<span>{{ t('sharedTexts.enabled') }}</span>
							</template>
							<template v-else>
								<ToggleLeftIcon class="mr-2 h-4 w-4" />
								<span>{{ t('sharedTexts.disabled') }}</span>
							</template>
						</UiDropdownMenuItem>
						<UiDropdownMenuItem @click="badgesActions.showModalDeleteBadge(badge)">
							<TrashIcon class="mr-2 h-4 w-4" />
							<span>{{ t('sharedButtons.delete') }}</span>
						</UiDropdownMenuItem>
					</UiDropdownMenuContent>
				</UiDropdownMenu>
			</UiCardHeader>
			<UiCardContent class="flex flex-col gap-4 p-4">
				<BadgesPreview :image="badge.fileUrl" />
				<div class="flex flex-wrap gap-2">
					<UiBadge :variant="badge.enabled ? 'success' : 'destructive'">
						<template v-if="badge.enabled">
							{{ t('sharedTexts.enabled') }}
						</template>
						<template v-else>
							{{ t('sharedTexts.disabled') }}
						</template>
					</UiBadge>
					<UiBadge>
						{{ t('adminPanel.manageBadges.usesCount', { count: badge.users?.length ?? 0 }) }}
					</UiBadge>
					<UiBadge>
						{{ t('adminPanel.manageBadges.badgeSlot', { slot: badge.ffzSlot }) }}
					</UiBadge>
					<UiBadge>
						<span>
							Created {{ new Date(badge.createdAt).toLocaleString('en-GB', {
								day: '2-digit',
								month: '2-digit',
								year: 'numeric',
								hour: '2-digit',
								minute: '2-digit',
								second: '2-digit'
							}) }}
						</span>
					</UiBadge>
				</div>
			</UiCardContent>
		</UiCard>
	</div>

	<UiActionConfirm
		v-model:open="badgesActions.isShowModalDelete.value"
		@confirm="badgesActions.deleteBadge()"
	/>
</template>
