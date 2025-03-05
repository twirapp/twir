<script setup lang="ts">
import { GavelIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { formSchema, useVotebanForm } from './composables/use-voteban-form'

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
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'
import { Textarea } from '@/components/ui/textarea'
import CommandButton from '@/features/commands/ui/command-button.vue'
import { VoteBanGameVotingMode } from '@/gql/graphql'

const isDialogOpen = ref(false)
const { t } = useI18n()
const { initialValues, save, settings } = useVotebanForm()

const votebanForm = useForm({
	validationSchema: formSchema,
	initialValues,
	validateOnMount: false,
	keepValuesOnUnmount: true,
})

watch(settings, (newSettings) => {
	if (!newSettings) return
	votebanForm.setValues(toRaw(newSettings.gamesVoteban))
}, { immediate: true })

const onSubmit = votebanForm.handleSubmit(async (values) => {
	await save(values)
	isDialogOpen.value = false
})

function resetSettings() {
	votebanForm.setValues(initialValues)
}
</script>

<template>
	<Dialog v-model:open="isDialogOpen">
		<DialogTrigger asChild>
			<Card
				title="Voteban"
				:icon="GavelIcon"
				:icon-stroke="1"
				:description="t('games.voteban.description')"
			/>
		</DialogTrigger>

		<DialogContent class="sm:max-w-[625px] max-h-[80vh] overflow-y-auto">
			<DialogHeader>
				<DialogTitle>Voteban</DialogTitle>
			</DialogHeader>

			<form @submit="onSubmit">
				<div class="grid gap-4 py-4">
					<div class="flex items-center gap-6">
						<FormField
							v-slot="{ value, handleChange }"
							name="enabled"
						>
							<FormItem class="flex flex-col items-center gap-1">
								<FormLabel>{{ t('sharedTexts.enabled') }}</FormLabel>
								<FormControl>
									<Switch
										:checked="value"
										@update:checked="handleChange"
									/>
								</FormControl>
							</FormItem>
						</FormField>

						<CommandButton name="voteban" />
					</div>

					<Separator />

					<div class="space-y-4">
						<h4 class="font-medium">
							Messages
						</h4>

						<FormField
							v-slot="{ componentField }"
							name="initMessage"
						>
							<FormItem>
								<FormLabel>{{ t('games.voteban.initialMessage') }}</FormLabel>
								<FormControl>
									<Textarea v-bind="componentField" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ componentField }"
							name="banMessage"
						>
							<FormItem>
								<FormLabel>{{ t('games.voteban.banMessage') }}</FormLabel>
								<FormControl>
									<Textarea v-bind="componentField" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ componentField }"
							name="surviveMessage"
						>
							<FormItem>
								<FormLabel>{{ t('games.voteban.surviveMessage') }}</FormLabel>
								<FormControl>
									<Textarea v-bind="componentField" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<Separator />

					<div class="space-y-4">
						<h4 class="font-medium">
							{{ t('sharedTexts.settings') }}
						</h4>

						<FormField
							v-slot="{ componentField }"
							name="votingMode"
						>
							<FormItem>
								<FormLabel>{{ t('games.voteban.voteMode') }}</FormLabel>
								<Select v-bind="componentField">
									<FormControl>
										<SelectTrigger>
											<SelectValue />
										</SelectTrigger>
									</FormControl>
									<SelectContent>
										<SelectItem :value="VoteBanGameVotingMode.Chat">
											Chat
										</SelectItem>
										<SelectItem :value="VoteBanGameVotingMode.Polls" disabled>
											Twitch polls (soon)
										</SelectItem>
									</SelectContent>
								</Select>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ value }"
							name="chatVotesWordsPositive"
						>
							<FormItem>
								<FormControl>
									<FormLabel>{{ t('games.voteban.wordsPositive') }}</FormLabel>
									<TagsInput
										:model-value="value"
										:placeholder="t('games.voteban.wordsPositive')"
										:max="10"
									>
										<TagsInputItem v-for="item in value" :key="item" :value="item">
											<TagsInputItemText />
											<TagsInputItemDelete />
										</TagsInputItem>

										<TagsInputInput placeholder="Enter words..." />
									</TagsInput>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ value }"
							name="chatVotesWordsNegative"
						>
							<FormItem>
								<FormControl>
									<FormLabel>{{ t('games.voteban.wordsNegative') }}</FormLabel>
									<TagsInput
										:model-value="value"
										:placeholder="t('games.voteban.wordsNegative')"
										:max="10"
									>
										<TagsInputItem v-for="item in value" :key="item" :value="item">
											<TagsInputItemText />
											<TagsInputItemDelete />
										</TagsInputItem>

										<TagsInputInput placeholder="Enter words..." />
									</TagsInput>
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ componentField }"
							name="voteDuration"
						>
							<FormItem>
								<FormLabel>{{ t('games.voteban.voteDuration') }}</FormLabel>
								<FormControl>
									<Input type="number" v-bind="componentField" min="1" max="86400" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ componentField }"
							name="neededVotes"
						>
							<FormItem>
								<FormLabel>{{ t('games.voteban.neededVotes') }}</FormLabel>
								<FormControl>
									<Input type="number" v-bind="componentField" min="1" max="999999" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ componentField }"
							name="timeoutSeconds"
						>
							<FormItem>
								<FormLabel>{{ t('games.voteban.banDuration') }}</FormLabel>
								<FormControl>
									<Input type="number" v-bind="componentField" min="1" max="86400" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
					</div>

					<Separator />

					<div class="space-y-4">
						<h4 class="font-medium">
							Moderators
						</h4>

						<FormField
							v-slot="{ value, handleChange }"
							name="timeoutModerators"
						>
							<FormItem class="flex items-center justify-between">
								<FormLabel>{{ t('games.voteban.timeoutModerators') }}</FormLabel>
								<FormControl>
									<Switch
										:checked="value"
										@update:checked="handleChange"
									/>
								</FormControl>
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ componentField }"
							name="banMessageModerators"
						>
							<FormItem>
								<FormLabel>{{ t('games.voteban.banMessageModerators') }}</FormLabel>
								<FormControl>
									<Textarea v-bind="componentField" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>

						<FormField
							v-slot="{ componentField }"
							name="surviveMessageModerators"
						>
							<FormItem>
								<FormLabel>{{ t('games.voteban.surviveMessageModerators') }}</FormLabel>
								<FormControl>
									<Textarea v-bind="componentField" />
								</FormControl>
								<FormMessage />
							</FormItem>
						</FormField>
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

						<Button type="submit">
							{{ t('sharedButtons.save') }}
						</Button>
					</div>
				</div>
			</form>
		</DialogContent>
	</Dialog>
</template>
