<script setup lang="ts">
import { type Font, FontSelector } from '~/lib/fontsource'
import { computed, ref, watch } from 'vue'

import { toast } from 'vue-sonner'

import { useNowPlayingForm } from './use-now-playing-form'

import { useProfile, useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useNowPlayingOverlayApi } from '#layers/dashboard/api/overlays/now-playing'
import { useCopyOverlayLink } from '#layers/dashboard/components/overlays/copyOverlayLink'









import { ChannelRolePermissionEnum } from '~/gql/graphql'

const { t } = useI18n()
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

	toast.success(t('sharedTexts.saved'), {
		duration: 1500,
	})
}

const fontData = ref<Font | null>(null)
watch(
	() => fontData.value,
	(font) => {
		if (!font) return
		formValue.value.fontFamily = font.id
	},
	{ deep: true }
)

const fontWeightOptions = computed(() => {
	if (!fontData.value) return []
	return fontData.value.weights.map((weight) => ({ label: `${weight}`, value: weight }))
})
</script>

<template>
	<UiCard v-if="formValue" class="card">
		<UiCardContent class="pt-4 flex flex-col gap-4">
			<div class="flex flex-col gap-2">
				<UiLabel for="preset">Style</UiLabel>
				<UiSelect id="preset" v-model:model-value="formValue.preset" default-value="AIDEN_REDESIGN">
					<UiSelectTrigger class="w-45">
						<UiSelectValue placeholder="Select a preset" />
					</UiSelectTrigger>
					<UiSelectContent>
						<UiSelectGroup>
							<UiSelectItem value="AIDEN_REDESIGN"> Aiden Redesign </UiSelectItem>
							<UiSelectItem value="TRANSPARENT"> Transparent </UiSelectItem>
							<UiSelectItem value="SIMPLE_LINE"> Simple line </UiSelectItem>
						</UiSelectGroup>
					</UiSelectContent>
				</UiSelect>
			</div>

			<div class="flex flex-col gap-2">
				<UiLabel for="showImage">Show image</UiLabel>
				<UiSwitch
					id="showImage"
					:model-value="formValue.showImage"
					@update:model-value="formValue.showImage = $event"
				/>
			</div>

			<div class="flex flex-col gap-2">
				<UiLabel for="backgroundColor">Background color</UiLabel>
				<UiColorPicker v-model="formValue.backgroundColor" />
			</div>

			<div class="flex flex-col gap-2">
				<UiLabel for="fontFamily">{{ t('overlays.chat.fontFamily') }}</UiLabel>
				<FontSelector
					id="fontFamily"
					v-model:font="fontData"
					:font-family="formValue.fontFamily"
					:font-weight="formValue.fontWeight"
					font-style="normal"
				/>
			</div>

			<div class="flex flex-col gap-2">
				<UiLabel for="fontWeight">{{ t('overlays.chat.fontWeight') }}</UiLabel>

				<UiPopover>
					<UiPopoverTrigger as-child>
						<UiButton variant="outline" size="sm" class="w-37.5 justify-start">
							<template v-if="formValue.fontWeight">
								{{ formValue.fontWeight }}
							</template>
							<template v-else> + Set font weight </template>
						</UiButton>
					</UiPopoverTrigger>
					<UiPopoverContent class="p-0" side="right" align="start">
						<UiCommand>
							<UiCommandList>
								<UiCommandGroup>
									<UiCommandItem
										v-for="weight in fontWeightOptions"
										:key="weight.value"
										:value="weight.value"
										@select="
											() => {
												formValue.fontWeight = weight.value
											}
										"
									>
										{{ weight.label }}
									</UiCommandItem>
								</UiCommandGroup>
							</UiCommandList>
						</UiCommand>
					</UiPopoverContent>
				</UiPopover>
			</div>

			<div class="flex flex-col gap-2">
				<UiLabel for="hideTimeout">{{ t('overlays.chat.hideTimeout') }}</UiLabel>
				<UiInput
					id="hideTimeout"
					v-model:model-value="formValue.hideTimeout"
					type="number"
					:min="0"
					:max="600"
				/>
			</div>
		</UiCardContent>

		<UiCardFooter class="flex justify-end gap-2">
			<UiButton variant="destructive" @click="deleter.executeMutation({ id: formValue.id! })">
				{{ t('sharedButtons.delete') }}
			</UiButton>
			<UiButton
				:disabled="!formValue.id || !canCopyLink"
				variant="secondary"
				@click="copyOverlayLink({ id: formValue.id! })"
			>
				{{ t('overlays.copyOverlayLink') }}
			</UiButton>
			<UiButton @click="save">
				{{ t('sharedButtons.save') }}
			</UiButton>
		</UiCardFooter>
	</UiCard>
</template>
