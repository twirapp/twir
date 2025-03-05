<script setup lang="ts">
import { Bomb } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { onMounted, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import Card from './card.vue'

import type { GamesQuery } from '@/gql/graphql'

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
import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { useToast } from '@/components/ui/toast/use-toast'
import CommandButton from '@/features/commands/ui/command-button.vue'

const isModalOpened = ref(false)

const gamesManager = useGamesApi()
const { data: settings } = gamesManager.useGamesQuery()
const updater = gamesManager.useRussianRouletteMutation()

type Duel = Omit<GamesQuery['gamesRussianRoulette'], '__typename'>

const initialSettings: Duel = {
	enabled: false,
	canBeUsedByModerator: false,
	timeoutSeconds: 60,
	decisionSeconds: 2,
	chargedBullets: 1,
	initMessage: '{sender} has initiated a game of roulette. Is luck on their side?',
	surviveMessage: '{sender} survives the game of roulette! Luck smiles upon them.',
	deathMessage: `{sender} couldn't make it through the game of roulette. Unfortunately, luck wasn't on their side this time.`,
	tumberSize: 6,
}

const rouletteForm = useForm<Duel>({
	initialValues: initialSettings,
	validateOnMount: false,
	keepValuesOnUnmount: true,
})

onMounted(() => {
	if (!settings.value) return
	rouletteForm.setValues(structuredClone(toRaw(settings.value.gamesRussianRoulette)))
})

watch(settings, (v) => {
	if (!v) return
	const raw = toRaw(v)
	rouletteForm.setValues(structuredClone(raw.gamesRussianRoulette))
}, { immediate: true })

const { t } = useI18n()
const { toast } = useToast()

const save = rouletteForm.handleSubmit(async () => {
	await updater.executeMutation({ opts: rouletteForm.values })
	toast({
		title: t('sharedTexts.saved'),
		variant: 'success',
	})
	isModalOpened.value = false
})

async function resetSettings() {
	rouletteForm.setValues(initialSettings)
}
</script>

<template>
	<Dialog v-model:open="isModalOpened">
		<DialogTrigger asChild>
			<Card
				title="Russian Roulette"
				:icon="Bomb"
				:icon-stroke="1"
				:description="t('games.russianRoulette.description')"
			/>
		</DialogTrigger>

		<DialogContent class="sm:max-w-[500px]">
			<DialogHeader>
				<DialogTitle>Russian Roulette</DialogTitle>
			</DialogHeader>

			<form>
				<div class="space-y-4">
					<div class="flex flex-col gap-4">
						<FormField v-slot="{ value, handleChange }" name="enabled">
							<FormItem class="flex gap-2 space-y-0 items-center">
								<FormLabel>{{ t('sharedTexts.enabled') }}</FormLabel>
								<FormControl>
									<Switch
										:checked="value"
										@update:checked="handleChange"
									/>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<CommandButton name="roulette" />
					</div>

					<Separator />

					<FormField v-slot="{ value, handleChange }" name="canBeUsedByModerator">
						<FormItem class="flex gap-2 space-y-0 items-center">
							<FormLabel>{{ t('games.russianRoulette.canBeUsedByModerator') }}</FormLabel>
							<FormControl>
								<Switch
									:checked="value"
									@update:checked="handleChange"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>

				<Separator />

				<FormField v-slot="{ componentField }" name="timeoutSeconds">
					<FormItem>
						<FormLabel>{{ t('games.russianRoulette.timeoutSeconds') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								type="number"
								:max="86400"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="decisionSeconds">
					<FormItem>
						<FormLabel>{{ t('games.russianRoulette.decisionSeconds') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								type="number"
								:max="60"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<Separator />

				<FormField v-slot="{ componentField }" name="initMessage">
					<FormItem>
						<FormLabel>{{ t('games.russianRoulette.initMessage') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								:maxlength="450"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="surviveMessage">
					<FormItem>
						<FormLabel>{{ t('games.russianRoulette.surviveMessage') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								:maxlength="450"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="deathMessage">
					<FormItem>
						<FormLabel>{{ t('games.russianRoulette.deathMessage') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								:maxlength="450"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="chargedBullets">
					<FormItem>
						<FormLabel>{{ t('games.russianRoulette.chargedBullets') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								type="number"
								:min="1"
								:max="6"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="tumberSize">
					<FormItem>
						<FormLabel>{{ t('games.russianRoulette.tumberSize') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								type="number"
								:min="6"
								:max="12"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

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

					<Button type="submit" @click="save">
						{{ t('sharedButtons.save') }}
					</Button>
				</div>
			</form>
		</DialogContent>
	</Dialog>
</template>
