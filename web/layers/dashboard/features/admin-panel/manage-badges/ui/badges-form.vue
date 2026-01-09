<script setup lang="ts">


import BadgesPreview from './badges-preview.vue'
import { useBadgesForm } from '../composables/use-badges-form.js'






const { t } = useI18n()
const badgesForm = useBadgesForm()
</script>

<template>
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.manageBadges.formTitle') }}
	</h4>

	<UiCard>
		<form class="flex flex-col gap-4" @submit="badgesForm.onSubmit">
			<UiCardContent class="flex flex-col gap-4 p-4">
				<div class="space-y-2">
					<UiLabel for="name">
						{{ t('adminPanel.manageBadges.name') }}
					</UiLabel>
					<UiInput
						v-model="badgesForm.nameField.fieldModel.value"
						name="name"
						type="text"
						placeholder=""
						required
					/>
				</div>

				<div class="space-y-2">
					<UiLabel for="slot">
						{{ t('adminPanel.manageBadges.slot') }}
					</UiLabel>
					<UiInput
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
					<UiLabel for="file">
						{{ t('adminPanel.manageBadges.image') }}
					</UiLabel>
					<div class="grid w-full items-center gap-1.5">
						<UiInput
							:required="!badgesForm.editableBadgeId"
							accept="image/*"
							type="file"
							@change="badgesForm.setImageField"
						/>
					</div>
				</div>

				<div v-if="badgesForm.isImageFile.value">
					<UiLabel>
						{{ t('adminPanel.manageBadges.preview') }}
					</UiLabel>
					<BadgesPreview class="mt-2" :image="badgesForm.formValues.value.image!" />
				</div>
			</UiCardContent>
			<UiCardFooter class="flex justify-end gap-2 p-4">
				<UiButton
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
				</UiButton>
				<UiButton type="submit">
					<template v-if="badgesForm.editableBadgeId">
						{{ t('sharedButtons.edit') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.create') }}
					</template>
				</UiButton>
			</UiCardFooter>
		</form>
	</UiCard>
</template>
