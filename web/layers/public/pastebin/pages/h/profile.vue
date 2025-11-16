<script setup lang="ts">
const api = useOapi()

const userStore = useAuth()
await callOnce(UserStoreKey, () => userStore.getUserDataWithoutDashboards())

const { data: pastes, error } = await useAsyncData('userPastesProfile', async () => {
	if (!userStore.userWithoutDashboards) {
		return
	}

	const req = await api.v1.pastebinGetUserList({
		page: 1,
		perPage: 20,
	})
	if (req.error) {
		throw req.error
	}

	return req.data.data
})

const requestUrl = useRequestURL()

async function deletePaste(id: string) {
	await api.v1.pastebinDelete(id)
	await refreshNuxtData('userPastesProfile')
}
</script>

<template>
	<div class="bg-[#1E1E1E] h-full min-h-screen">
		<div v-if="!userStore.userWithoutDashboards">
			<div class="container flex flex-col gap-2 pt-4 items-center justify-center h-full">
				<h1 class="text-2xl">You need to be logged in to view your pastes</h1>

				<button
					v-if="!userStore.userWithoutDashboards"
					class="flex flex-row px-4 py-2 items-center gap-2 bg-[#5D58F5] text-white rounded-lg font-medium focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-[#5D58F5]/50 cursor-pointer hover:bg-[#6964FF] transition-shadow"
					@click="() => userStore.login()"
				>
					Login
					<SvgoSocialTwitch :fontControlled="false" class="w-5 h-5 fill-white" />
				</button>
			</div>
		</div>
		<div v-else class="container flex flex-col gap-2 pt-4 items-center justify-center h-full">
			<div class="border-border w-full">
				<UiTable class="bg-zinc-700 rounded-md">
					<UiTableHeader>
						<UiTableRow>
							<UiTableHead class="w-[5%]">
								<!-- id -->
							</UiTableHead>
							<UiTableHead class="w-[20%]"> Created at </UiTableHead>
							<UiTableHead class="w-[20%]"> Expire at </UiTableHead>
							<UiTableHead class="w-[50%]"> Content </UiTableHead>
							<UiTableHead class="w-[5%]">
								<!-- actions -->
							</UiTableHead>
						</UiTableRow>
					</UiTableHeader>
					<UiTableBody>
						<template v-if="pastes?.items?.length">
							<UiTableRow v-for="paste in pastes?.items" :key="paste.id">
								<UiTableCell>
									<a class="underline" target="_blank" :href="`${requestUrl.origin}/h/${paste.id}`">
										{{ paste.id }}
									</a>
								</UiTableCell>
								<UiTableCell>
									{{ new Date(paste.created_at).toLocaleString() }}
								</UiTableCell>
								<UiTableCell>
									{{ paste.expire_at ? new Date(paste.expire_at).toLocaleString() : '-' }}
								</UiTableCell>
								<UiTableCell>
									{{ paste.content.slice(0, 200) }} {{ paste.content.length > 200 ? '...' : '' }}
								</UiTableCell>
								<UiTableCell>
									<UiButton
										variant="destructive"
										size="sm"
										class="flex items-center gap-2"
										@click="deletePaste(paste.id)"
									>
										<Icon name="lucide:trash" class="size-4" />
									</UiButton>
								</UiTableCell>
							</UiTableRow>
							<UiTableRow v-if="error">
								<UiTableCell :colspan="4">
									{{ error.message }}
								</UiTableCell>
							</UiTableRow>
						</template>
						<UiTableRow v-else>
							<UiTableCell :colspan="4" class="text-center"> No pastes found </UiTableCell>
						</UiTableRow>
					</UiTableBody>
				</UiTable>
			</div>
		</div>
	</div>
</template>
