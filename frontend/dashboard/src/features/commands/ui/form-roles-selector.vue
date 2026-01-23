<script setup lang="ts">
import { useField } from "vee-validate";
import { computed } from "vue";

import type { RoleTypeEnum } from "@/gql/graphql.ts";

import { Checkbox } from "@/components/ui/checkbox";
import { FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form";
import { Label } from "@/components/ui/label";
import { useCommandEditV2 } from "@/features/commands/composables/use-command-edit-v2";

const props = defineProps<{
	fieldName: string;
	excludedTypes?: RoleTypeEnum[];
	hideEveryone?: boolean;
	hideBroadcaster?: boolean;
	everyoneAlwaysActive?: boolean;
}>();

defineSlots<{
	extra?: (props: { roleId?: string; roleName?: string; roleType?: RoleTypeEnum }) => any;
}>();

const { channelRoles } = useCommandEditV2();

const { setValue } = useField(props.fieldName);

function uncheckAll() {
	setValue([]);
}

const roles = computed(() => {
	if (!channelRoles.value?.roles) return [];
	return channelRoles.value.roles.filter((r) => !props.excludedTypes?.includes(r.type));
});
</script>

<template>
	<div class="grid grid-cols-1 gap-1">
		<FormField
			v-for="(role, index) in roles"
			v-slot="{ value, handleChange }"
			:key="role.id"
			type="checkbox"
			:value="role.id"
			:unchecked-value="false"
			:name="fieldName"
		>
			<div v-if="index === 0 && everyoneAlwaysActive && !hideEveryone" class="role-container space-y-0">
				<div class="role w-2xs max-w-2xs" @click="uncheckAll">
					<Checkbox id="allRoles" :model-value="everyoneAlwaysActive ?? !value?.length" />
					<Label for="allRoles" class="capitalize">Everyone</Label>
				</div>
				<slot name="extra" />
			</div>

			<FormItem class="space-y-0">
				<div class="role-container">
					<FormLabel class="role w-2xs max-w-2xs">
						<FormControl>
							<Checkbox
								:model-value="value?.includes(role.id)"
								@update:model-value="handleChange"
							/>
						</FormControl>
						<span>{{ role.name.at(0) + role.name.slice(1).toLowerCase() }}</span>
					</FormLabel>
					<slot name="extra" :role-id="role.id" :role-name="role.name" :role-type="role.type" />
				</div>
			</FormItem>
		</FormField>
	</div>
</template>

<style scoped>
@reference '@/assets/index.css';

.role-container {
	@apply flex flex-wrap items-center gap-2 w-full;
}

.role {
	@apply flex flex-row items-center gap-2 space-y-0 bg-accent px-3 py-2 rounded-md leading-5;
}
</style>
