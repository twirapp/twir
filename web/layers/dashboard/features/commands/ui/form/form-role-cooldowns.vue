<script setup lang="ts">
import { PlusIcon, TrashIcon } from "lucide-vue-next";
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import { useFieldArray, useFormContext } from "vee-validate";

import { Button } from "@/components/ui/button";
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { useCommandEditV2 } from "@/features/commands/composables/use-command-edit-v2";

const { t } = useI18n();
const { channelRoles } = useCommandEditV2();
const { values } = useFormContext();

const { fields, push, remove } = useFieldArray("roleCooldowns");

const availableRoles = computed(() => {
	const selectedRoleIds = new Set(
		((values.roleCooldowns as Array<{ roleId: string; cooldown: number }>) || []).map(
			(rc) => rc.roleId,
		),
	);

	return (
		channelRoles.value?.roles.filter((role) => {
			// Filter out broadcaster role and already selected roles
			if (role.type === "BROADCASTER") return false;
			return !selectedRoleIds.has(role.id);
		}) || []
	);
});

const getRoleName = (roleId: string) => {
	const role = channelRoles.value?.roles.find((r) => r.id === roleId);
	return role?.name || roleId;
};

function addRoleCooldown() {
	if (availableRoles.value.length === 0) return;

	push({
		roleId: availableRoles.value[0].id,
		cooldown: 0,
	});
}

function removeRoleCooldown(index: number) {
	remove(index);
}
</script>

<template>
	<div class="flex flex-col gap-4">
		<div class="flex justify-between items-center">
			<span class="text-sm font-medium">
				{{ t("commands.modal.cooldown.roleCooldowns.title") }}
			</span>
			<Button
				type="button"
				size="sm"
				variant="outline"
				:disabled="availableRoles.length === 0"
				@click="addRoleCooldown"
			>
				<PlusIcon class="size-4 mr-1" />
				{{ t("commands.modal.cooldown.roleCooldowns.add") }}
			</Button>
		</div>

		<div v-if="fields.length === 0" class="text-sm text-muted-foreground">
			{{ t("commands.modal.cooldown.roleCooldowns.empty") }}
		</div>

		<div v-for="(field, index) in fields" :key="field.key" class="flex gap-2 items-end">
			<FormField v-slot="{ componentField }" :name="`roleCooldowns.${index}.roleId`">
				<FormItem class="flex-1">
					<FormLabel v-if="index === 0">
						{{ t("commands.modal.cooldown.roleCooldowns.role") }}
					</FormLabel>
					<FormControl>
						<Select v-bind="componentField">
							<FormControl>
								<SelectTrigger>
									<SelectValue>
										{{ getRoleName((field.value as any).roleId as string) }}
									</SelectValue>
								</SelectTrigger>
							</FormControl>
							<SelectContent>
								<SelectGroup>
									<SelectItem
										v-for="role in channelRoles?.roles.filter((r) => r.type !== 'BROADCASTER')"
										:key="role.id"
										:value="role.id"
									>
										{{ role.name }}
									</SelectItem>
								</SelectGroup>
							</SelectContent>
						</Select>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" :name="`roleCooldowns.${index}.cooldown`">
				<FormItem class="flex-1">
					<FormLabel v-if="index === 0">
						{{ t("commands.modal.cooldown.roleCooldowns.cooldownValue") }}
					</FormLabel>
					<FormControl>
						<Input type="number" v-bind="componentField" min="0" max="84600" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<Button
				type="button"
				size="icon"
				variant="destructive"
				:class="{ 'mb-0': index === 0 }"
				@click="removeRoleCooldown(index)"
			>
				<TrashIcon class="size-4" />
			</Button>
		</div>

		<p v-if="fields.length > 0" class="text-xs text-muted-foreground">
			{{ t("commands.modal.cooldown.roleCooldowns.description") }}
		</p>
	</div>
</template>

<style scoped></style>
