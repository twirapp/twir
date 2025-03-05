<script setup lang="ts">
import { Bomb } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { onMounted, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { formSchema, useDuelForm } from './composables/use-duel-form'

import Card from '@/components/games/card.vue'
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
import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import CommandButton from '@/features/commands/ui/command-button.vue'

const isDialogOpen = ref(false)
const { t } = useI18n()
const { initialValues, save, settings } = useDuelForm()

const duelForm = useForm({
	validationSchema: formSchema,
	initialValues,
	validateOnMount: false,
	keepValuesOnUnmount: true,
})

watch(settings, (newSettings) => {
	if (!newSettings) return
	duelForm.setValues(toRaw(newSettings.gamesDuel))
}, { immediate: true })

onMounted(() => {
	if (!settings.value) return
	duelForm.setValues(structuredClone(toRaw(settings.value.gamesDuel)))
})

const onSubmit = duelForm.handleSubmit(async (values) => {
	await save(values)
	isDialogOpen.value = false
})

function resetSettings() {
	duelForm.setValues(initialValues)
}
</script>

<template>
	<Dialog v-model:open="isDialogOpen">
		<DialogTrigger asChild>
			<Card
				title="Duel"
				:icon="Bomb"
				:icon-stroke="1"
				:description="t('games.duel.description')"
			/>
		</DialogTrigger>

		<DialogContent class="sm:max-w-[625px]">
			<DialogHeader>
				<DialogTitle>{{ t('games.duel.title') }}</DialogTitle>
			</DialogHeader>

			<form>
				<div class="grid gap-4 py-4">
					<FormField
						v-slot="{ value, handleChange }"
						name="enabled"
					>
						<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4 space-y-0">
							<FormLabel>Enabled</FormLabel>
							<FormControl>
								<Switch
									:checked="value"
									@update:checked="handleChange"
								/>
							</FormControl>
						</FormItem>
					</FormField>

					<div class="rounded-lg border p-4">
						<h4 class="mb-4 text-sm font-medium">
							{{ t('games.duel.commands.title') }}
						</h4>
						<div class="flex gap-2">
							<CommandButton name="duel" :title="t('games.duel.commands.duel')" />
							<CommandButton name="duel accept" :title="t('games.duel.commands.accept')" />
							<CommandButton name="duel stats" :title="t('games.duel.commands.stats')" />
						</div>
					</div>

					<div class="rounded-lg border p-4">
						<h4 class="mb-4 text-sm font-medium">
							{{ t('games.duel.cooldown.title') }}
						</h4>
						<div class="grid grid-cols-2 gap-4">
							<FormField
								v-slot="{ componentField }"
								name="userCooldown"
							>
								<FormItem>
									<FormLabel>{{ t('games.duel.cooldown.user') }}</FormLabel>
									<FormControl>
										<Input type="number" v-bind="componentField" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								v-slot="{ componentField }"
								name="globalCooldown"
							>
								<FormItem>
									<FormLabel>{{ t('games.duel.cooldown.global') }}</FormLabel>
									<FormControl>
										<Input type="number" v-bind="componentField" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</div>

					<div class="rounded-lg border p-4">
						<h4 class="mb-4 text-sm font-medium">
							{{ t('games.duel.messages.title') }}
						</h4>
						<div class="space-y-4">
							<FormField
								v-slot="{ componentField }"
								name="startMessage"
							>
								<FormItem>
									<FormLabel>{{ t('games.duel.messages.start.title') }}</FormLabel>
									<FormControl>
										<Input v-bind="componentField" />
									</FormControl>
									<p class="text-sm text-muted-foreground">
										{{ t('games.duel.messages.start.description') }}
									</p>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								v-slot="{ componentField }"
								name="resultMessage"
							>
								<FormItem>
									<FormLabel>{{ t('games.duel.messages.result.title') }}</FormLabel>
									<FormControl>
										<Input v-bind="componentField" />
									</FormControl>
									<p class="text-sm text-muted-foreground">
										{{ t('games.duel.messages.result.description') }}
									</p>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField
								v-slot="{ componentField }"
								name="bothDieMessage"
							>
								<FormItem>
									<FormLabel>{{ t('games.duel.messages.bothDie.title') }}</FormLabel>
									<FormControl>
										<Input v-bind="componentField" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>
						</div>
					</div>

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
									<AlertDialogDescription>
										Are you sure?
									</AlertDialogDescription>
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

						<Button type="submit" @click="onSubmit">
							{{ t('sharedButtons.save') }}
						</Button>
					</div>
				</div>
			</form>
		</DialogContent>
	</Dialog>
</template>
