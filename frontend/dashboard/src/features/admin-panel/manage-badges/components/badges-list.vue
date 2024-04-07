<script setup lang="ts">
import type { Badge } from '@twir/api/messages/badges_unprotected/badges_unprotected';
import { PencilIcon, TrashIcon } from 'lucide-vue-next';
import { NCard } from 'naive-ui';
import { storeToRefs } from 'pinia';

import BadgesPreview from './badges-preview.vue';
import { useBadges } from '../composables/use-badges';
import { useBadgesForm } from '../composables/use-badges-form';

import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';

const badgesForm = useBadgesForm();
const badgesStore = useBadges();
const { badges } = storeToRefs(badgesStore);

async function removeBadge(badgeId: string) {
	await badgesStore.badgesDeleter.mutateAsync(badgeId);
}

function editBadge(badge: Badge) {
	badgesForm.editableBadgeId = badge.id;
	badgesForm.form.setFieldValue('name', badge.name);
	badgesForm.form.setFieldValue('image', badge.fileUrl);
}
</script>

<template>
	<n-card v-for="badge of badges" :key="badge.id" size="small" bordered>
		<Label>
			{{ badge.name }}
		</Label>
		<div class="flex justify-between gap-2 max-sm:flex-col">
			<badges-preview :image="badge.fileUrl" />
			<div class="flex items-end gap-2">
				<Button class="max-sm:w-full" size="icon" @click="editBadge(badge)">
					<PencilIcon class="h-4 w-4" />
				</Button>
				<Button class="max-sm:w-full" size="icon" variant="destructive" @click="removeBadge(badge.id)">
					<TrashIcon class="h-4 w-4" />
				</Button>
			</div>
		</div>
	</n-card>
</template>
