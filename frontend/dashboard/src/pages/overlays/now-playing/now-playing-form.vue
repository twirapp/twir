<script setup lang="ts">
import { type Font, FontSelector } from '@twir/fontsource'
import { NColorPicker } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useNowPlayingForm } from './use-now-playing-form'

import {
	useNowPlayingOverlayApi,
	useProfile,
	useUserAccessFlagChecker,
} from '@/api'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import {
	Command,
	CommandGroup,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const { t } = useI18n()

const discrete = useNaiveDiscrete()
const { copyOverlayLink } = useCopyOverlayLink('now-playing')
const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)

const { data: profile } = useProfile()
const { data: formValue } = useNowPlayingForm()

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays
})

const manager = useNowPlayingOverlayApi()
const updater = manager.useNowPlayingUpdate()
const deleter = manager.useNowPlayingDelete()

async function save() {
	if (!formValue.value?.id) return

	await updater.executeMutation({
		id: formValue.value.id,
		input: {
			preset: formValue.value.preset,
			fontFamily: formValue.value.fontFamily,
			fontWeight: formValue.value.fontWeight,
			backgroundColor: formValue.value.backgroundColor,
			showImage: formValue.value.showImage,
			hideTimeout: formValue.value.hideTimeout,
		},
	})

	discrete.notification.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	})
}

const fontData = ref<Font | null>(null)
watch(() => fontData.value, (font) => {
	if (!font) return
	formValue.value.fontFamily = font.id
}, { deep: true })

const fontWeightOptions = computed(() => {
	if (!fontData.value) return []
	return fontData.value.weights.map((weight) => ({ label: `${weight}`, value: weight }))
})
</script>

<template>
	<Card v-if="formValue" class="card">
		<CardContent class="pt-4 flex flex-col gap-4">
			<div class="flex flex-col gap-2">
				<Label for="preset">Style</Label>
				<Select id="preset" v-model:model-value="formValue.preset" default-value="AIDEN_REDESIGN">
					<SelectTrigger class="w-[180px]">
						<SelectValue placeholder="Select a preset" />
					</SelectTrigger>
					<SelectContent>
						<SelectGroup>
							<SelectItem value="AIDEN_REDESIGN">
								Aiden Redesign
							</SelectItem>
							<SelectItem value="TRANSPARENT">
								Transparent
							</SelectItem>
							<SelectItem value="SIMPLE_LINE">
								Simple line
							</SelectItem>
						</SelectGroup>
					</SelectContent>
				</Select>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="showImage">Show image</Label>
				<Switch
					id="showImage"
					v-model:checked="formValue.showImage"
					@update:checked="formValue.showImage = $event"
				/>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="backgroundColor">Background color</Label>
				<NColorPicker
					v-model:value="formValue.backgroundColor"
				/>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="fontFamily">{{ t('overlays.chat.fontFamily') }}</Label>
				<FontSelector
					id="fontFamily"
					v-model:font="fontData"
					:font-family="formValue.fontFamily"
					:font-weight="formValue.fontWeight"
					font-style="normal"
				/>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="fontWeight">{{ t('overlays.chat.fontWeight') }}</Label>

				<Popover>
					<PopoverTrigger as-child>
						<Button
							variant="outline"
							size="sm"
							class="w-[150px] justify-start"
						>
							<template v-if="formValue.fontWeight">
								{{ formValue.fontWeight }}
							</template>
							<template v-else>
								+ Set font weight
							</template>
						</Button>
					</PopoverTrigger>
					<PopoverContent class="p-0" side="right" align="start">
						<Command>
							<CommandList>
								<CommandGroup>
									<CommandItem
										v-for="weight in fontWeightOptions"
										:key="weight.value"
										:value="weight.value"
										@select="() => {
											formValue.fontWeight = weight.value
										}"
									>
										{{ weight.label }}
									</CommandItem>
								</CommandGroup>
							</CommandList>
						</Command>
					</PopoverContent>
				</Popover>
			</div>

			<div class="flex flex-col gap-2">
				<Label for="hideTimeout">{{ t('overlays.chat.hideTimeout') }}</Label>
				<Input
					id="hideTimeout"
					v-model:model-value="formValue.hideTimeout"
					type="number"
					:min="0"
					:max="600"
				/>
			</div>
		</CardContent>

		<CardFooter class="flex justify-end gap-2">
			<Button variant="destructive" @click="deleter.executeMutation({ id: formValue.id! })">
				{{ t('sharedButtons.delete') }}
			</Button>
			<Button
				:disabled="!formValue.id || !canCopyLink"
				variant="secondary"
				@click="copyOverlayLink({ id: formValue.id! })"
			>
				{{ t('overlays.copyOverlayLink') }}
			</Button>
			<Button @click="save">
				{{ t('sharedButtons.save') }}
			</Button>
		</CardFooter>
	</Card>
</template>
