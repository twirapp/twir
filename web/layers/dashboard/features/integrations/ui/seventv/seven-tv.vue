<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import Settings from './settings.vue'

import SevenTVSvg from '@/assets/integrations/seventv.svg?use'
import WithSettings from '@/components/integrations/variants/withSettings.vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import CommandsList from '@/features/commands/ui/list.vue'
import { useSeventvData } from '@/features/integrations/composables/seventv/use-seventv-data.ts'
import { useSeventvSteps } from '@/features/integrations/composables/seventv/use-seventv-steps.ts'
import Steps from '@/features/integrations/ui/seventv/steps/steps.vue'

const { t } = useI18n()

const { userProfile, sevenTvCommands } = useSeventvData()

const { steps, currentStep } = useSeventvSteps()
</script>

<template>
	<WithSettings title="7TV" :icon="SevenTVSvg" icon-width="48px" dialog-content-class="w-[600px]">
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
