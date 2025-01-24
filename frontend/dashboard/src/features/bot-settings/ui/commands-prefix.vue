<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { z } from 'zod'

import { useCommandsPrefixApi } from '@/api/commands-prefix'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useToast } from '@/components/ui/toast'

const { toast } = useToast()
const { t } = useI18n()

const api = useCommandsPrefixApi()
const { data: currentPrefix } = api.usePrefix()
const update = api.usePrefixUpdate()
const reset = api.usePrefixReset()

const formSchema = toTypedSchema(z.object({
	prefix: z.string().min(1).max(10),
}))

const form = useForm({
	validationSchema: formSchema,
})

watch(currentPrefix, (v) => {
	if (!v?.channelsCommandsPrefix) return

	form.setFieldValue('prefix', v.channelsCommandsPrefix)
})

const onSubmit = form.handleSubmit(async (values) => {
	try {
		await update.executeMutation({
			input: { newPrefix: values.prefix },
		})
		toast({
			title: 'Updated',
		})
	} catch {
		toast({
			title: 'Error happend on update',
			variant: 'destructive',
		})
	}
})

async function onReset() {
	try {
		await reset.executeMutation({})
		toast({
			title: 'Resetted',
		})
	} catch {
		toast({
			title: 'Error happend on reset',
			variant: 'destructive',
		})
	}
}
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
			Commands prefix
		</h4>
		<Card>
			<form @submit="onSubmit">
				<CardContent class="pt-6 flex gap-2">
					<FormField v-slot="{ componentField }" name="prefix">
						<FormItem class="w-full">
							<FormLabel>Prefix used for all commands, "!" by default</FormLabel>
							<FormControl>
								<Input type="text" placeholder="Global commands prefix" v-bind="componentField" class="w-full" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</CardContent>
				<CardFooter class="flex flex-row gap-2 justify-end">
					<Button type="button" variant="destructive" class="place-self-end" @click="onReset">
						Reset
					</Button>
					<Button type="submit">
						{{ t('sharedButtons.save') }}
					</Button>
				</CardFooter>
			</form>
		</Card>
	</div>
</template>
