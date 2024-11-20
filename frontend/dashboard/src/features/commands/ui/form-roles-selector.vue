<script setup lang="ts">
import { Checkbox } from '@/components/ui/checkbox'
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form'
import { Label } from '@/components/ui/label'
import { useCommandEditV2 } from '@/features/commands/composables/use-command-edit-v2'

defineProps<{ fieldName: string }>()

const { channelRoles } = useCommandEditV2()
</script>

<template>
	<div class="grid grid-cols-1 md:grid-cols-2 gap-1 xl:max-w-[50%]">
		<FormField
			v-for="(role, index) in channelRoles?.roles" v-slot="{ value, handleChange }"
			:key="role.id"
			type="checkbox"
			:value="role.id"
			:unchecked-value="false"
			:name="fieldName"
		>
			<div v-if="index === 0" class="role">
				<Checkbox id="allRoles" :checked="!value?.length" disabled />
				<Label for="allRoles" class="capitalize">Everyone</Label>
			</div>

			<FormItem class="space-y-0">
				<FormLabel class="role">
					<FormControl>
						<Checkbox
							:checked="value?.includes(role.id)"
							@update:checked="handleChange"
						/>
					</FormControl>
					<span>{{ role.name }}</span>
				</FormLabel>
			</FormItem>
		</FormField>
	</div>
</template>

<style scoped>
.role {
	@apply flex flex-row items-center gap-2 space-y-0 bg-accent px-3 py-2 rounded-md leading-5
}
</style>
