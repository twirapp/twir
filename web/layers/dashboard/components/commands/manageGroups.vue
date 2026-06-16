<script setup lang="ts">
import { ref, toRaw, watch } from 'vue'
import { useCommandsGroupsApi } from '~~/layers/dashboard/api/commands/commands-groups'

import type { CommandGroup } from '~/gql/graphql.js'

import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { ColorPicker } from '@/components/ui/color-picker'
import { Input } from '@/components/ui/input'

const { t } = useI18n()

const groupsManager = useCommandsGroupsApi()
const { data } = groupsManager.useQueryGroups()
const groupsCreator = groupsManager.useMutationCreateGroup()
const groupsDeleter = groupsManager.useMutationDeleteGroup()
const groupsUpdater = groupsManager.useMutationUpdateGroup()

type FormGroup = Omit<CommandGroup, 'id'> & { id?: string }

const groups = ref<FormGroup[]>([])

watch(
	data,
	(data) => {
		groups.value = data?.commandsGroups ? toRaw(data?.commandsGroups) : []
	},
	{ immediate: true }
)

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
	await groupsUpdater.executeMutation({
		id: group.id,
		opts: { name: group.name, color: group.color },
	})
}

const swatches = ['#74f2ca', '#d03050']
</script>

<template>
	<div class="flex flex-col gap-2">
		<Alert v-if="!groups.length">
			<AlertDescription> No groups created yet </AlertDescription>
		</Alert>
		<div
			v-for="(group, index) of groups"
			v-else
			:key="group.id"
			class="border-border flex flex-col flex-wrap gap-2 rounded-md border px-4 py-2 md:flex-row"
		>
			<div class="flex w-full flex-col gap-1 md:w-[40%]">
				<span>Name</span>
				<Input v-model:model-value="group.name" />
			</div>
			<div class="flex w-full flex-col gap-1 md:w-[40%]">
				<span>Color</span>
				<ColorPicker
					v-model="group.color"
					:presets="swatches"
					show-presets
				/>
			</div>

			<div class="flex items-end justify-end gap-2">
				<Button
					size="icon"
					@click="update(index)"
				>
					<Icon name="lucide:save" />
				</Button>
				<Button
					size="icon"
					variant="destructive"
					@click="deleteGroup(index)"
				>
					<Icon name="lucide:trash" />
				</Button>
			</div>
		</div>
	</div>

	<Button
		variant="outline"
		class="mt-2 flex w-full items-center gap-2"
		:disabled="groups.length >= 10"
		@click="() => create(`New Group #${groups.length + 1}`, swatches[0])"
	>
		<Icon
			name="lucide:badge-plus"
			class="size-4"
		/>
		<span>{{ t('sharedButtons.create') }} ({{ groups.length }}/10)</span>
	</Button>
</template>
