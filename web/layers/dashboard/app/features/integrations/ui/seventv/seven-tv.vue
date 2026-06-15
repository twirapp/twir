<script setup lang="ts">

import Settings from './settings.vue'

import WithSettings from '~~/layers/dashboard/app/components/integrations/variants/withSettings.vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import CommandsList from '~~/layers/dashboard/app/features/commands/ui/list.vue'
import { useSeventvData } from '~~/layers/dashboard/app/features/integrations/composables/seventv/use-seventv-data.js'
import { useSeventvSteps } from '~~/layers/dashboard/app/features/integrations/composables/seventv/use-seventv-steps.js'
import Steps from '~~/layers/dashboard/app/features/integrations/ui/seventv/steps/steps.vue'

const { t } = useI18n()

const { userProfile, sevenTvCommands } = useSeventvData()

const { steps, currentStep } = useSeventvSteps()
</script>

<template>
	<WithSettings title="7TV" icon="twir-integrations:seventv" icon-width="48px" dialog-content-class="w-[600px]">
		<template #description>
			{{ t('integrations.sevenTv.description') }}
		</template>

		<template #settings>
			<Steps v-if="currentStep !== -1" />

			<Tabs v-else default-value="settings" class="flex flex-col w-full">
				<TabsList class="ml-auto">
					<TabsTrigger value="settings"> Settings </TabsTrigger>
					<TabsTrigger value="commands">
						{{ t('sidebar.commands.label') }}
					</TabsTrigger>
				</TabsList>
				<TabsContent value="settings">
					<Settings />
				</TabsContent>
				<TabsContent value="commands">
					<CommandsList v-if="sevenTvCommands" :commands="sevenTvCommands" show-background />
				</TabsContent>
			</Tabs>
		</template>

		<template #additionalFooter>
			<div
				v-if="steps.every((s) => s.completed) && userProfile"
				class="flex items-center gap-2 p-2 border-2 border-gray-700 rounded-md px-4"
			>
				<img :src="userProfile.avatarUri" class="h-5 w-5 rounded-full" />
				<span class="text-sm font-medium">{{ userProfile.displayName }}</span>
			</div>
			<div v-else class="flex items-center gap-2 p-2 bg-destructive/50 rounded-md px-4">
				Not configured
			</div>
		</template>
	</WithSettings>
</template>
