<script setup lang="ts">
import { BadgePlus, SaveIcon, TrashIcon } from 'lucide-vue-next'
import type { CommandGroup } from '@/gql/graphql'

import { useCommandsGroupsApi } from '#layers/dashboard/api/commands/commands-groups'

const { t } = useI18n()

const groupsManager = useCommandsGroupsApi()
const { data } = groupsManager.useQueryGroups()
const groupsCreator = groupsManager.useMutationCreateGroup()
const groupsDeleter = groupsManager.useMutationDeleteGroup()
const groupsUpdater = groupsManager.useMutationUpdateGroup()

type FormGroup = Omit<CommandGroup, 'id'> & { id?: string }

const groups = ref<FormGroup[]>([])

watch(data, (data) => {
	groups.value = data?.commandsGroups ? toRaw(data?.commandsGroups) : []
}, { immediate: true })

async function create(name: string, color: string) {
	await groupsCreator.executeMutation({ opts: { color, name } })
}

async function deleteGroup(index: number) {
	const group = groups.value[index]
	if (!group?.id) return
	await groupsDeleter.executeMutation({ id: group.id })
}

async function update(index: number) {
	const group = groups.value[index]
	if (!group?.id) return
	await groupsUpdater.executeMutation({ id: group.id, opts: { name: group.name, color: group.color } })
}

const swatches = [
	'#74f2ca',
	'#d03050',
]
</script>

<template>
	<div class="flex flex-col gap-2">
		<UiAlert v-if="!groups.length">
			<UiAlertDescription>
				No groups created yet
			</UiAlertDescription>
		</UiAlert>
		<div
			v-for="(group, index) of groups"
			v-else
			:key="group.id"
			class="flex flex-wrap flex-col md:flex-row gap-2 border border-border py-2 px-4 rounded-md"
		>
			<div class="flex flex-col gap-1 w-full md:w-[40%]">
				<span>Name</span>
				<UiInput v-model:model-value="group.name" />
			</div>
			<div class="flex flex-col gap-1 w-full md:w-[40%]">
				<span>Color</span>
				<ColorPicker
					v-model="group.color"
					:presets="swatches"
					show-presets
				/>
			</div>

			<div class="flex gap-2 items-end justify-end">
				<UiButton size="icon" @click="update(index)">
					<SaveIcon />
				</UiButton>
				<UiButton size="icon" variant="destructive" @click="deleteGroup(index)">
					<TrashIcon />
				</UiButton>
			</div>
		</div>
	</div>

	<UiButton
		variant="outline"
		class="mt-2 w-full flex gap-2 items-center"
		:disabled="groups.length >= 10"
		@click="() => create(`New Group #${groups.length + 1}`, swatches[0]!)"
	>
		<BadgePlus class="size-4" />
		<span>{{ t('sharedButtons.create') }} ({{ groups.length }}/10)</span>
	</UiButton>
</template>
