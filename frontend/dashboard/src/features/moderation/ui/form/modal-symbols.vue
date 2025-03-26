<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'

const { t } = useI18n()
</script>

<template>
	<div>
		<FormField v-slot="{ componentField }" name="triggerLength">
			<FormItem>
				<FormLabel>{{ t('moderation.types.symbols.triggerLength') }}</FormLabel>
				<FormControl>
					<Input type="number" v-bind="componentField" />
				</FormControl>
				<FormMessage />
			</FormItem>
		</FormField>

		<FormField v-slot="{ componentField }" name="maxPercentage">
			<FormItem>
				<FormLabel>{{ t('moderation.types.symbols.maxPercentage') }}</FormLabel>
				<FormControl>
					<Input type="number" v-bind="componentField" />
				</FormControl>
				<FormMessage />
			</FormItem>
		</FormField>
	</div>
</template>
