<script setup lang="ts">
import { EditIcon, MoreVerticalIcon, ToggleLeftIcon, ToggleRightIcon, TrashIcon, UserIcon } from 'lucide-vue-next';
import { NCard, NTime } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import BadgesPreview from './badges-preview.vue';
import { useUsersTableFilters } from '../../manage-users/composables/use-users-table-filters.js';
import { useBadgesForm } from '../composables/use-badges-form.js';
import { useBadges } from '../composables/use-badges.js';

import { Badge as UiBadge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import DeleteConfirm from '@/components/ui/delete-confirm.vue';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
	DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu';
import { useLayout } from '@/composables/use-layout.js';
import type { Badge } from '@/gql/graphql';

const { t } = useI18n();
const layout = useLayout();
const badgesForm = useBadgesForm();
const badgesStore = useBadges();
const { badges } = storeToRefs(badgesStore);

async function removeBadge(badgeId: string) {
	await badgesStore.badgesDelete.executeMutation({ id: badgeId });
	deleteBadgeId.value = null;
}

function editBadge(badge: Badge) {
	badgesForm.editableBadgeId = badge.id;
	badgesForm.nameField.fieldModel = badge.name;
	badgesForm.fileField.fieldModel = badge.fileUrl;
	badgesForm.slotField.fieldModel = badge.ffzSlot;
	layout.scrollToTop();
}

function toggleEnableBadge(badge: Badge) {
	badgesStore.badgesUpdate.executeMutation({
		id: badge.id,
		opts: { enabled: !badge.enabled },
	});
}

const showDelete = ref(false);
const deleteBadgeId = ref<string | null>(null);
function deleteBadge(badge: Badge): void {
	showDelete.value = true;
	deleteBadgeId.value = badge.id;
}

const router = useRouter();
const userFilters = useUsersTableFilters();

function applyUserSearchBadgeFilter(badge: Badge): void {
	userFilters.clearFilters();
	userFilters.selectedBadges.push(badge.id);
	router.push({ query: { tab: 'users' } });
}
</script>

<template>
	<h4 v-if="badges.length" class="scroll-m-20 text-xl font-semibold tracking-tight w-full">
		{{ t('adminPanel.manageBadges.title') }}
	</h4>

	<div class="grid grid-cols-1 gap-4 xl:grid-cols-2 w-full">
		<n-card
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
						<DropdownMenuItem @click="applyUserSearchBadgeFilter(badge)">
							<UserIcon class="mr-2 h-4 w-4" />
							<span>{{ t('adminPanel.manageBadges.users') }}</span>
						</DropdownMenuItem>
						<DropdownMenuItem @click="editBadge(badge)">
							<EditIcon class="mr-2 h-4 w-4" />
							<span>{{ t('sharedButtons.edit') }}</span>
						</DropdownMenuItem>
						<DropdownMenuSeparator />
						<DropdownMenuItem @click="toggleEnableBadge(badge)">
							<template v-if="badge.enabled">
								<ToggleRightIcon class="mr-2 h-4 w-4" />
								<span>{{ t('sharedTexts.enabled') }}</span>
							</template>
							<template v-else>
								<ToggleLeftIcon class="mr-2 h-4 w-4" />
								<span>{{ t('sharedTexts.disabled') }}</span>
							</template>
						</DropdownMenuItem>
						<DropdownMenuItem @click="deleteBadge(badge)">
							<TrashIcon class="mr-2 h-4 w-4" />
							<span>{{ t('sharedButtons.delete') }}</span>
						</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
			</template>
			<div class="flex flex-col gap-3">
				<badges-preview :image="badge.fileUrl" />
				<div class="flex flex-wrap gap-2">
					<ui-badge :variant="badge.enabled ? 'secondary' : 'destructive'">
						<template v-if="badge.enabled">
							{{ t('sharedTexts.enabled') }}
						</template>
						<template v-else>
							{{ t('sharedTexts.disabled') }}
						</template>
					</ui-badge>
					<ui-badge variant="secondary">
						{{ t('adminPanel.manageBadges.usesCount', { count: badge.users?.length ?? 0 }) }}
					</ui-badge>
					<ui-badge variant="secondary">
						{{ t('adminPanel.manageBadges.badgeSlot', { slot: badge.ffzSlot }) }}
					</ui-badge>
					<ui-badge variant="secondary">
						<span>
							Created <n-time :to="new Date(badge.createdAt)" format="dd.MM.yyyy HH:mm:ss" type="datetime" />
						</span>
					</ui-badge>
				</div>
			</div>
		</n-card>
	</div>

	<delete-confirm
		v-model:open="showDelete"
		@confirm="removeBadge(deleteBadgeId!)"
	/>
</template>
