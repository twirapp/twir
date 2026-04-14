<script setup lang="ts">
import { SlidersHorizontalIcon } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import TwitchCategorySearchShadcnMultiple from '@/components/twitch-category-search-shadcn-multiple.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'

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
	<Card>
		<CardHeader class="flex flex-row place-content-center flex-wrap">
			<CardTitle class="flex items-center gap-2">
				<SlidersHorizontalIcon />
				Conditions
			</CardTitle>
		</CardHeader>
		<CardContent class="flex flex-col gap-2 pt-4">
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-2">
				<FormField
					v-for="checkbox of checkboxes"
					:key="checkbox.name"
					v-slot="{ value, handleChange }"
					type="checkbox"
					:name="checkbox.name"
				>
					<FormItem class="flex flex-row items-start gap-x-3 space-y-0 rounded-md border p-4">
						<FormControl>
							<Checkbox :model-value="value" @update:model-value="handleChange" />
						</FormControl>
						<div class="space-y-1 leading-none">
							<FormLabel>{{ checkbox.label }}</FormLabel>
							<FormDescription v-if="checkbox.description">
								{{ checkbox.description }}
							</FormDescription>
							<FormMessage />
						</div>
					</FormItem>
				</FormField>
			</div>

			<div>
				<FormField v-slot="{ field }" name="enabledCategories">
					<FormItem>
						<FormLabel>
							{{ t('commands.modal.gameCategories.label') }}
						</FormLabel>
						<FormControl>
							<TwitchCategorySearchShadcnMultiple
								:id="field.name"
								:model-value="field.value"
								@update:model-value="field['onUpdate:modelValue']"
							/>
						</FormControl>
					</FormItem>
				</FormField>
			</div>
		</CardContent>
	</Card>
</template>
