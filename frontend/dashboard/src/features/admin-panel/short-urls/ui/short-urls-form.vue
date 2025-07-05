<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader,CardTitle } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
	useAdminShortUrlsApi,
} from '@/features/admin-panel/short-urls/composables/use-admin-short-urls-api.ts'

const api = useAdminShortUrlsApi()

const formSchema = z.object({
	link: z.string().url().min(1).max(2000).trim(),
	shortId: z.string().min(1).max(10).trim().optional(),
})

const urlForm = useForm({
	validationSchema: toTypedSchema(formSchema),
	initialValues: {
		shortId: undefined,
		link: '',
	},
})

const handleSubmit = urlForm.handleSubmit(async (values) => {
	try {
		await api.createShortUrl(values.link, values.shortId)
		urlForm.resetForm()
	} catch (err) {
		console.error(err)
	}
})
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>
				Create short url
			</CardTitle>
		</CardHeader>
		<CardContent>
			<form class="flex flex-col gap-2" @submit.prevent="handleSubmit">
				<FormField v-slot="{ componentField }" name="link">
					<FormItem>
						<FormLabel>Link</FormLabel>
						<FormControl>
							<Input type="text" v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="shortId">
					<FormItem>
						<FormLabel>Short ID</FormLabel>
						<FormControl>
							<Input type="text" v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<Button type="submit" class="place-self-end mt-2">
					Create
				</Button>
			</form>
		</CardContent>
	</Card>
</template>
