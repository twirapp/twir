<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next';
import { NCard } from 'naive-ui';
import { storeToRefs } from 'pinia';

import BadgesPreview from './badges-preview.vue';
import { useBadges } from '../composables/use-badges';

import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';

const badgesStore = useBadges();
const { badges } = storeToRefs(badgesStore);

async function removeBadge(badgeId: string) {
	await badgesStore.badgesDeleter.mutateAsync(badgeId);
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
				<Button class="max-sm:w-full" size="icon">
					<PencilIcon class="h-4 w-4" />
				</Button>
				<Button class="max-sm:w-full" size="icon" variant="destructive" @click="removeBadge(badge.id)">
					<TrashIcon class="h-4 w-4" />
				</Button>
			</div>
		</div>
	</n-card>
</template>
