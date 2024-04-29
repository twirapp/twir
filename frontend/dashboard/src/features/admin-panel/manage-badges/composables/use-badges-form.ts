import { toTypedSchema } from '@vee-validate/zod'
import { defineStore, storeToRefs } from 'pinia'
import { computed, ref, watch } from 'vue'
import * as z from 'zod'

import { useBadges } from './use-badges.js'

import { useFormField } from '@/composables/use-form-field.js'

const formSchema = toTypedSchema(z.object({
	name: z.string(),
	slot: z.number(),
	image: z.any()
}))

export const useBadgesForm = defineStore('admin-panel/badges-form', () => {
	const badgesApi = useBadges()
	const { badges } = storeToRefs(badgesApi)

	const editableBadgeId = ref<string | null>(null)
	const isEditableForm = computed(() => Boolean(editableBadgeId.value))

	const nameField = useFormField<string>('name', '')
	const fileField = useFormField<string | File | null>('image', '')
	const slotField = useFormField<number>('ffzSlot', 256)
	const fileInputRef = computed({
		get() {
			return fileField.fieldRef.value
		},
		set(el: any) {
			fileField.fieldRef.value = el?.$el
		}
	})

	watch(badges, (badges) => {
		const lastBadge = badges.at(-1)
		if (lastBadge) {
			slotField.fieldModel.value = lastBadge.ffzSlot + 1
		}
	})

	const formValues = computed(() => {
		return {
			name: nameField.fieldModel.value,
			image: fileField.fieldModel.value,
			slot: slotField.fieldModel.value
		}
	})
	const isFormDirty = computed(() => Boolean(formValues.value.name || formValues.value.image))
	const isImageFile = computed(() => formValues.value.image instanceof File)

	async function onSubmit(event: Event) {
		event.preventDefault()

		try {
			const { value } = await formSchema.parse(formValues.value)
			if (!value)
				return

			if (editableBadgeId.value) {
				await badgesApi.badgesUpdate.executeMutation({
					id: editableBadgeId.value,
					opts: {
						name: value.name,
						ffzSlot: value.slot,
						file: value.image instanceof File ? value.image : undefined
					}
				})
			} else {
				await badgesApi.badgesCreate.executeMutation({
					opts: {
						name: value.name,
						file: value.image,
						ffzSlot: value.slot
					}
				})
			}
		} catch (err) {
			console.error(err)
		}

		onReset()
	}

	function onReset(): void {
		nameField.reset()
		fileField.reset()
		editableBadgeId.value = null
	}

	function setImageField(event: Event): void {
		const files = (event.target as HTMLInputElement).files
		if (!files)
			return
		fileField.fieldModel.value = files[0]
	}

	return {
		formValues,

		nameField,
		fileField,
		fileInputRef,
		slotField,

		isFormDirty,
		isImageFile,

		isEditableForm,
		editableBadgeId,

		onSubmit,
		onReset,
		setImageField
	}
})
