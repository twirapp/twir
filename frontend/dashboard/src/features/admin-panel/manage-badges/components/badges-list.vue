<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next';
import { NCard } from 'naive-ui';
import { storeToRefs } from 'pinia';

import BadgesPreview from './badges-preview.vue';
import { useBadges } from '../composables/use-badges';

import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';

const { badges } = storeToRefs(useBadges());
</script>

<template>
	<n-card v-for="badge of badges" :key="badge.name" size="small" bordered>
		<Label>
			{{ badge.name }}
		</Label>
		<div class="flex justify-between gap-2 max-sm:flex-col">
			<badges-preview :image="badge.image" />
			<div class="flex items-end gap-2">
				<Button class="max-sm:w-full" size="icon">
					<PencilIcon class="h-4 w-4" />
				</Button>
				<Button class="max-sm:w-full" size="icon" variant="destructive">
					<TrashIcon class="h-4 w-4" />
				</Button>
			</div>
		</div>
	</n-card>
</template>
