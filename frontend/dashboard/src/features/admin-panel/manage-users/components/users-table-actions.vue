<script setup lang="ts">
import { BanIcon, ShieldIcon } from 'lucide-vue-next';
import { computed } from 'vue';

import UsersBadgeSelector from './users-badge-selector.vue';
import { useUsers } from '../composables/use-users.js';

import { useProfile } from '@/api';
import { Button } from '@/components/ui/button';
import { storeToRefs } from 'pinia';

const props = defineProps<{
	userId: string
	isBanned: boolean
	isBotAdmin: boolean
}>();

const { switchBan, switchAdmin } = useUsers();

const user = storeToRefs(useProfile());
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
			@click="switchBan.executeMutation({ userId })"
		>
			<BanIcon class="h-4 w-4" />
		</Button>
		<Button
			v-if="isVisibleButton"
			:class="{ 'bg-green-600 hover:bg-green-600/90': isBotAdmin }"
			variant="secondary"
			size="icon"
			@click="switchAdmin.executeMutation({ userId })"
		>
			<ShieldIcon class="h-4 w-4" />
		</Button>
	</div>
</template>
