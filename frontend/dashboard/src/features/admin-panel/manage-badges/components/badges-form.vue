<script setup lang="ts">
import { NCard } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import BadgesPreview from './badges-preview.vue'
import { useBadgesForm } from '../composables/use-badges-form'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const { t } = useI18n()
const badgesForm = useBadgesForm()
</script>

<template>
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.manageBadges.formTitle') }}
	</h4>
	<NCard size="small" bordered>
		<form class="flex flex-col gap-4" @submit="badgesForm.onSubmit">
			<div class="space-y-2">
				<Label for="name">
					{{ t('adminPanel.manageBadges.name') }}
				</Label>
				<Input
					v-model="badgesForm.nameField.fieldModel"
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
					v-model="badgesForm.slotField.fieldModel"
					name="slot"
					type="text"
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
				<div className="grid w-full items-center gap-1.5">
					<Input
						:required="!badgesForm.editableBadgeId"
						accept="image/*"
						type="file"
						@change="badgesForm.setImageField"
					/>
				</div>
			</div>

			<div v-if="badgesForm.isImageFile">
				<Label>
					{{ t('adminPanel.manageBadges.preview') }}
				</Label>
				<BadgesPreview class="mt-2" :image="badgesForm.formValues.image!" />
			</div>

			<div class="flex justify-end gap-4">
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
			</div>
		</form>
	</NCard>
</template>
