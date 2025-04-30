<script setup lang="ts">
import { AlertTriangle, CaseSensitiveIcon, Plus, RegexIcon, WholeWordIcon, X } from 'lucide-vue-next'
import { FieldArray, useField } from 'vee-validate'
import { useI18n } from 'vue-i18n'

import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'

const { t } = useI18n()

const { value: denyList, setValue: setDenyList } = useField<string[]>('denyList')

function addItem() {
	setDenyList([...(denyList.value || []), ''])
}
</script>

<template>
	<div class="flex flex-col gap-4">
		<Separator />

		<FormField v-slot="{ field }" name="denyListRegexpEnabled">
			<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
				<div class="space-y-0.5">
					<FormLabel class="flex gap-2 items-center text-base">
						<RegexIcon />
						Regexp
					</FormLabel>
					<FormDescription>
						<i18n-t keypath="moderation.types.deny_list.regexp">
							<a
								href="https://yourbasic.org/golang/regexp-cheat-sheet/#cheat-sheet"
								target="_blank"
								class="text-primary underline-offset-4 underline"
							>
								{{ t('moderation.types.deny_list.regexpCheatSheet') }}
							</a>
						</i18n-t>
					</FormDescription>
				</div>
				<FormControl>
					<Switch
						:checked="field.value"
						default-checked
						@update:checked="field['onUpdate:modelValue']"
					/>
				</FormControl>
			</FormItem>
		</FormField>

		<FormField v-slot="{ field }" name="denyListWordBoundaryEnabled">
			<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
				<div class="space-y-0.5">
					<FormLabel class="flex gap-2 items-center text-base">
						<WholeWordIcon />
						{{ t('moderation.types.deny_list.wordBoundary.label') }}
					</FormLabel>
					<FormDescription>
						{{ t('moderation.types.deny_list.wordBoundary.description') }}
					</FormDescription>
				</div>
				<FormControl>
					<Switch
						:checked="field.value"
						default-checked
						@update:checked="field['onUpdate:modelValue']"
					/>
				</FormControl>
			</FormItem>
		</FormField>

		<FormField v-slot="{ field }" name="denyListSensitivityEnabled">
			<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
				<div class="space-y-0.5">
					<FormLabel class="flex gap-2 items-center text-base">
						<CaseSensitiveIcon />
						{{ t('moderation.types.deny_list.caseSensitive.label') }}
					</FormLabel>
					<FormDescription>
						{{ t('moderation.types.deny_list.caseSensitive.description') }}
					</FormDescription>
				</div>
				<FormControl>
					<Switch
						:checked="field.value"
						default-checked
						@update:checked="field['onUpdate:modelValue']"
					/>
				</FormControl>
				<FormMessage />
			</FormItem>
		</FormField>

		<FieldArray v-slot="{ fields, remove }" name="denyList">
			<Separator />

			<Alert
				v-if="!denyList?.length"
				class="flex items-center gap-2"
			>
				<AlertTriangle class="h-4 w-4" />
				<AlertDescription>
					{{ t('moderation.types.deny_list.empty') }}
				</AlertDescription>
			</Alert>

			<div v-for="(field, index) in fields" :key="`word-${field.key}`">
				<FormField v-slot="{ componentField }" :name="`denyList[${index}]`">
					<FormItem class="flex gap-2">
						<FormControl>
							<Textarea
								:model-value="componentField.modelValue"
								:maxlength="500"
								rows="1"
								class="flex-1"
								placeholder="Word"
								@update:model-value="componentField.onChange"
							/>
							<Button
								variant="ghost"
								size="icon"
								class="shrink-0"
								type="button"
								@click="remove(index)"
							>
								<X class="h-4 w-4" />
							</Button>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<Button
				:disabled="(denyList?.length || 0) >= 100"
				variant="outline"
				class="w-full"
				type="button"
				@click="addItem"
			>
				<Plus class="h-4 w-4 mr-2" />
				{{ t('sharedButtons.create') }}
			</Button>
		</FieldArray>
	</div>
</template>
