<script lang="ts" setup>
import { EditIcon, RefreshCcwIcon } from 'lucide-vue-next'
import { useFormContext } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { useCommandEditV2 } from '@/features/commands/composables/use-command-edit-v2.ts'
import CommunityRolesModal from '@/features/community-roles/community-roles-modal.vue'

import type { FormSchema } from '@/features/commands/composables/use-command-edit-v2.ts'

const { t } = useI18n()
const { channelRoles } = useCommandEditV2()
const { values, setFieldValue } = useFormContext<FormSchema>()

const roleCooldowns = computed({
	get: () => {
		const cooldowns = values.roleCooldowns || []
		return new Map(cooldowns.map((rc) => [rc.roleId, rc.cooldown]))
	},
	set: (map: Map<string, number>) => {
		const cooldowns = Array.from(map.entries()).map(([roleId, cooldown]) => ({ roleId, cooldown }))

		setFieldValue('roleCooldowns', cooldowns)
	},
})

function updateRoleCooldown(roleId: string, cooldown: number | null) {
	const map = new Map(roleCooldowns.value)
	if (cooldown !== null) {
		const safe = Number.isFinite(cooldown) ? Math.min(Math.max(Math.trunc(cooldown), 0), 86400) : 0
		map.set(roleId, safe)
	} else {
		map.delete(roleId)
	}
	roleCooldowns.value = map
}

function getRoleCooldown(roleId: string): number {
	return roleCooldowns.value.get(roleId) ?? 0
}
</script>

<template>
	<Card>
		<CardHeader class="flex flex-row flex-wrap place-content-center">
			<CardTitle class="flex items-center gap-2">
				<RefreshCcwIcon />
				{{ t('commands.modal.cooldown.label') }}
			</CardTitle>
		</CardHeader>
		<CardContent class="pt-4">
			<div class="flex flex-col gap-4">
				<FormField
					v-slot="{ componentField }"
					name="cooldownType"
				>
					<FormItem>
						<FormLabel>{{ t('commands.modal.cooldown.type.name') }}</FormLabel>
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
											{{ t('commands.modal.cooldown.type.global') }}
										</SelectItem>
										<SelectItem value="PER_USER">
											{{ t('commands.modal.cooldown.type.user') }}
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
						{{ t('commands.modal.cooldown.roleCooldowns.title') }}
						<CommunityRolesModal>
							<template #trigger>
								<span class="flex cursor-pointer flex-row items-center gap-1 underline">
									{{ t('sidebar.roles') }}
									<EditIcon class="size-4" />
								</span>
							</template>
						</CommunityRolesModal>
					</span>

					<p class="text-muted-foreground text-xs">
						{{ t('commands.modal.cooldown.roleCooldowns.description') }}
					</p>

					<div class="@container w-full max-w-2xl">
						<div class="grid grid-cols-1 gap-1 xl:max-w-[70%]">
							<div class="flex flex-row flex-wrap items-center gap-2 space-y-0">
								<div
									class="bg-accent leading flex w-56 min-w-56 flex-row gap-2 rounded-md px-3 py-2"
								>
									<Checkbox
										id="allRoles"
										:model-value="true"
										disabled
									/>
									<Label
										class="capitalize"
										for="allRoles"
										>Viewer</Label
									>
								</div>
								<Input
									:model-value="values.cooldown ?? 0"
									:placeholder="t('commands.modal.cooldown.value')"
									class="w-auto"
									max="86400"
									min="0"
									type="number"
									@update:model-value="
										(v) => {
											const n = parseInt(String(v), 10)
											setFieldValue(
												'cooldown',
												Number.isFinite(n) ? Math.min(Math.max(n, 0), 86400) : 0
											)
										}
									"
								/>
								sec
							</div>

							<div
								v-for="role in channelRoles?.roles"
								:key="role!.id"
								class="flex flex-row items-center gap-2 space-y-0"
							>
								<div
									class="fle bg-accent leading flex w-56 min-w-56 flex-row flex-wrap gap-2 space-y-0 rounded-md px-3 py-2"
								>
									<Checkbox
										:id="`cooldown-${role.id}`"
										:model-value="roleCooldowns.has(role.id)"
										@update:model-value="
											roleCooldowns.has(role.id)
												? updateRoleCooldown(role.id, null)
												: updateRoleCooldown(role.id, 0)
										"
									/>
									<Label
										:for="`cooldown-${role.id}`"
										class="cursor-pointer overflow-hidden text-ellipsis"
									>
										{{ role.name }}
									</Label>
								</div>
								<Input
									:disabled="!roleCooldowns.has(role.id)"
									:model-value="getRoleCooldown(role.id)"
									:placeholder="t('commands.modal.cooldown.value')"
									class="w-auto"
									max="86400"
									min="0"
									type="number"
									@update:model-value="
										(val) => {
											const n = parseInt(String(val), 10)
											updateRoleCooldown(role.id, Number.isFinite(n) ? n : 0)
										}
									"
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
