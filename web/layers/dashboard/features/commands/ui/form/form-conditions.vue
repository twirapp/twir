<script setup lang="ts">
import { SlidersHorizontalIcon } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed, watch } from 'vue'


import TwitchCategorySearchShadcnMultiple from '#layers/dashboard/components/twitch-category-search-shadcn-multiple.vue'




const { value: onlineOnly, setValue: setOnlineOnly } = useField<boolean>('onlineOnly')
const { value: offlineOnly, setValue: setOfflineOnly } = useField<boolean>('offlineOnly')

watch(onlineOnly, (v) => {
	if (v) {
		setOfflineOnly(false)
	}
})

watch(offlineOnly, (v) => {
	if (v) {
		setOnlineOnly(false)
	}
})

const { t } = useI18n()
const checkboxes = computed(() => {
	return [
		{
			name: 'isReply',
			label: t('sharedTexts.reply.label'),
			description: t('sharedTexts.reply.text'),
		},
		{
			name: 'visible',
			label: t('commands.modal.settings.visible.label'),
			description: t('commands.modal.settings.visible.text'),
		},
		{
			name: 'keepResponsesOrder',
			label: t('commands.modal.settings.keepOrder.label'),
			description: t('commands.modal.settings.keepOrder.text'),
		},
		{
			name: 'onlineOnly',
			label: t('commands.modal.settings.onlineOnly.label'),
			description: t('commands.modal.settings.onlineOnly.text'),
		},
		{
			name: 'offlineOnly',
			label: t('commands.modal.settings.offlineOnly.label'),
			description: t('commands.modal.settings.offlineOnly.text'),
		},
	]
})
</script>

<template>
	<UiCard>
		<UiCardHeader class="flex flex-row place-content-center flex-wrap">
			<UiCardTitle class="flex items-center gap-2">
				<SlidersHorizontalIcon />
				Conditions
			</UiCardTitle>
		</UiCardHeader>
		<UiCardContent class="flex flex-col gap-2 pt-4">
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
				<UiFormField
					v-for="checkbox of checkboxes"
					:key="checkbox.name"
					v-slot="{ value, handleChange }"
					type="checkbox"
					:name="checkbox.name"
				>
					<UiFormItem class="flex flex-row items-start gap-x-3 space-y-0 rounded-md border p-4">
						<UiFormControl>
							<UiCheckbox :model-value="value" @update:model-value="handleChange" />
						</UiFormControl>
						<div class="space-y-1 leading-none">
							<UiFormLabel>{{ checkbox.label }}</UiFormLabel>
							<UiFormDescription v-if="checkbox.description">
								{{ checkbox.description }}
							</UiFormDescription>
							<UiFormMessage />
						</div>
					</UiFormItem>
				</UiFormField>
			</div>

			<div>
				<UiFormField v-slot="{ field }" name="enabledCategories">
					<UiFormItem>
						<UiFormLabel>
							{{ t('commands.modal.gameCategories.label') }}
						</UiFormLabel>
						<UiFormControl>
							<TwitchCategorySearchShadcnMultiple
								:id="field.name"
								:model-value="field.value"
								@update:model-value="field['onUpdate:modelValue']"
							/>
						</UiFormControl>
					</UiFormItem>
				</UiFormField>
			</div>
		</UiCardContent>
	</UiCard>
</template>
