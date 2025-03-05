<script setup lang="ts">
import { ref } from 'vue'

import type { FunctionalComponent } from 'vue'

import { useUserAccessFlagChecker } from '@/api'
import { Card, CardFooter, CardHeader } from '@/components/ui/card'
import SettingsModal from '@/components/ui/settings-modal.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = withDefaults(defineProps<{
	title: string
	icon?: FunctionalComponent
	iconFill?: string
	save?: () => void | Promise<void>
	isLoading?: boolean
	saveDisabled?: boolean
}>(), {})

defineSlots<{
	settings: FunctionalComponent
	customDescriptionSlot?: FunctionalComponent
	additionalFooter?: FunctionalComponent
	description: FunctionalComponent
}>()

const showSettings = ref(false)
const userCanManageIntegrations = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageIntegrations)

async function callSave() {
	if (props.save) {
		await props.save()
	}
	showSettings.value = false
}
</script>

<template>
	<Card class="flex flex-col h-full">
		<CardHeader>
			<div class="flex items-center gap-2">
				<component
					:is="icon"
					:style="{ fill: iconFill }"
					class="w-8 h-8"
				/>
				<h3 class="text-lg font-medium">
					{{ title }}
				</h3>
			</div>
		</CardHeader>

		<div class="flex-1 p-6">
			<slot name="description" />
			<slot name="customDescriptionSlot" />
		</div>

		<CardFooter class="mt-auto">
			<div class="flex justify-between flex-wrap w-full items-center gap-2">
				<div>
					<SettingsModal
						v-model:open="showSettings"
						:title="title"
						:show-save="!!save"
						:save-disabled="saveDisabled"
						:button-disabled="!userCanManageIntegrations"
						@save="callSave"
					>
						<slot name="settings" />
					</SettingsModal>
				</div>

				<div>
					<slot name="additionalFooter" />
				</div>
			</div>
		</CardFooter>
	</Card>
</template>
