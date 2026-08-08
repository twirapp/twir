<script setup lang="ts">
import { useForm } from 'vee-validate'
import { ref } from 'vue'
import { toast } from 'vue-sonner'
import { z } from 'zod'

import type { WidgetItem } from '~~/layers/dashboard/components/dashboard/widgets.js'

import {
	useDashboardWidgetsDelete,
	useDashboardWidgetsUpdateCustom,
} from '~~/layers/dashboard/api/dashboard-widgets-layout.js'
import { Button } from '@/components/ui/button'
import { CardContent } from '@/components/ui/card'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'

import ActionConfirm from '@/components/ui/action-confirm'

import Card from './card.vue'

defineOptions({
	inheritAttrs: false,
})

interface Props {
	item: WidgetItem
	url: string
}

const props = defineProps<Props>()

const isEditDialogOpen = ref(false)
const isDeleteConfirmOpen = ref(false)
const updateMutation = useDashboardWidgetsUpdateCustom()
const deleteMutation = useDashboardWidgetsDelete()

// The widgetId is already the full identifier (e.g., "custom-{uuid}")
const widgetId = props.item.i.toString()

const formSchema = z.object({
	name: z.string().min(2, 'Name must be at least 2 characters.'),
	url: z.string().url('Must be a valid URL'),
})

const { handleSubmit, setValues } = useForm({
	validationSchema: formSchema,
	initialValues: {
		name: props.item.displayName || '',
		url: props.url,
	},
})

const onEdit = () => {
	setValues({
		name: props.item.displayName || '',
		url: props.url,
	})
	isEditDialogOpen.value = true
}

const onSubmitEdit = handleSubmit(async (values) => {
	const result = await updateMutation.executeMutation({
		input: {
			widgetId: widgetId,
			name: values.name,
			url: values.url,
		},
	})

	if (result.error) {
		toast.error('Failed to update widget', {
			description: result.error.message,
		})
	} else {
		toast.success('Widget updated successfully')
		isEditDialogOpen.value = false
	}
})

const onDelete = () => {
	isDeleteConfirmOpen.value = true
}

const confirmDelete = async () => {
	const result = await deleteMutation.executeMutation({ widgetId: widgetId })

	if (result.error) {
		toast.error('Failed to delete widget', {
			description: result.error.message,
		})
	} else {
		toast.success('Widget deleted successfully')
	}
}
</script>

<template>
	<div v-bind="$attrs">
		<Card
			:item="item"
			class="h-full"
		>
			<template #header-extra>
				<Button
					size="sm"
					variant="ghost"
					@click="onEdit"
				>
					<Icon
						name="lucide:pencil"
						class="h-4 w-4"
					/>
				</Button>
				<Button
					size="sm"
					variant="ghost"
					@click="onDelete"
				>
					<Icon
						name="lucide:trash"
						class="h-4 w-4"
					/>
				</Button>
			</template>

			<CardContent class="flex-1 p-0">
				<iframe
					:src="url"
					class="h-full w-full border-0"
					sandbox="allow-scripts allow-same-origin"
					loading="lazy"
				/>
			</CardContent>
		</Card>

		<!-- Edit Dialog -->
		<Dialog v-model:open="isEditDialogOpen">
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Edit Custom Widget</DialogTitle>
					<DialogDescription> Update the name and URL of your custom widget. </DialogDescription>
				</DialogHeader>

				<form
					@submit="onSubmitEdit"
					class="space-y-4"
				>
					<FormField
						v-slot="{ componentField }"
						name="name"
					>
						<FormItem>
							<FormLabel>Widget Name</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
									placeholder="My Custom Widget"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField
						v-slot="{ componentField }"
						name="url"
					>
						<FormItem>
							<FormLabel>Website URL</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
									placeholder="https://example.com"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<div class="flex justify-end gap-2">
						<Button
							type="button"
							variant="outline"
							@click="isEditDialogOpen = false"
						>
							Cancel
						</Button>
						<Button
							type="submit"
							:disabled="updateMutation.fetching.value"
						>
							{{ updateMutation.fetching.value ? 'Updating...' : 'Update Widget' }}
						</Button>
					</div>
				</form>
			</DialogContent>
		</Dialog>

		<ActionConfirm
			v-model:open="isDeleteConfirmOpen"
			:confirm-text="`Are you sure you want to delete widget &quot;${props.item.displayName}&quot;?`"
			@confirm="confirmDelete"
		/>
	</div>
</template>
