<script setup lang="ts">
import {
	AlertTriangle,
	CaseSensitiveIcon,
	Plus,
	RegexIcon,
	WholeWordIcon,
	X,
} from 'lucide-vue-next'
import { FieldArray, useField } from 'vee-validate'









const { t } = useI18n()

const { value: denyList, setValue: setDenyList } = useField<string[]>('denyList')
const { value: isRegexpEnabled } = useField<boolean>('denyListRegexpEnabled')

function addItem() {
	setDenyList([...(denyList.value || []), ''])
}
</script>

<template>
	<div class="flex flex-col gap-4">
		<UiSeparator />

		<UiFormField v-slot="{ field }" name="denyListRegexpEnabled">
			<UiFormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
				<div class="space-y-0.5">
					<UiFormLabel class="flex gap-2 items-center text-base">
						<RegexIcon />
						Regexp
					</UiFormLabel>
					<UiFormDescription>
						<i18n-t keypath="moderation.types.deny_list.regexp">
							<a
								href="https://yourbasic.org/golang/regexp-cheat-sheet/#cheat-sheet"
								target="_blank"
								class="text-primary underline-offset-4 underline"
							>
								{{ t('moderation.types.deny_list.regexpCheatSheet') }}
							</a>
						</i18n-t>
					</UiFormDescription>
				</div>
				<UiFormControl>
					<UiSwitch
						:model-value="field.value"
						default-checked
						@update:model-value="field['onUpdate:modelValue']"
					/>
				</UiFormControl>
			</UiFormItem>
		</UiFormField>

		<UiFormField v-slot="{ field }" name="denyListWordBoundaryEnabled">
			<UiFormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
				<div class="space-y-0.5">
					<UiFormLabel class="flex gap-2 items-center text-base">
						<WholeWordIcon />
						{{ t('moderation.types.deny_list.wordBoundary.label') }}
					</UiFormLabel>
					<UiFormDescription class="flex flex-col">
						<span>
							{{ t('moderation.types.deny_list.wordBoundary.description') }}
						</span>
						<span class="text-xs text-orange-600">
							Word boundary cannot be used within regular expressions.
						</span>
					</UiFormDescription>
				</div>
				<UiFormControl>
					<UiSwitch
						:disabled="isRegexpEnabled"
						:model-value="field.value"
						default-checked
						@update:model-value="field['onUpdate:modelValue']"
					/>
				</UiFormControl>
			</UiFormItem>
		</UiFormField>

		<UiFormField v-slot="{ field }" name="denyListSensitivityEnabled">
			<UiFormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
				<div class="space-y-0.5">
					<UiFormLabel class="flex gap-2 items-center text-base">
						<CaseSensitiveIcon />
						{{ t('moderation.types.deny_list.caseSensitive.label') }}
					</UiFormLabel>
					<UiFormDescription>
						{{ t('moderation.types.deny_list.caseSensitive.description') }}
					</UiFormDescription>
				</div>
				<UiFormControl>
					<UiSwitch
						:model-value="field.value"
						default-checked
						@update:model-value="field['onUpdate:modelValue']"
					/>
				</UiFormControl>
				<UiFormMessage />
			</UiFormItem>
		</UiFormField>

		<FieldArray v-slot="{ fields, remove }" name="denyList">
			<UiSeparator />

			<UiAlert v-if="!denyList?.length" class="flex items-center gap-2">
				<AlertTriangle class="h-4 w-4" />
				<UiAlertDescription>
					{{ t('moderation.types.deny_list.empty') }}
				</UiAlertDescription>
			</UiAlert>

			<div v-for="(field, index) in fields" :key="`word-${field.key}`">
				<UiFormField v-slot="{ componentField }" :name="`denyList[${index}]`">
					<UiFormItem class="flex gap-2">
						<UiFormControl>
							<UiTextarea
								:model-value="componentField.modelValue"
								:maxlength="500"
								rows="1"
								class="flex-1"
								placeholder="Word"
								@update:model-value="componentField.onChange"
							/>
							<UiButton
								variant="ghost"
								size="icon"
								class="shrink-0"
								type="button"
								@click="remove(index)"
							>
								<X class="h-4 w-4" />
							</UiButton>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<UiButton
				:disabled="(denyList?.length || 0) >= 100"
				variant="outline"
				class="w-full"
				type="button"
				@click="addItem"
			>
				<Plus class="h-4 w-4 mr-2" />
				{{ t('sharedButtons.create') }}
			</UiButton>
		</FieldArray>
	</div>
</template>
