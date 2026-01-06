<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import List from './ui/list.vue'

import { useProfile, useUserAccessFlagChecker } from '@/api'
import { useCommandsApi } from '@/api/commands/commands.js'
import ManageGroups from '@/components/commands/manageGroups.vue'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import PageLayout from '@/layout/page-layout.vue'

const route = useRoute()
const { t } = useI18n()
const userCanManageCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageCommands)

const { data: profile } = useProfile()
const commandsManager = useCommandsApi()
const { data: commandsResponse } = commandsManager.useQueryCommands()

const excludedModules: string[] = []

const commandsFilter = ref('')
const commands = computed(() => {
	if (!commandsResponse.value?.commands) return []

	const system =
		(Array.isArray(route.params.system) ? route.params.system[0] : route.params.system) ?? ''

	return commandsResponse.value.commands
		.filter((c) => {
			if (system.toUpperCase() === 'CUSTOM') {
				return c.module === 'CUSTOM'
			}

			return c.module !== 'CUSTOM' && !excludedModules.includes(c.module)
		})
		.filter((c) => {
			return (
				c.name.includes(commandsFilter.value) ||
				c.aliases.some((a) => a.includes(commandsFilter.value))
			)
		})
})

const showManageGroupsModal = ref(false)

const maxCommands = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan.maxCommands ?? 50
})

const isCreateDisabled = computed(() => {
	return commands.value.length >= maxCommands.value || !userCanManageCommands.value
})

const isCustom = computed(() => {
	return typeof route.params.system === 'string' && route.params.system.toLowerCase() === 'custom'
})

const title = computed(() => {
	if (isCustom.value) {
		return t('sidebar.commands.custom')
	}

	return t('sidebar.commands.builtin')
})
</script>

<template>
	<PageLayout>
		<template #title> {{ title }} {{ t('sidebar.commands.label').toLocaleLowerCase() }} </template>

		<template v-if="isCustom" #action>
			<div class="flex gap-2 flex-wrap">
				<Dialog>
					<DialogTrigger as-child>
						<Button variant="secondary" @click="showManageGroupsModal = true">
							{{ t('commands.groups.manageButton') }}
						</Button>
					</DialogTrigger>
					<DialogOrSheet>
						<DialogTitle>{{ t('commands.groups.manageButton') }}</DialogTitle>
						<ManageGroups />
					</DialogOrSheet>
				</Dialog>
				<RouterLink v-slot="{ href, navigate }" custom to="/dashboard/commands/custom/create">
					<Button as="a" :href="href" :disabled="isCreateDisabled" @click="navigate">
						<PlusIcon class="size-4 mr-2" />
						{{ t('sharedButtons.create') }} ({{ commands.length }}/{{ maxCommands }})
					</Button>
				</RouterLink>
			</div>
		</template>

		<template #content>
			<div class="flex flex-col gap-2">
				<Input v-model="commandsFilter" placeholder="Search..." class="w-full lg:w-[40%]" />
				<List :commands="commands" show-background enable-groups />
			</div>
		</template>
	</PageLayout>
</template>
