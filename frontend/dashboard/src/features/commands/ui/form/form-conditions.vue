<script setup lang="ts">
import { computed } from 'vue'
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
	]
})
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>Conditions</CardTitle>
		</CardHeader>
		<CardContent class="flex flex-col gap-2">
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
							<Checkbox :checked="value" @update:checked="handleChange" />
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
