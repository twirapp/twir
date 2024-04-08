<script setup lang="ts">
import { BanIcon, ShieldIcon } from 'lucide-vue-next';
import { computed } from 'vue';

import UsersBadgeSelector from './users-badge-selector.vue';
import { useUsersActions } from '../composables/use-users-actions';

import { useProfile } from '@/api';
import { Button } from '@/components/ui/button';

const props = defineProps<{
	userId: string
	isBanned: boolean
	isAdmin: boolean
}>();

const { switchBan, switchAdmin } = useUsersActions();

const user = useProfile();
const isVisibleButton = computed(() => user.data.value?.id !== props.userId);
</script>

<template>
	<div class="flex items-center gap-2">
		<users-badge-selector :userId="userId" />
		<Button
			v-if="isVisibleButton"
			:class="{ 'bg-red-600 hover:bg-red-600/90': isBanned }"
			variant="secondary"
			size="icon"
			@click="switchBan(userId)"
		>
			<BanIcon class="h-4 w-4" />
		</Button>
		<Button
			v-if="isVisibleButton"
			:class="{ 'bg-green-600 hover:bg-green-600/90': isAdmin }"
			variant="secondary"
			size="icon"
			@click="switchAdmin(userId)"
		>
			<ShieldIcon class="h-4 w-4" />
		</Button>
	</div>
</template>
