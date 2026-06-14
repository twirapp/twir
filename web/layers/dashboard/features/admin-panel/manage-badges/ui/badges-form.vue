<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import BadgesPreview from './badges-preview.vue'
import { useBadgesForm } from '../composables/use-badges-form.js'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const { t } = useI18n()
const badgesForm = useBadgesForm()
</script>

<template>
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.manageBadges.formTitle') }}
	</h4>

	<Card>
		<form class="flex flex-col gap-4" @submit="badgesForm.onSubmit">
			<CardContent class="flex flex-col gap-4 p-4">
				<div class="space-y-2">
					<Label for="name">
						{{ t('adminPanel.manageBadges.name') }}
					</Label>
					<Input
						v-model="badgesForm.nameField.fieldModel.value"
						name="name"
						type="text"
						placeholder=""
						required
					/>
				</div>

				<div class="space-y-2">
					<Label for="slot">
						{{ t('adminPanel.manageBadges.slot') }}
					</Label>
					<Input
						v-model.number="badgesForm.slotField.fieldModel.value"
						type="number"
						name="slot"
						inputmode="numeric"
						pattern="[0-9]*"
						placeholder=""
						required
					/>
				</div>

				<div class="space-y-2">
					<Label for="file">
						{{ t('adminPanel.manageBadges.image') }}
					</Label>
					<div class="grid w-full items-center gap-1.5">
						<Input
							:required="!badgesForm.editableBadgeId"
							accept="image/*"
							type="file"
							@change="badgesForm.setImageField"
						/>
					</div>
				</div>

				<div v-if="badgesForm.isImageFile.value">
					<Label>
						{{ t('adminPanel.manageBadges.preview') }}
					</Label>
					<BadgesPreview class="mt-2" :image="badgesForm.formValues.value.image!" />
				</div>
			</CardContent>
			<CardFooter class="flex justify-end gap-2 p-4">
				<Button
					type="button"
					variant="secondary"
					:disabled="!badgesForm.isFormDirty"
					@click="badgesForm.onReset"
				>
					<template v-if="badgesForm.editableBadgeId">
						{{ t('sharedButtons.cancel') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.reset') }}
					</template>
				</Button>
				<Button type="submit">
					<template v-if="badgesForm.editableBadgeId">
						{{ t('sharedButtons.edit') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.create') }}
					</template>
				</Button>
			</CardFooter>
		</form>
	</Card>
</template>
