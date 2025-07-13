<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Plus, X } from 'lucide-vue-next'
import { useFieldArray } from 'vee-validate'

const {
	fields: excludedEmotes,
	push: addExcludedEmote,
	remove: removeExcludedEmote,
} = useFieldArray('excludedEmotes')
</script>

<template>
	<div class="space-y-6">
		<Card>
			<CardHeader>
				<CardTitle>General Settings</CardTitle>
				<CardDescription>Configure basic Kappagen overlay behavior</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="flex flex-col md:flex-row gap-4">
					<FormField name="enableSpawn" v-slot="{ value, handleChange }">
						<FormItem class="flex flex-col rounded-lg border p-3 shadow-sm">
							<div class="flex flex-row w-full items-center justify-between">
								<FormLabel>Enable Spawn</FormLabel>
								<div>
									<Switch :checked="value" @update:checked="handleChange" />
								</div>
							</div>
							<div class="text-[0.8rem] text-muted-foreground">
								Spawn emote on each chat message
							</div>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField name="enableRave" v-slot="{ value, handleChange }">
						<FormItem class="flex flex-col rounded-lg border p-3 shadow-sm">
							<div class="flex flex-row w-full items-center justify-between">
								<FormLabel>Enable Rave Mode</FormLabel>
								<div>
									<Switch :checked="value" @update:checked="handleChange" />
								</div>
							</div>
							<div class="text-[0.8rem] text-muted-foreground">
								Enable special rave animations and effects
							</div>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<CardTitle>Excluded Emotes</CardTitle>
				<CardDescription>Specify emotes that should not appear in the overlay</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="space-y-2">
					<div
						v-for="(field, index) in excludedEmotes"
						:key="field.key"
						class="flex items-center gap-2"
					>
						<FormField :name="`excludedEmotes[${index}]`" v-slot="{ componentField }">
							<FormItem class="flex-1">
								<Input v-bind="componentField" placeholder="Enter emote name" />
							</FormItem>
						</FormField>
						<Button type="button" variant="outline" size="sm" @click="removeExcludedEmote(index)">
							<X class="h-4 w-4" />
						</Button>
					</div>
				</div>

				<Button type="button" variant="outline" size="sm" @click="() => addExcludedEmote('')">
					<Plus class="h-4 w-4 mr-2" />
					Add Excluded Emote
				</Button>
			</CardContent>
		</Card>
	</div>
</template>
