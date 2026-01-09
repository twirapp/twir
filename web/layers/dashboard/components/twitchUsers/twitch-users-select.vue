<script setup lang="ts">
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
const userId = defineModel<string[]>({ required: true, default: [] })

const selectedUsersQuery = useTwitchGetUsers({ ids: userId })
const selectedUsers = computed(() => {
	const users: Record<string, { label: string; value: string; profileImageUrl: string }> = {}
	if (!userId.value.length) {
		return users
	}

	selectedUsersQuery.data.value?.forEach((user) => {
		users[user.id] = {
			label: resolveUserName(user.login, user.displayName),
			value: user.id,
			profileImageUrl: user.profileImageUrl,
		}
	})

	return users
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
		twitchSearch.data?.value?.filter(Boolean).map((channel) => ({
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
	if (userId.value?.includes(event.detail.value)) {
		userId.value = userId.value?.filter((id) => id !== event.detail.value)
	} else {
		userId.value = [...(userId.value ?? []), event.detail.value]
	}

	search.value = ''
	open.value = false
}
</script>

<template>
	<UiPopover :open="!!selectOptions.length">
		<UiPopoverTrigger as-child>
			<UiTagsInput v-model="userId">
				<UiTagsInputItem
					v-for="item in selectedUsers"
					:key="item.value"
					:value="item.value"
					class="rounded-full"
				>
					<div class="flex gap-1 items-center py-1 px-2 text-sm rounded bg-transparent">
						<img :src="item.profileImageUrl" class="size-4 rounded-full" />
						<span>{{ item.label }}</span>
					</div>
					<UITagsInputItemDelete />
				</UiTagsInputItem>

				<input
					v-model="search"
					type="text"
					:placeholder="placeholder ?? 'Search...'"
					class="text-sm min-h-6 focus:outline-hidden flex-1 bg-transparent px-1"
				/>
			</UiTagsInput>
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
