<script setup lang="ts">
import { TrashIcon } from 'lucide-vue-next'
import {
	NAlert,
	NCard,
} from 'naive-ui'
import { type VNode, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { type FormKey, useForm } from '../composables/use-form.js'

import { useUserAccessFlagChecker } from '@/api/index.js'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { ChannelRolePermissionEnum } from '@/gql/graphql.js'

const props = defineProps<{
	formKey: FormKey
	alertMessage?: string
	count?: {
		label: string
	}
	maxMessages: number
	defaultMessageText: string
	minCount?: number
	minCooldown: number
}>()

defineSlots<{
	additionalSettings: VNode
}>()

const { formValue, formInited, formRef } = useForm()
const hasAccessToManageAlerts = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageAlerts)

watch(formInited, (v) => {
	if (!v) return

	if (!formValue?.value?.[props.formKey]?.messages.length) {
		createMessage()
	}
}, { immediate: true })

function createMessage() {
	if (!hasAccessToManageAlerts) return

	if (props.count) {
		const latest = formValue?.value?.[props.formKey]?.messages.at(-1)
		const countForSet = latest && 'count' in latest ? latest.count + 1 : 1

		formValue?.value?.[props.formKey]?.messages.push({
			count: countForSet,
			text: props.defaultMessageText,
		})
	} else {
		formValue?.value?.[props.formKey]?.messages.push({ text: props.defaultMessageText } as any)
	}
}

function removeMessage(index: number) {
	if (!hasAccessToManageAlerts) return
	if (!formValue?.value?.[props.formKey]?.messages) return

	// eslint-disable-next-line ts/ban-ts-comment
	// @ts-expect-error
	formValue.value[props.formKey]!.messages = formValue.value[props.formKey]!.messages.filter((_, i) => i !== index)
}

const { t } = useI18n()
</script>

<template>
	<form
		v-if="formValue?.[formKey]"
		ref="formRef"
		class="flex flex-col gap-4"
		@submit.prevent
	>
		<div class="relative">
			<div v-if="!hasAccessToManageAlerts" class="absolute w-full h-full z-50">
				<div class="flex flex-col items-center justify-center h-full gap-2">
					<h2
						class="scroll-m-20 text-3xl font-semibold tracking-tight transition-colors first:mt-0"
					>
						{{ t('haveNoAccess.title') }}
					</h2>
					<p class="leading-7 [&:not(:first-child)]:mt-6">
						{{ t('haveNoAccess.description') }}
					</p>
				</div>
			</div>

			<NCard :class="{ 'absolute opacity-20': !hasAccessToManageAlerts }" :title="t('sharedButtons.settings')" size="small" bordered>
				<div class="flex flex-col gap-4">
					<div class="flex items-center gap-4">
						<Label for="enabled">{{ t('sharedTexts.enabled') }}</Label>
						<Switch
							id="enabled"
							:checked="formValue[formKey]!.enabled"
							@update:checked="(v) => formValue[formKey]!.enabled = v"
						/>
					</div>

					<div class="grid items-center gap-1.5">
						<Label for="cooldown">{{ t('chatAlerts.cooldown') }}</Label>
						<Input
							id="cooldown"
							v-model="formValue[formKey]!.cooldown"
							:min="minCooldown"
							:max="9999"
							class="w-[10%] min-w-[100px]"
							type="number"
							pattern="\d*"
						/>
					</div>

					<slot name="additionalSettings" />
				</div>
			</NCard>
		</div>

		<div class="relative">
			<div v-if="!hasAccessToManageAlerts" class="absolute w-full h-full z-50">
				<div class="flex flex-col items-center justify-center h-full gap-2">
					<h2
						class="scroll-m-20 text-3xl font-semibold tracking-tight transition-colors first:mt-0"
					>
						{{ t('haveNoAccess.title') }}
					</h2>
					<p class="leading-7 [&:not(:first-child)]:mt-6">
						{{ t('haveNoAccess.description') }}
					</p>
				</div>
			</div>

			<NCard
				:class="{ 'absolute opacity-20': !hasAccessToManageAlerts }"
				:title="t('sharedTexts.messages')"
				size="small"
				bordered
			>
				<template #header-extra>
					<Button
						:disabled="(formValue[formKey]!.messages?.length === maxMessages) || !hasAccessToManageAlerts"
						variant="default"
						@click="createMessage"
					>
						<span v-if="formValue[formKey]!.messages?.length">{{ t('sharedButtons.create') }} ({{ formValue[formKey]!.messages.length }} / {{ maxMessages }})</span>
						<span v-else>{{ t('sharedButtons.create') }}</span>
					</Button>
				</template>

				<p class="leading-7 [&:not(:first-child)]:mt-6" v-html="alertMessage" />

				<NAlert v-if="!formValue[formKey]!.messages?.length" type="warning">
					No messages
				</NAlert>

				<ul v-else class="flex flex-col gap-3.5 p-0 mx-0 my-3.5">
					<li
						v-for="(message, index) of formValue[formKey]!.messages"
						:key="index"
						class="flex justify-between gap-3.5"
					>
						<div class="flex w-full gap-x-3.5 gap-y-2 items-end">
							<div v-if="count && 'count' in message" class="grid max-w-sm items-center gap-1.5">
								<Label for="count">{{ count.label }} >=</Label>
								<Input
									id="count"
									v-model="message.count"
									:min="minCount ?? 1"
									:max="9999999"
									class="flex-1"
									type="number"
								/>
							</div>

							<Input
								v-model="message.text"
							/>

							<Button
								:disabled="!hasAccessToManageAlerts"
								variant="secondary"
								size="icon"
								@click="removeMessage(index)"
							>
								<TrashIcon class="h-4 w-4" />
							</Button>
						</div>
					</li>
				</ul>
			</NCard>
		</div>
	</form>
</template>
