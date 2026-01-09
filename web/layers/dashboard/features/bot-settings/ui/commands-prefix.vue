<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { onMounted, watch } from 'vue'

import { z } from 'zod'

import { useCommandsPrefixApi } from '#layers/dashboard/api/commands-prefix'




import { toast } from 'vue-sonner'

const { t } = useI18n()

const api = useCommandsPrefixApi()
const { data: currentPrefix } = api.usePrefix()
const update = api.usePrefixUpdate()
const reset = api.usePrefixReset()

const formSchema = toTypedSchema(
	z.object({
		prefix: z.string().min(1).max(10),
	})
)

const form = useForm({
	validationSchema: formSchema,
})

watch(currentPrefix, (v) => {
	if (!v?.channelsCommandsPrefix) return

	form.setFieldValue('prefix', v.channelsCommandsPrefix)
})

onMounted(() => {
	if (!currentPrefix.value?.channelsCommandsPrefix) return

	form.setFieldValue('prefix', currentPrefix.value?.channelsCommandsPrefix)
})

const onSubmit = form.handleSubmit(async (values) => {
	try {
		await update.executeMutation({
			input: { newPrefix: values.prefix },
		})
		toast.success('Updated')
	} catch {
		toast.error('Error happend on update')
	}
})

async function onReset() {
	try {
		await reset.executeMutation({})
		toast.success('Resetted')
	} catch {
		toast.error('Error happend on reset')
	}
}
</script>

<template>
	<div class="flex flex-col w-full gap-4">
		<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">Commands prefix</h4>
		<UiCard>
			<form @submit="onSubmit">
				<UiCardContent class="pt-6 flex gap-2">
					<UiFormField v-slot="{ componentField }" name="prefix">
						<UiFormItem class="w-full">
							<UiFormLabel>Prefix used for all commands, "!" by default</UiFormLabel>
							<UiFormControl>
								<UiInput
									type="text"
									placeholder="Global commands prefix"
									v-bind="componentField"
									class="w-full"
								/>
							</UiFormControl>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>
				</UiCardContent>
				<UiCardFooter class="flex flex-row gap-2 justify-end mt-4">
					<UiButton type="button" variant="destructive" class="place-self-end" @click="onReset">
						Reset
					</UiButton>
					<UiButton type="submit">
						{{ t('sharedButtons.save') }}
					</UiButton>
				</UiCardFooter>
			</form>
		</UiCard>
	</div>
</template>
