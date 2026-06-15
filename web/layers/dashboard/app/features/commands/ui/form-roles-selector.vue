<script setup lang="ts">
import { useField } from 'vee-validate'
import { computed } from 'vue'

import type { RoleTypeEnum } from '~/gql/graphql.js'

import { Checkbox } from '@/components/ui/checkbox'
import { FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form'
import { Label } from '@/components/ui/label'
import { useCommandEditV2 } from '~~/layers/dashboard/app/features/commands/composables/use-command-edit-v2'

const props = defineProps<{
	fieldName: string
	excludedTypes?: RoleTypeEnum[]
	hideEveryone?: boolean
	hideBroadcaster?: boolean
}>()

const { channelRoles } = useCommandEditV2()

const { setValue } = useField(props.fieldName)

function uncheckAll() {
	setValue([])
}

const roles = computed(() => {
	if (!channelRoles.value?.roles) return []
	return channelRoles.value.roles.filter((r) => !props.excludedTypes?.includes(r.type))
})
</script>

<template>
	<div class="grid grid-cols-1 md:grid-cols-2 gap-1 xl:max-w-[50%]">
		<FormField
			v-for="(role, index) in roles"
			v-slot="{ value, handleChange }"
			:key="role.id"
			type="checkbox"
			:value="role.id"
			:unchecked-value="false"
			:name="fieldName"
		>
			<div
				v-if="index === 0 && !hideEveryone"
				class="flex flex-row items-center gap-2 space-y-0 rounded-md bg-accent px-3 py-2 leading-5"
				@click="uncheckAll"
			>
				<Checkbox id="allRoles" :model-value="!value?.length" />
				<Label for="allRoles" class="capitalize">Everyone</Label>
			</div>

			<FormItem class="space-y-0">
				<FormLabel class="flex flex-row items-center gap-2 space-y-0 rounded-md bg-accent px-3 py-2 leading-5">
					<FormControl>
						<Checkbox :model-value="value?.includes(role.id)" @update:model-value="handleChange" />
					</FormControl>
					<span>{{ role.name.at(0) + role.name.slice(1).toLowerCase() }}</span>
				</FormLabel>
			</FormItem>
		</FormField>
	</div>
</template>
