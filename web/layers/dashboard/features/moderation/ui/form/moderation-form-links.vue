<script setup lang="ts">
import { FormControl, FormDescription, FormField, FormItem, FormLabel } from '@/components/ui/form'
import { Switch } from '@/components/ui/switch'
</script>

<template>
	<div>
		<FormField v-slot="{ field }" name="checkClips">
			<FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
				<div class="space-y-0.5">
					<FormLabel class="text-base"> Check twitch clips </FormLabel>
					<FormDescription> If enabled bot will check clips as links. </FormDescription>
				</div>
				<FormControl>
					<Switch
						:model-value="field.value"
						default-checked
						@update:model-value="field['onUpdate:modelValue']"
					/>
				</FormControl>
			</FormItem>
		</FormField>
	</div>
</template>
