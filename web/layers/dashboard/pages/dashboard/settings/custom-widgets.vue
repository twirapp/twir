<script setup lang="ts">
import { useForm } from 'vee-validate'
import { computed, ref } from 'vue'
import { toast } from 'vue-sonner'
import { z } from 'zod'
import {
	useDashboardWidgetsCreateCustom,
	useDashboardWidgetsDelete,
	useDashboardWidgetsLayout,
} from '~~/layers/dashboard/api/dashboard-widgets-layout.js'

import { Button } from '@/components/ui/button'
import ActionConfirm from '@/components/ui/action-confirm'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { layout, fetching } = useDashboardWidgetsLayout()
const createMutation = useDashboardWidgetsCreateCustom()
const deleteMutation = useDashboardWidgetsDelete()

const customWidgets = computed(() => {
	return layout.value
		.filter((w) => w.type === 'CUSTOM')
		.map((w) => ({
			id: w.widgetId,
			name: w.customName || 'Unnamed Widget',
			url: w.customUrl || '',
		}))
})

const isDialogOpen = ref(false)
const widgetPendingDeletion = ref<{ id: string; name: string } | null>(null)

const formSchema = z.object({
	name: z.string().min(2, 'Name must be at least 2 characters.'),
	url: z.string().url('Must be a valid URL'),
})

const { handleSubmit, resetForm } = useForm({
	validationSchema: formSchema,
})

const onSubmit = handleSubmit(async (values) => {
	const result = await createMutation.executeMutation({
		input: {
			name: values.name,
			url: values.url,
			x: 0,
			y: 0,
			w: 4,
			h: 8,
		},
	})

	if (result.error) {
		toast.error('Failed to create widget', {
			description: result.error.message,
		})
	} else {
		toast.success('Widget created successfully')
		resetForm()
		isDialogOpen.value = false
	}
})

async function deleteWidget() {
	if (!widgetPendingDeletion.value) {
		return
	}

	const result = await deleteMutation.executeMutation({ widgetId: widgetPendingDeletion.value.id })

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
	<div class="p-6">
		<h1 class="mb-6 text-3xl font-bold">Custom Dashboard Widgets</h1>

		<div class="max-w-2xl space-y-6">
			<div class="bg-card rounded-lg border p-6">
				<h2 class="mb-4 text-xl font-semibold">Your Custom Widgets</h2>

				<div
					v-if="fetching"
					class="text-muted-foreground py-8 text-center"
				>
					Loading...
				</div>

				<div
					v-else-if="customWidgets.length === 0"
					class="text-muted-foreground py-8 text-center"
				>
					No custom widgets yet. Click the + button to create one!
				</div>

				<div
					v-else
					class="space-y-3"
				>
					<div
						v-for="widget in customWidgets"
						:key="widget.id"
						class="bg-secondary/50 flex items-center justify-between rounded-lg p-4"
					>
						<div class="flex-1">
							<h3 class="font-medium">{{ widget.name }}</h3>
							<p class="text-muted-foreground truncate text-sm">
								{{ widget.url }}
							</p>
						</div>
						<Button
							variant="destructive"
							size="icon"
							@click="widgetPendingDeletion = { id: widget.id, name: widget.name }"
							:disabled="deleteMutation.fetching.value"
						>
							<Icon
								name="lucide:trash"
								class="h-4 w-4"
							/>
						</Button>
					</div>
				</div>
			</div>

			<div class="rounded-lg border border-blue-500/20 bg-blue-500/10 p-4">
				<h3 class="mb-2 font-semibold text-blue-500">How to use:</h3>
				<ol class="text-muted-foreground list-inside list-decimal space-y-1 text-sm">
					<li>Create a custom widget by clicking the + button</li>
					<li>Go to your dashboard page</li>
					<li>Click the "+" button in the bottom right</li>
					<li>Select your custom widget to add it to the dashboard</li>
					<li>Drag and resize it like any other widget</li>
				</ol>
			</div>
		</div>
	</div>

	<Dialog v-model:open="isDialogOpen">
		<DialogTrigger as-child>
			<Button
				size="icon"
				class="fixed right-8 bottom-8 z-50 h-14 w-14 rounded-full shadow-lg"
			>
				<Icon
					name="lucide:plus"
					class="h-6 w-6"
				/>
			</Button>
		</DialogTrigger>
		<DialogContent>
			<DialogHeader>
				<DialogTitle>Create Custom Widget</DialogTitle>
				<DialogDescription>
					Add a new custom widget to your dashboard by providing a name and URL.
				</DialogDescription>
			</DialogHeader>

			<form
				@submit="onSubmit"
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
						@click="isDialogOpen = false"
					>
						Cancel
					</Button>
					<Button
						type="submit"
						:disabled="createMutation.fetching.value"
					>
						{{ createMutation.fetching.value ? 'Creating...' : 'Create Widget' }}
					</Button>
				</div>
			</form>
		</DialogContent>
	</Dialog>

	<ActionConfirm
		:open="widgetPendingDeletion !== null"
		:confirm-text="`Are you sure you want to delete widget &quot;${widgetPendingDeletion?.name}&quot;?`"
		@update:open="widgetPendingDeletion = null"
		@confirm="deleteWidget"
	/>
</template>
