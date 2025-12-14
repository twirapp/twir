<script setup lang="ts">
import { refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'

import { useTwitchGetUsers, useTwitchSearchChannels } from '@/api'
import { Command, CommandGroup, CommandItem, CommandList } from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { TagsInput, TagsInputItem, TagsInputItemDelete } from '@/components/ui/tags-input'
import { resolveUserName } from '@/helpers'

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
		value?: string | number | boolean | Record<string, any>
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
	<Popover :open="!!selectOptions.length">
		<PopoverTrigger as-child>
			<TagsInput v-model="userId">
				<TagsInputItem
					v-for="item in selectedUsers"
					:key="item.value"
					:value="item.value"
					class="rounded-full"
				>
					<div class="flex gap-1 items-center py-1 px-2 text-sm rounded bg-transparent">
						<img :src="item.profileImageUrl" class="size-4 rounded-full" />
						<span>{{ item.label }}</span>
					</div>
					<TagsInputItemDelete />
				</TagsInputItem>

				<input
					v-model="search"
					type="text"
					:placeholder="placeholder ?? 'Search...'"
					class="text-sm min-h-6 focus:outline-hidden flex-1 bg-transparent px-1"
				/>
			</TagsInput>
		</PopoverTrigger>
		<PopoverContent class="p-0">
			<Command v-model:open="open">
				<CommandList>
					<CommandGroup>
						<CommandItem
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
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
