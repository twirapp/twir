<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { SkullIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { onMounted, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import Card from './card.vue'

import type { SeppukuGame } from '@/gql/graphql'

import { useGamesApi } from '@/api/games/games'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { toast } from 'vue-sonner'
import CommandButton from '@/features/commands/ui/command-button.vue'

const formSchema = toTypedSchema(
	z.object({
		enabled: z.boolean(),
		message: z.string().max(500),
		messageModerators: z.string().max(500),
		timeoutModerators: z.boolean(),
		timeoutSeconds: z.number().min(1).max(86400),
	})
)

const isModalOpened = ref(false)
const { t } = useI18n()

const gamesManager = useGamesApi()
const { data } = gamesManager.useGamesQuery()
const updater = gamesManager.useSeppukuMutation()

const initialValues = {
	enabled: false,
	message:
		'{sender} said: my honor tarnished, I reclaim it through death. May my spirit find peace. Farewell.',
	messageModerators: '{sender} drew his sword and ripped open his belly for the sad emperor.',
	timeoutModerators: false,
	timeoutSeconds: 60,
}

const seppukuForm = useForm<SeppukuGame>({
	validationSchema: formSchema,
	initialValues,
	validateOnMount: false,
	keepValuesOnUnmount: true,
})

watch(
	data,
	(v) => {
		if (!v) return
		const raw = toRaw(v)
		seppukuForm.setValues({
			enabled: raw.gamesSeppuku.enabled,
			message: raw.gamesSeppuku.message,
			messageModerators: raw.gamesSeppuku.messageModerators,
			timeoutModerators: raw.gamesSeppuku.timeoutModerators,
			timeoutSeconds: raw.gamesSeppuku.timeoutSeconds,
		})
	},
	{ immediate: true }
)

onMounted(() => {
	if (!data.value) return
	seppukuForm.setValues(structuredClone(toRaw(data.value.gamesSeppuku)))
})

const save = seppukuForm.handleSubmit(async () => {
	await updater.executeMutation({
		opts: seppukuForm.values,
	})
	toast.success(t('sharedTexts.saved'))
	isModalOpened.value = false
})

function resetSettings() {
	seppukuForm.setValues(initialValues)
}
</script>

<template>
	<Dialog v-model:open="isModalOpened">
		<DialogTrigger asChild>
			<Card
				title="Seppuku"
				:icon="SkullIcon"
				:icon-stroke="1"
				:description="t('games.seppuku.description')"
			/>
		</DialogTrigger>

		<DialogContent class="sm:max-w-[500px]">
			<DialogHeader>
				<DialogTitle>Seppuku</DialogTitle>
			</DialogHeader>

			<form>
				<div class="space-y-4">
					<div class="flex gap-4 flex-col">
						<FormField v-slot="{ value, handleChange }" name="enabled">
							<FormItem class="flex gap-2 space-y-0 items-center">
								<FormLabel>
									{{ t('sharedTexts.enabled') }}
								</FormLabel>
								<FormControl>
									<Switch :checked="value" @update:checked="handleChange" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<CommandButton name="seppuku" />
					</div>

					<Separator />

					<FormField v-slot="{ componentField }" name="message">
						<FormItem>
							<FormLabel>{{ t('games.seppuku.message') }}</FormLabel>
							<FormControl>
								<Input v-bind="componentField" :maxlength="500" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="timeoutSeconds">
						<FormItem>
							<FormLabel>{{ t('games.seppuku.timeoutSeconds') }}</FormLabel>
							<FormControl>
								<Input v-bind="componentField" type="number" :min="1" :max="86400" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ value, handleChange }" name="timeoutModerators">
						<FormItem>
							<FormLabel>{{ t('games.seppuku.timeoutModerators') }}</FormLabel>
							<FormControl>
								<Switch :checked="value" @update:checked="handleChange" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="messageModerators">
						<FormItem>
							<FormLabel>{{ t('games.seppuku.messageModerators') }}</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
									:maxlength="500"
									:disabled="!seppukuForm.values.timeoutModerators"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>

				<Separator class="my-4" />

				<div class="flex justify-between">
					<AlertDialog>
						<AlertDialogTrigger asChild>
							<Button variant="outline">
								{{ t('sharedButtons.setDefaultSettings') }}
							</Button>
						</AlertDialogTrigger>
						<AlertDialogContent>
							<AlertDialogHeader>
								<AlertDialogTitle>{{ t('sharedTexts.areYouSure') }}</AlertDialogTitle>
								<AlertDialogDescription> Are you sure? </AlertDialogDescription>
							</AlertDialogHeader>
							<AlertDialogFooter>
								<AlertDialogCancel>
									{{ t('sharedButtons.cancel') }}
								</AlertDialogCancel>
								<AlertDialogAction @click="resetSettings">
									{{ t('sharedButtons.confirm') }}
								</AlertDialogAction>
							</AlertDialogFooter>
						</AlertDialogContent>
					</AlertDialog>

					<Button type="submit" @click="save">
						{{ t('sharedButtons.save') }}
					</Button>
				</div>
			</form>
		</DialogContent>
	</Dialog>
</template>
