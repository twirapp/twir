<script setup lang="ts">
import { Trash2, Plus } from "lucide-vue-next";
import { ref, computed } from "vue";

import {
	useDashboardWidgetsCreateCustom,
	useDashboardWidgetsDelete,
	useDashboardWidgetsLayout,
} from "@/api/dashboard-widgets-layout.ts";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { z } from "zod";
import { toast } from "vue-sonner";

const { layout, fetching } = useDashboardWidgetsLayout();
const createMutation = useDashboardWidgetsCreateCustom();
const deleteMutation = useDashboardWidgetsDelete();

// Filter only custom widgets from the layout
const customWidgets = computed(() => {
	return layout.value
		.filter((w) => w.type === "CUSTOM")
		.map((w) => ({
			id: w.widgetId,
			name: w.customName || "Unnamed Widget",
			url: w.customUrl || "",
		}));
});

const isDialogOpen = ref(false);

const formSchema = z.object({
	name: z.string().min(2, "Name must be at least 2 characters."),
	url: z.string().url("Must be a valid URL"),
});

const { handleSubmit, resetForm } = useForm({
	validationSchema: toTypedSchema(formSchema),
});

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
	});

	if (result.error) {
		toast.error("Failed to create widget", {
			description: result.error.message,
		});
	} else {
		toast.success("Widget created successfully");
		resetForm();
		isDialogOpen.value = false;
	}
});

async function deleteWidget(id: string, name: string) {
	if (!confirm(`Are you sure you want to delete widget "${name}"?`)) {
		return;
	}

	const result = await deleteMutation.executeMutation({ widgetId: id });

	if (result.error) {
		toast.error("Failed to delete widget", {
			description: result.error.message,
		});
	} else {
		toast.success("Widget deleted successfully");
	}
}
</script>

<template>
	<div class="p-6">
		<h1 class="text-3xl font-bold mb-6">Custom Dashboard Widgets</h1>

		<div class="max-w-2xl space-y-6">
			<!-- Widgets List -->
			<div class="bg-card p-6 rounded-lg border">
				<h2 class="text-xl font-semibold mb-4">Your Custom Widgets</h2>

				<div v-if="fetching" class="text-center py-8 text-muted-foreground">Loading...</div>

				<div v-else-if="customWidgets.length === 0" class="text-center py-8 text-muted-foreground">
					No custom widgets yet. Click the + button to create one!
				</div>

				<div v-else class="space-y-3">
					<div
						v-for="widget in customWidgets"
						:key="widget.id"
						class="flex items-center justify-between p-4 bg-secondary/50 rounded-lg"
					>
						<div class="flex-1">
							<h3 class="font-medium">{{ widget.name }}</h3>
							<p class="text-sm text-muted-foreground truncate">
								{{ widget.url }}
							</p>
						</div>
						<Button
							variant="destructive"
							size="icon"
							@click="deleteWidget(widget.id, widget.name)"
							:disabled="deleteMutation.fetching.value"
						>
							<Trash2 class="h-4 w-4" />
						</Button>
					</div>
				</div>
			</div>

			<!-- Info Box -->
			<div class="bg-blue-500/10 border border-blue-500/20 p-4 rounded-lg">
				<h3 class="font-semibold text-blue-500 mb-2">How to use:</h3>
				<ol class="text-sm space-y-1 list-decimal list-inside text-muted-foreground">
					<li>Create a custom widget by clicking the + button</li>
					<li>Go to your dashboard page</li>
					<li>Click the "+" button in the bottom right</li>
					<li>Select your custom widget to add it to the dashboard</li>
					<li>Drag and resize it like any other widget</li>
				</ol>
			</div>
		</div>
	</div>

	<!-- Floating Action Button -->
	<Dialog v-model:open="isDialogOpen">
		<DialogTrigger as-child>
			<Button size="icon" class="fixed right-8 bottom-8 h-14 w-14 rounded-full shadow-lg z-50">
				<Plus class="h-6 w-6" />
			</Button>
		</DialogTrigger>
		<DialogContent>
			<DialogHeader>
				<DialogTitle>Create Custom Widget</DialogTitle>
				<DialogDescription>
					Add a new custom widget to your dashboard by providing a name and URL.
				</DialogDescription>
			</DialogHeader>

			<form @submit="onSubmit" class="space-y-4">
				<FormField v-slot="{ componentField }" name="name">
					<FormItem>
						<FormLabel>Widget Name</FormLabel>
						<FormControl>
							<Input v-bind="componentField" placeholder="My Custom Widget" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="url">
					<FormItem>
						<FormLabel>Website URL</FormLabel>
						<FormControl>
							<Input v-bind="componentField" placeholder="https://example.com" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<div class="flex justify-end gap-2">
					<Button type="button" variant="outline" @click="isDialogOpen = false"> Cancel </Button>
					<Button type="submit" :disabled="createMutation.fetching.value">
						{{ createMutation.fetching.value ? "Creating..." : "Create Widget" }}
					</Button>
				</div>
			</form>
		</DialogContent>
	</Dialog>
</template>
