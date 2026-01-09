<script setup lang="ts">
import { XIcon } from 'lucide-vue-next'

import { useTwitchGetUsers, useTwitchSearchChannels } from '#layers/dashboard/api/twitch'
import { resolveUserName } from '#layers/dashboard/helpers'
import type { AcceptableValue } from 'reka-ui'

interface Props {
	twirOnly?: boolean
	placeholder?: string
}
const props = withDefaults(defineProps<Props>(), {
	twirOnly: false,
	initial: null,
})

const open = ref(false)
const userId = defineModel<string | null>({ required: true })

const selectedUsersQuery = useTwitchGetUsers({ ids: userId })
const selectedUser = computed(() => {
	if (!userId.value) {
		return null
	}

	const user = selectedUsersQuery.data.value?.[0]

	return user
		? {
				label: resolveUserName(user.login, user.displayName),
				value: user.id,
				profileImageUrl: user.profileImageUrl,
			}
		: null
})

const search = ref('')
const searchDebounced = refDebounced(search, 500)

const searchParams = computed(() => ({
	query: searchDebounced.value,
	twirOnly: props.twirOnly,
}))
const twitchSearch = useTwitchSearchChannels(searchParams)
const selectOptions = computed(() => {
	return (
		twitchSearch.data?.value?.map((channel) => ({
			label: resolveUserName(channel.login, channel.displayName),
			value: channel.id,
			profileImageUrl: channel.profileImageUrl,
		})) ?? []
	)
})

function handleSelect(
	event: CustomEvent<{
		originalEvent: PointerEvent
		value?: AcceptableValue
	}>
) {
	if (typeof event.detail.value !== 'string') return
	userId.value = event.detail.value

	search.value = ''
	open.value = false
}
</script>

<template>
	<UiPopover :open="!!selectOptions.length">
		<UiPopoverTrigger as-child>
			<div
				class="flex flex-wrap gap-2 items-center rounded-md border border-input bg-background px-3 py-2 text-sm w-full"
			>
				<div
					v-if="selectedUser"
					class="flex h-6 items-center bg-secondary gap-1 py-1 px-2 text-sm rounded-full"
				>
					<img :src="selectedUser.profileImageUrl" class="size-4 rounded-full" />
					<span>{{ selectedUser.label }}</span>
					<XIcon class="size-4 cursor-pointer" @click="userId = null" />
				</div>
				<input
					v-if="!selectedUser"
					v-model="search"
					type="text"
					:placeholder="placeholder ?? 'Search...'"
					class="text-sm min-h-6 focus:outline-hidden flex-1 bg-transparent px-1"
				/>
			</div>
		</UiPopoverTrigger>
		<UiPopoverContent class="p-0">
			<UiCommand v-model:open="open">
				<UiCommandList>
					<UiCommandGroup>
						<UiCommandItem
							v-for="option in selectOptions"
							:key="option.value"
							:value="option.value"
							@select="handleSelect"
						>
							<div class="flex gap-2 items-center">
								<img :src="option.profileImageUrl" class="size-5 rounded-full" />
								<span>
									{{ option.label }}
								</span>
							</div>
						</UiCommandItem>
					</UiCommandGroup>
				</UiCommandList>
			</UiCommand>
		</UiPopoverContent>
	</UiPopover>
</template>
