<script setup lang="ts">
import { NModal } from 'naive-ui'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import List from './ui/list.vue'

import { useUserAccessFlagChecker } from '@/api'
import { useCommandsApi } from '@/api/commands/commands.js'
import ManageGroups from '@/components/commands/manageGroups.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import PageLayout from '@/layout/page-layout.vue'

const route = useRoute()
const { t } = useI18n()
const userCanManageCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageCommands)

const commandsManager = useCommandsApi()
const { data: commandsResponse } = commandsManager.useQueryCommands()

const excludedModules = ['7tv']

const commandsFilter = ref('')
const commands = computed(() => {
	if (!commandsResponse.value?.commands) return []

	const system = Array.isArray(route.params.system) ? route.params.system[0] : route.params.system

	return commandsResponse.value.commands
		.filter(c => {
			if (system.toUpperCase() === 'CUSTOM') {
				return c.module === 'CUSTOM'
			}

			return c.module !== 'CUSTOM' && !excludedModules.includes(c.module)
		})
		.filter(c => {
			return c.name.includes(commandsFilter.value) || c.aliases.some(a => a.includes(commandsFilter.value))
		})
})

const showManageGroupsModal = ref(false)

const isCreateDisabled = computed(() => {
	return commands.value.length >= 50 || !userCanManageCommands.value
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
		<template #title>
			{{ title }} {{ t('sidebar.commands.label').toLocaleLowerCase() }}
		</template>

		<template #title-footer>
			<Input v-model="commandsFilter" placeholder="Search..." size="" />
		</template>

		<template v-if="isCustom" #action>
			<div class="flex gap-2 flex-wrap">
				<Button variant="secondary" @click="showManageGroupsModal = true">
					{{ t('commands.groups.manageButton') }}
				</Button>
				<RouterLink v-slot="{ href, navigate }" custom to="/dashboard/commands/custom/create">
					<Button as="a" :href="href" :disabled="isCreateDisabled" @click="navigate">
						{{ t('sharedButtons.create') }} ({{ commands.length }}/50)
					</Button>
				</RouterLink>
			</div>
		</template>

		<template #content>
			<List
				:commands="commands"
				show-background
				enable-groups
			/>
		</template>
	</PageLayout>

	<NModal
		v-model:show="showManageGroupsModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="t('commands.groups.manageButton')"
		class="modal"
		:style="{
			width: '600px',
		}"
		:on-close="() => showManageGroupsModal = false"
	>
		<ManageGroups />
	</NModal>
</template>
