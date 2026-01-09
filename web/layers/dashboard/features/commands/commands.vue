<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { computed, ref } from 'vue'

import { useRoute } from 'vue-router'

import List from './ui/list.vue'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useCommandsApi } from '#layers/dashboard/api/commands/commands'
import ManageGroups from '#layers/dashboard/components/commands/manageGroups.vue'
import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'



import { ChannelRolePermissionEnum } from '~/gql/graphql'
import PageLayout from '~/layout/page-layout.vue'

const route = useRoute()
const { t } = useI18n()
const userCanManageCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageCommands)

const { user: profile } = storeToRefs(useDashboardAuth())
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
	return selectedDashboard?.plan?.maxCommands ?? 50
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
				<UiDialog>
					<UiDialogTrigger as-child>
						<UiButton variant="secondary" @click="showManageGroupsModal = true">
							{{ t('commands.groups.manageButton') }}
						</UiButton>
					</UiDialogTrigger>
					<DialogOrSheet>
						<UiDialogTitle>{{ t('commands.groups.manageButton') }}</UiDialogTitle>
						<ManageGroups />
					</DialogOrSheet>
				</UiDialog>
				<RouterLink v-slot="{ href, navigate }" custom to="/dashboard/commands/custom/create">
					<UiButton as="a" :href="href" :disabled="isCreateDisabled" @click="navigate">
						<PlusIcon class="size-4 mr-2" />
						{{ t('sharedButtons.create') }} ({{ commands.length }}/{{ maxCommands }})
					</UiButton>
				</RouterLink>
			</div>
		</template>

		<template #content>
			<div class="flex flex-col gap-2">
				<UiInput v-model="commandsFilter" placeholder="Search..." class="w-full lg:w-[40%]" />
				<List :commands="commands" show-background enable-groups />
			</div>
		</template>
	</PageLayout>
</template>
