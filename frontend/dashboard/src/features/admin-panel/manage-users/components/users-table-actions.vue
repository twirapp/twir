<script setup lang="ts">
import { useMutation, useQueryClient } from '@tanstack/vue-query';
import { BanIcon, ShieldIcon } from 'lucide-vue-next';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';

import { useUsersActions } from '../composables/use-users-actions';

import { useProfile } from '@/api';
import { adminApiClient } from '@/api/twirp';
import { Button } from '@/components/ui/button';
import { useBadges } from '@/features/admin-panel/manage-badges/composables/use-badges';

defineProps<{
	userId: string
	isBanned: boolean
	isAdmin: boolean
}>();

const user = useProfile();
const currentUserId = computed(() => user.data.value?.id);
const { switchBan, switchAdmin } = useUsersActions();

const { badges } = storeToRefs(useBadges());

const queryClient = useQueryClient();
const badgerAdder = useMutation({
	mutationFn: async (opts: { badgeId: string, userId: string }) => {
		const req = await adminApiClient.badgeAddUser({
			badgeId: opts.badgeId,
			userId: opts.userId,
		});
		return req.response;
	},
	onSuccess(response, opts) {
			// тут для того чтобы добавить юзера, нужно сделать badges.find(badge => badge.id === opts.badgeId)
			// и запихнуть туда response полностью
	},
});

const badgerRemover = useMutation({
	mutationFn: async (opts: { badgeId: string, userId: string }) => {
		await adminApiClient.badgeDeleteUser({
			badgeId: opts.badgeId,
			userId: opts.userId,
		});
	},
	onSuccess(_, opts) {
		queryClient.setQueriesData(['admin/badges'], (data) => {
			if (!data) return data;
			return {
				badges: data.badges.map(badge => ({
					...badge,
					users: badge.users.filter(u => u.userId != opts.userId),
				})),
			};
		});
	},
});
</script>

<template>
	<div>
		<Button
			v-for="badge of badges"
			:key="badge.id"
			:disabled="badge.users.some(u => u.userId === userId)"
			@click="badgerAdder.mutateAsync({ badgeId: badge.id, userId })"
		>
			add user to {{ badge.name }}
		</Button>
		<Button
			v-for="badge of badges"
			:key="badge.id"
			:disabled="!badge.users.some(u => u.userId === userId)"
			@click="badgerRemover.mutateAsync({ badgeId: badge.id, userId })"
		>
			remove user from {{ badge.name }}
		</Button>
		<div v-if="currentUserId !== userId" class="flex items-center gap-2">
			<Button
				:class="{ 'bg-red-600 hover:bg-red-600/90': isBanned }"
				variant="secondary"
				size="icon"
				@click="switchBan(userId)"
			>
				<BanIcon class="h-4 w-4" />
			</Button>
			<Button
				:class="{ 'bg-green-600 hover:bg-green-600/90': isAdmin }"
				variant="secondary"
				size="icon"
				@click="switchAdmin(userId)"
			>
				<ShieldIcon class="h-4 w-4" />
			</Button>
		</div>
	</div>
</template>
