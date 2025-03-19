<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { InfoIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { computed, onMounted, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import { useObsOverlayManager, useProfile } from '@/api/index.js'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
	Stepper,
	StepperDescription,
	StepperItem,
	StepperSeparator,
	StepperTrigger,
} from '@/components/ui/stepper'
import { useToast } from '@/components/ui/toast'

const { t } = useI18n()
const { toast } = useToast()

const obsSettingsManager = useObsOverlayManager()
const { refetch, data: settings } = obsSettingsManager.getSettings()
const obsSettingsUpdater = obsSettingsManager.updateSettings()

const schema = z.object({
	serverAddress: z.string().min(1).trim().default('localhost'),
	serverPassword: z.string().min(1).trim(),
	serverPort: z.number().min(1).default(4455),
})

const obsForm = useForm({
	validationSchema: toTypedSchema(schema),
	validateOnMount: false,
	keepValuesOnUnmount: true,
})

onMounted(async () => {
	const settings = await refetch()
	if (!settings.data) return
	obsForm.setValues(toRaw(settings.data))
})

const onSubmit = obsForm.handleSubmit(async (values) => {
	await obsSettingsUpdater.mutateAsync(values)
	toast({
		title: 'Settings updated, now you can paste overlay link into obs',
		duration: 10000,
	})
})

const { copyOverlayLink } = useCopyOverlayLink('obs')
const { data: profile } = useProfile()

const currentStep = computed(() => {
	if (!settings.value) return 1

	return settings.value.isConnected ? 3 : 2
})
</script>

<template>
	<Alert type="info">
		<InfoIcon class="size-4" />
		<AlertTitle>OBS Overlay</AlertTitle>
		<AlertDescription>
			{{ t('overlays.obs.description') }}
		</AlertDescription>
	</Alert>

	<Stepper orientation="vertical" class="flex w-full flex-col justify-start gap-10" :model-value="currentStep">
		<StepperItem :step="1" class="relative flex w-full items-start gap-6">
			<StepperSeparator
				class="absolute left-[18px] top-[38px] block h-[105%] w-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-primary"
			/>

			<StepperTrigger as-child>
				<Button
					variant="outline"
					size="icon"
					class="z-10 rounded-full shrink-0 ring-1 ring-ring ring-offset-2 ring-offset-background"
				>
					1
				</Button>
			</StepperTrigger>

			<StepperDescription class="text-white">
				<span class="text-xl">Configure settings of twir overlay for connect to OBS</span>

				<form class="flex flex-col gap-2 mt-4" @submit="onSubmit">
					<FormField v-slot="{ componentField }" name="serverAddress">
						<FormItem>
							<FormLabel>{{ t('overlays.obs.address') }}</FormLabel>
							<FormControl>
								<Input type="text" placeholder="localhost" v-bind="componentField" />
							</FormControl>
							<FormDescription>
								Usually it's localhost, but we gives you opportunity to connect to remote servers.
							</FormDescription>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="serverPort">
						<FormItem>
							<FormLabel>{{ t('overlays.obs.port') }}</FormLabel>
							<FormControl>
								<Input type="number" placeholder="4455" v-bind="componentField" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="serverPassword">
						<FormItem>
							<FormLabel>Password</FormLabel>
							<FormControl>
								<Input type="password" placeholder="Password" v-bind="componentField" />
							</FormControl>
							<FormMessage />
							<FormDescription>
								Ensure the password matches, as many users misconfigure it and encounter issues when clicking 'generate password' after copying it.
							</FormDescription>
						</FormItem>
					</FormField>

					<Button type="submit">
						{{ t('sharedButtons.save') }}
					</Button>
				</form>
			</StepperDescription>
		</StepperItem>

		<StepperItem :step="2" class="relative flex w-full items-start gap-6">
			<StepperSeparator
				class="absolute left-[18px] top-[38px] block h-[105%] w-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-primary"
			/>

			<StepperTrigger as-child>
				<Button
					variant="outline"
					size="icon"
					class="z-10 rounded-full shrink-0 ring-offset-2 ring-offset-background"
					:class="currentStep >= 2 && 'ring-1 ring-ring'"
				>
					2
				</Button>
			</StepperTrigger>

			<StepperDescription class="text-white w-full">
				<span class="text-xl">Add overlay to OBS</span>

				<div class="flex flex-col gap-2">
					<span>Add overlay on you scene or scenes, if you need functinality on all scenes.</span>
					<Button
						:disabled="profile?.id !== profile?.selectedDashboardId"
						variant="secondary"
						@click="copyOverlayLink()"
					>
						{{ t('overlays.copyOverlayLink') }}
					</Button>
				</div>
			</StepperDescription>
		</StepperItem>

		<StepperItem :step="3" class="relative flex w-full items-start gap-6">
			<StepperTrigger as-child>
				<Button
					variant="outline"
					size="icon"
					class="z-10 rounded-full shrink-0 ring-offset-2 ring-offset-background"
					:class="currentStep === 3 && 'ring-1 ring-ring'"
				>
					3
				</Button>
			</StepperTrigger>

			<StepperDescription class="text-white w-full">
				<Alert
					type="info"
					class="w-full"
					:class="settings?.isConnected ? 'bg-green-800/50' : 'bg-red-800/50'"
				>
					<AlertTitle>State</AlertTitle>

					<AlertDescription>
						<span v-if="settings?.isConnected">{{ t('overlays.obs.connected') }}</span>
						<span v-else>{{ t('overlays.obs.notConnected') }}</span>
					</AlertDescription>
				</Alert>
			</StepperDescription>
		</StepperItem>
	</Stepper>
</template>
