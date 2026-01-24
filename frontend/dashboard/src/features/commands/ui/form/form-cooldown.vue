<script setup lang="ts">
import { EditIcon, RefreshCcwIcon } from "lucide-vue-next";
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import { useFormContext } from "vee-validate";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
import { Separator } from "@/components/ui/separator";
import CommunityRolesModal from "@/features/community-roles/community-roles-modal.vue";
import type {
	FormSchema} from '@/features/commands/composables/use-command-edit-v2.ts';
import {
	useCommandEditV2,
} from '@/features/commands/composables/use-command-edit-v2.ts';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';

const { t } = useI18n();
const { channelRoles } = useCommandEditV2()
const { values, setFieldValue } = useFormContext<FormSchema>();

const roleCooldowns = computed({
	get: () => {
		const cooldowns = values.roleCooldowns || [];
		return new Map(cooldowns.map((rc) => [rc.roleId, rc.cooldown]));
	},
	set: (map: Map<string, number>) => {
		const cooldowns = Array.from(map.entries())
			.filter(([_, cooldown]) => cooldown > 0)
			.map(([roleId, cooldown]) => ({ roleId, cooldown }));

		setFieldValue("roleCooldowns", cooldowns);
	},
});

function updateRoleCooldown(roleId: string, cooldown: number) {
	const map = new Map(roleCooldowns.value);
	if (cooldown > 0) {
		map.set(roleId, cooldown);
	} else {
		map.delete(roleId);
	}
	roleCooldowns.value = map;
}

function getRoleCooldown(roleId: string): number {
	return roleCooldowns.value.get(roleId) ?? 0;
}
</script>

<template>
	<Card>
		<CardHeader class="flex flex-row place-content-center flex-wrap">
			<CardTitle class="flex items-center gap-2">
				<RefreshCcwIcon />
				{{ t("commands.modal.cooldown.label") }}
			</CardTitle>
		</CardHeader>
		<CardContent class="pt-4">
			<div class="flex flex-col gap-4">
				<FormField v-slot="{ componentField }" name="cooldownType">
					<FormItem>
						<FormLabel>{{ t("commands.modal.cooldown.type.name") }}</FormLabel>
						<FormControl>
							<Select v-bind="componentField">
								<FormControl>
									<SelectTrigger>
										<SelectValue />
									</SelectTrigger>
								</FormControl>
								<SelectContent>
									<SelectGroup>
										<SelectItem value="GLOBAL">
											{{ t("commands.modal.cooldown.type.global") }}
										</SelectItem>
										<SelectItem value="PER_USER">
											{{ t("commands.modal.cooldown.type.user") }}
										</SelectItem>
									</SelectGroup>
								</SelectContent>
							</Select>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<Separator />

				<div class="flex flex-col gap-2">
					<span class="inline-flex gap-1 text-sm font-medium">
						{{ t("commands.modal.cooldown.roleCooldowns.title") }}
						<CommunityRolesModal>
							<template #trigger>
								<span class="flex flex-row gap-1 items-center cursor-pointer underline">
									{{ t("sidebar.roles") }}
									<EditIcon class="size-4" />
								</span>
							</template>
						</CommunityRolesModal>
					</span>

					<p class="text-xs text-muted-foreground">
						{{ t("commands.modal.cooldown.roleCooldowns.description") }}
					</p>

					<div class="@container w-full max-w-2xl">
						<div class="grid grid-cols-1 gap-1 xl:max-w-[50%]">
							<div class="flex flex-row items-center gap-2 space-y-0">
								<div class="flex flex-row gap-2 bg-accent px-3 py-2 rounded-md leading w-56 min-w-56">
									<Checkbox id="allRoles" :model-value="true" disabled />
									<Label for="allRoles" class="capitalize">Everyone</Label>
								</div>
								<Input
									v-model="values.cooldown"
									min="0"
									max="86400"
									class="w-auto"
									:placeholder="t('commands.modal.cooldown.value')"
								/>
								sec
							</div>

							<div
								v-for="(role) in channelRoles?.roles"
								:key="role!.id"
								class="flex flex-row items-center gap-2 space-y-0"
							>
								<div class="flex flex-row fle flex-wrap gap-2 space-y-0 bg-accent px-3 py-2 rounded-md leading w-56 min-w-56">
									<Checkbox
										:model-value="roleCooldowns.has(role.id)"
										@update:model-value="roleCooldowns.has(role.id) ? updateRoleCooldown(role.id, 0) : updateRoleCooldown(role.id, 5)"
										:id="`cooldown-${role.id}`"
									/>
									<Label :for="`cooldown-${role.id}`" class="cursor-pointer text-ellipsis overflow-hidden">
										{{ role.name }}
									</Label>
								</div>
								<Input
									:disabled="!roleCooldowns.has(role.id)"
									:model-value="getRoleCooldown(role.id)"
									@update:model-value="(val) => updateRoleCooldown(role.id, Number(val) ?? 5)"
									min="0"
									max="86400"
									class="w-auto"
									:placeholder="t('commands.modal.cooldown.value')"
								/>
								sec
							</div>
						</div>
					</div>
				</div>
			</div>
		</CardContent>
	</Card>
</template>
