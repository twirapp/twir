<script setup lang="ts">
import type { Badge } from '@twir/api/messages/badges_unprotected/badges_unprotected';
import { NCard } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { ref, unref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import BadgesPreview from './badges-preview.vue';
import { useUsersTableFilters } from '../../manage-users/composables/use-users-table-filters';
import { useBadges } from '../composables/use-badges';
import { useBadgesForm } from '../composables/use-badges-form';

import { Button } from '@/components/ui/button';
import DeleteConfirm from '@/components/ui/delete-confirm.vue';
import { useLayout } from '@/composables/use-layout';

const { t } = useI18n();
const layout = useLayout();
const badgesForm = useBadgesForm();
const badgesStore = useBadges();
const { badges } = storeToRefs(badgesStore);

async function removeBadge(badgeId: string) {
	await badgesStore.badgesDeleter.mutateAsync(badgeId);
	deleteBadgeId.value = null;
}

function editBadge(badge: Badge) {
	badgesForm.editableBadgeId = badge.id;
	badgesForm.form.setFieldValue('name', badge.name);
	badgesForm.form.setFieldValue('image', badge.fileUrl);
	layout.scrollToTop();
}

const showDelete = ref(false);
const deleteBadgeId = ref<string | null>(null);
function deleteBadge(badgeId: string): void {
	showDelete.value = true;
	deleteBadgeId.value = badgeId;
}

const router = useRouter();
const userFilters = useUsersTableFilters();

function applyUserSearchBadgeFilter(badge: Badge): void {
	userFilters.selectedBadges.push(badge.id);
	router.push({ query: { tab: 'users' } });
}

</script>

<template>
	<h4 v-if="badges.length" class="scroll-m-20 text-xl font-semibold tracking-tight w-full">
		{{ t('adminPanel.manageBadges.title') }}
	</h4>

	<div class="grid grid-cols-1 gap-4 xl:grid-cols-2 w-full">
		<n-card v-for="badge of badges" :key="badge.id" size="small" bordered>
			<badges-preview class="mt-2" :image="badge.fileUrl" />
			<div class="flex justify-between items-center gap-4 mt-4 max-sm:flex-col max-sm:items-start">
				<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
					{{ badge.name }}
				</h4>
				<div class="flex items-end gap-2">
					<Button
						class="max-sm:grow flex items-center space-x-4"
						variant="secondary"
						size="sm"
						@click="applyUserSearchBadgeFilter(unref(badge))"
					>
						{{ t('adminPanel.manageBadges.usersCount', { count: badge.users.length }) }}
					</Button>
					<Button
						class="max-sm:grow"
						variant="secondary"
						size="sm"
						@click="editBadge(unref(badge))"
					>
						{{ t('sharedButtons.edit') }}
					</Button>
					<Button
						class="max-sm:grow"
						variant="destructive"
						size="sm"
						@click="deleteBadge(badge.id)"
					>
						{{ t('sharedButtons.delete') }}
					</Button>
				</div>
			</div>
		</n-card>
	</div>

	<delete-confirm
		v-model:open="showDelete"
		@confirm="removeBadge(deleteBadgeId!)"
	/>
</template>
