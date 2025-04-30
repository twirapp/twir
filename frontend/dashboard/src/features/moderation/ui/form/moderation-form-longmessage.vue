<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'

const { t } = useI18n()
</script>

<template>
	<FormField v-slot="{ componentField }" name="triggerLength">
		<FormItem>
			<FormLabel>{{ t('moderation.types.long_message.triggerLength') }}</FormLabel>
			<FormControl>
				<Input type="number" v-bind="componentField" />
			</FormControl>
			<FormMessage />
		</FormItem>
	</FormField>

	<Separator />
</template>
