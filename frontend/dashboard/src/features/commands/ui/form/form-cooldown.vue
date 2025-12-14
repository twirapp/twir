<script setup lang="ts">
import { EditIcon, RefreshCcwIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import FormRolesSelector from '@/features/commands/ui/form-roles-selector.vue'
import CommunityRolesModal from '@/features/community-roles/community-roles-modal.vue'

const { t } = useI18n()
</script>

<template>
	<Card>
		<CardHeader class="flex flex-row place-content-center flex-wrap">
			<CardTitle class="flex items-center gap-2">
				<RefreshCcwIcon />
				{{ t('commands.modal.cooldown.label') }}
			</CardTitle>
		</CardHeader>
		<CardContent class="pt-4">
			<div class="flex flex-col gap-4">
				<FormField v-slot="{ componentField }" name="cooldown">
					<FormItem>
						<FormLabel class="flex gap-2">
							{{ t('commands.modal.cooldown.value') }}
						</FormLabel>
						<FormControl>
							<Input type="number" v-bind="componentField" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="cooldownType">
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

				<div class="flex flex-col gap-2">
					<span class="inline-flex gap-1">
						Affected roles
						<CommunityRolesModal>
							<template #trigger>
								<span class="flex flex-row gap-1 items-center cursor-pointer underline">
									{{ t('sidebar.roles') }}
									<EditIcon class="size-4" />
								</span>
							</template>
						</CommunityRolesModal>
					</span>
					<FormRolesSelector field-name="cooldownRolesIds" hide-broadcaster />
				</div>
			</div>
		</CardContent>
	</Card>
</template>

<style scoped></style>
