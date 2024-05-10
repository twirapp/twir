<script setup lang="ts">
import { refDebounced } from '@vueuse/core'
import { Check, ChevronsUpDown } from 'lucide-vue-next'
import { computed, ref } from 'vue'

import type { Channel, TwitchSearchChannelsRequest,TwitchUser } from '@twir/api/messages/twitch/twitch'
import type { SelectEvent } from 'radix-vue/dist/Listbox/ListboxItem'
import type { AcceptableValue } from 'radix-vue/dist/shared/types'

import { useTwitchGetUsers,useTwitchSearchChannels } from '@/api'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import { resolveUserName } from '@/helpers'
import { cn } from '@/lib/utils'

interface Props {
	initial?: string | string[] | null
	twirOnly?: boolean
	multiple?: boolean
}
const props = withDefaults(defineProps<Props>(), {
	twirOnly: false,
	initial: null,
	multiple: false,
})

const userId = defineModel<string | string[] | undefined | null>()

const ids = computed<string[]>(() => {
	const selectedIds = (Array.isArray(userId.value) ? userId.value : [userId.value]).filter(i => !!i) as string[]
	const initialArray = (Array.isArray(props.initial) ? props.initial : [props.initial]).filter(i => !!i) as string[]

	return [...initialArray, ...selectedIds]
})

const getUsers = useTwitchGetUsers({
	ids,
})

const userName = ref('')
const userNameDebounced = refDebounced(userName, 500)

const searchParams = computed<TwitchSearchChannelsRequest>(() => ({
	query: userNameDebounced.value,
	twirOnly: props.twirOnly,
}))
const twitchSearch = useTwitchSearchChannels(searchParams)

function mapOptions(users: (TwitchUser | Channel)[]) {
	return users.map((user) => ({
		label: resolveUserName(user.login, user.displayName),
		value: user.id,
		profileImageUrl: user.profileImageUrl,
	}))
}
const options = computed(() => {
	const searchUsers = twitchSearch.data.value?.channels ?? []
	const initialUsers = getUsers.data.value?.users ?? []

	const allOptions = [
		...mapOptions(searchUsers),
		...mapOptions(initialUsers),
	]

	const uniqueOptions = allOptions.filter((option, index, self) =>
		index === self.findIndex((t) => (
			t.value === option.value
		)),
	)

	return uniqueOptions
})

const open = ref(false)

function handleSelect(event: SelectEvent<AcceptableValue>) {
	if (typeof event.detail.value !== 'string') return
	if (props.multiple && Array.isArray(userId.value)) {
		if (userId.value?.includes(event.detail.value)) {
			userId.value = userId.value?.filter((id) => id !== event.detail.value)
		} else {
			userId.value = [...userId.value ?? [], event.detail.value]
		}
	} else {
		userId.value = event.detail.value
	}

	if (!props.multiple) {
		open.value = false
	}
}

function getCheckedClass(value: string) {
	if (Array.isArray(userId.value)) {
		return userId.value.includes(value) ? 'opacity-100' : 'opacity-0'
	}
	return userId.value === value ? 'opacity-100' : 'opacity-0'
}
</script>

<template>
	<Popover v-model:open="open">
		<PopoverTrigger as-child>
			<Button
				variant="outline"
				role="combobox"
				:aria-expanded="open"
				class="w-full justify-between"
				@click="open = true"
			>
				<template v-if="multiple">
					<div v-if="userId?.length" class="flex gap-2 items-center">
						<span v-if="userId?.length">
							{{ userId?.length }} users selected
						</span>
						<div class="flex flex-row gap-0.5">
							<img v-for="id in userId" :key="id" :src="options.find((option) => option.value === id)?.profileImageUrl" class="size-4 rounded-full" />
						</div>
					</div>
					<span v-else>
						Select users...
					</span>
				</template>
				<template v-else>
					<template v-if="userId">
						<div class="flex gap-2 items-center">
							<img :src="options.find((option) => option.value === userId)?.profileImageUrl" class="size-5 rounded-full" />
							<div>{{ options.find((option) => option.value === userId)?.label }}</div>
						</div>
					</template>
					<span v-else>
						Select user...
					</span>
				</template>
				<ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
			</Button>
		</PopoverTrigger>
		<PopoverContent class="w-full p-0 z-[9999]">
			<Command
				v-model:searchTerm="userName"
				:multiple="multiple"
				:filter-function="(l) => l"
				:reset-search-term-on-blur="false"
			>
				<CommandInput class="h-9" placeholder="Search user..." />
				<CommandEmpty>
					No Users found.
				</CommandEmpty>
				<CommandList>
					<CommandGroup>
						<CommandItem
							v-for="option in options"
							:key="option.value"
							:value="option.value"
							@select="handleSelect"
						>
							<div class="flex gap-2 items-center">
								<img :src="option.profileImageUrl" class="size-5 rounded-full" />
								<span>{{ option.label }}</span>
							</div>
							<Check
								:class="cn(
									'ml-auto h-4 w-4',
									getCheckedClass(option.value),
								)"
							/>
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
