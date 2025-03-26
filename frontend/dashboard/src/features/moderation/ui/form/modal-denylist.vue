<script setup lang="ts">
import { AlertTriangle, Info, Plus, X } from 'lucide-vue-next'
import { FieldArray, useField } from 'vee-validate'
import { useI18n } from 'vue-i18n'

import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Separator } from '@/components/ui/separator'
import { Textarea } from '@/components/ui/textarea'

const { t } = useI18n()

const { value: denyList, setValue: setDenyList } = useField<string[]>('denyList')

function addItem() {
	setDenyList([...(denyList.value || []), ''])
}
</script>

<template>
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

		<Alert
			v-else
			class="flex items-center gap-2"
		>
			<Info class="h-4 w-4" />
			<AlertDescription>
				<i18n-t keypath="moderation.types.deny_list.regexp">
					<a
						href="https://yourbasic.org/golang/regexp-cheat-sheet/#cheat-sheet"
						target="_blank"
						class="text-primary underline-offset-4 underline"
					>
						{{ t('moderation.types.deny_list.regexpCheatSheet') }}
					</a>
				</i18n-t>
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
</template>
