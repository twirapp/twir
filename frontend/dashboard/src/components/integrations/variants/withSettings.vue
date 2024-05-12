<script setup lang="ts">
import { IconSettings } from '@tabler/icons-vue'
import { NButton, NModal, NSpace } from 'naive-ui'
import { onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import type { FunctionalComponent } from 'vue'

import { useUserAccessFlagChecker } from '@/api'
import Card from '@/components/card/card.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = withDefaults(defineProps<{
	title: string
	icon?: FunctionalComponent
	iconFill?: string
	save?: () => void | Promise<void>
	modalWidth?: string
	iconWidth?: string
	isLoading?: boolean
	saveDisabled?: boolean
	modalContentStyle?: string
}>(), {
	modalWidth: '600px',
})

defineSlots<{
	settings: FunctionalComponent
	customDescriptionSlot?: FunctionalComponent
	additionalFooter?: FunctionalComponent
	description?: FunctionalComponent | string
}>()

const showSettings = ref(false)

async function callSave() {
	await props.save?.()
}

const userCanManageIntegrations = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageIntegrations)

const { t } = useI18n()

onUnmounted(() => showSettings.value = false)
</script>

<template>
	<Card
		:title="title"
		style="height: 100%;"
		:with-stroke="false"
		:icon="icon"
		:icon-fill="iconFill"
		:icon-width="iconWidth"
		:is-loading="isLoading"
	>
		<template #content>
			<slot name="customDescriptionSlot" />
			<slot class="description" name="description" />
		</template>

		<template #footer>
			<div class="flex justify-between flex-wrap w-full items-center gap-2">
				<div>
					<NButton
						:disabled="!userCanManageIntegrations"
						secondary
						size="large"
						@click="showSettings = true"
					>
						<div class="flex gap-1">
							<span>{{ t('sharedButtons.settings') }}</span>
							<IconSettings />
						</div>
					</NButton>
				</div>

				<div>
					<slot name="additionalFooter" />
				</div>
			</div>
		</template>
	</Card>

	<NModal
		v-model:show="showSettings"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="title"
		class="modal"
		:style="{
			width: modalWidth,
			top: '5%',
			bottom: '5%',
		}"
		:content-style="modalContentStyle"
	>
		<template #header>
			{{ title }}
		</template>
		<slot name="settings" />

		<template #action>
			<NSpace justify="end">
				<NButton secondary @click="showSettings = false">
					{{ t('sharedButtons.close') }}
				</NButton>
				<NButton v-if="save" secondary type="success" :disabled="saveDisabled" @click="callSave">
					{{ t('sharedButtons.save') }}
				</NButton>
			</NSpace>
		</template>
	</NModal>
</template>
