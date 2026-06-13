<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'

const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col gap-4">
		<div class="flex flex-row flex-wrap w-full gap-2">
			<FormField v-slot="{ componentField }" name="triggerLength">
				<FormItem class="w-full xl:w-[45%]">
					<FormLabel>{{ t('moderation.types.one_man_spam.triggerLength') }}</FormLabel>
					<FormControl>
						<Input type="number" v-bind="componentField" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="oneManSpamMinimumStoredMessages">
				<FormItem class="w-full xl:w-[45%]">
					<FormLabel>{{ t('moderation.types.one_man_spam.minimumStoredMessages') }}</FormLabel>
					<FormControl>
						<Input type="number" v-bind="componentField" min="1" max="20" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>
		</div>

		<div class="flex flex-row w-full flex-wrap gap-2">
			<FormField v-slot="{ componentField }" name="oneManSpamMessageMemorySeconds">
				<FormItem class="w-full xl:w-[45%]">
					<FormLabel>{{ t('moderation.types.one_man_spam.messageMemorySeconds') }}</FormLabel>
					<FormControl>
						<Input type="number" v-bind="componentField" min="1" max="600" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="maxPercentage">
				<FormItem class="w-full xl:w-[45%]">
					<FormLabel>{{ t('moderation.types.one_man_spam.maxPercentage') }}</FormLabel>
					<FormControl>
						<Input type="number" v-bind="componentField" min="1" max="100" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>
		</div>
	</div>
</template>
