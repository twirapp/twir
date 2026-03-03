<script setup lang="ts">
import PastebinPasteCard from "../../components/pastebin/paste-card.vue";

definePageMeta({
	layout: "landing",
});

const api = useOapi();
const userStore = useAuth();
await callOnce(UserStoreKey, () => userStore.getUserDataWithoutDashboards());

const currentPage = ref(1);
const perPage = ref(10);

const { data: pastesData, refresh } = await useAsyncData(
	"userPastesProfile",
	async () => {
		if (!userStore.userWithoutDashboards) {
			return null;
		}

		const req = await api.v1.pastebinGetUserList({
			page: currentPage.value,
			perPage: perPage.value,
		});
		if (req.error) {
			throw req.error;
		}

		return req.data;
	},
	{
		watch: [currentPage, perPage],
	},
);

const pastes = computed(() => pastesData.value?.data?.items ?? []);
const total = computed(() => pastesData.value?.data?.total ?? 0);
const totalPages = computed(() => Math.ceil(total.value / perPage.value));

function handlePasteDeleted() {
	refresh();
}
</script>

<template>
	<div class="h-full w-full">
		<div
			class="absolute inset-0 bg-[linear-gradient(to_right,hsl(240,11%,9%)_1px,transparent_1px),linear-gradient(to_bottom,hsl(240,11%,9%)_1px,transparent_1px)] bg-size-[36px_36px] mask-[linear-gradient(to_bottom,transparent_15%,black_100%)]"
		></div>
		<div class="container mx-auto py-8 relative min-h-[calc(100vh-73px)]">
			<!-- Not logged in state -->
			<div
				v-if="!userStore.userWithoutDashboards"
				class="max-w-6xl mx-auto flex flex-col items-center justify-center gap-6 py-20"
			>
				<div class="text-center space-y-2">
					<h1 class="text-3xl font-bold">Authentication Required</h1>
					<p class="text-[hsl(240,11%,65%)]">You need to be logged in to view your pastes</p>
				</div>
				<button
					class="flex flex-row px-6 py-3 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 cursor-pointer hover:bg-[#6964FF] transition-all"
					@click="() => userStore.login()"
				>
					Login with Twitch
					<SvgoSocialTwitch :fontControlled="false" class="w-5 h-5 fill-white" />
				</button>
			</div>

			<!-- Logged in state -->
			<div v-else class="max-w-6xl mx-auto space-y-8">
				<div class="space-y-4">
					<NuxtLink
						to="/h"
						class="inline-flex items-center gap-2 px-4 py-2 rounded-lg border border-[hsl(240,11%,30%)] hover:border-[hsl(240,11%,45%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,25%)] transition-colors text-sm font-medium"
					>
						<Icon name="lucide:arrow-left" class="w-4 h-4" />
						Back to Hastebin
					</NuxtLink>

					<div class="flex items-center justify-between flex-wrap gap-4">
						<div>
							<h1 class="text-3xl font-bold">My Pastes</h1>
							<p class="text-[hsl(240,11%,65%)]">Manage your pastes and view their content</p>
						</div>
					</div>
				</div>

				<!-- Empty state -->
				<div v-if="pastes.length === 0" class="text-center py-12">
					<p class="text-[hsl(240,11%,65%)]">You don't have any pastes yet</p>
					<NuxtLink
						to="/h"
						class="text-[hsl(240,11%,85%)] hover:text-white transition-colors underline"
					>
						Create your first paste
					</NuxtLink>
				</div>

				<!-- Pastes list -->
				<div v-else class="space-y-6">
					<!-- Pastes Grid -->
					<div class="grid gap-4">
						<PastebinPasteCard
							v-for="paste in pastes"
							:key="paste.id"
							:paste="paste"
							@deleted="handlePasteDeleted"
						/>
					</div>

					<!-- Pagination -->
					<div v-if="totalPages > 1" class="flex items-center justify-between">
						<div class="text-sm text-[hsl(240,11%,65%)]">
							Showing {{ (currentPage - 1) * perPage + 1 }} to
							{{ Math.min(currentPage * perPage, total) }} of {{ total }} pastes
						</div>
						<div class="flex items-center gap-2">
							<button
								:disabled="currentPage === 1"
								@click="currentPage--"
								class="px-3 py-2 rounded-lg border border-[hsl(240,11%,30%)] hover:border-[hsl(240,11%,45%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,25%)] transition-colors text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
							>
								<Icon name="lucide:chevron-left" class="w-4 h-4" />
							</button>
							<span class="text-sm text-[hsl(240,11%,65%)]">
								Page {{ currentPage }} of {{ totalPages }}
							</span>
							<button
								:disabled="currentPage === totalPages"
								@click="currentPage++"
								class="px-3 py-2 rounded-lg border border-[hsl(240,11%,30%)] hover:border-[hsl(240,11%,45%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,25%)] transition-colors text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
							>
								<Icon name="lucide:chevron-right" class="w-4 h-4" />
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
