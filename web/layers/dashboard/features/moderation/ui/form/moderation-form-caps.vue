<script setup lang="ts">





const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col gap-4">
		<UiFormField v-slot="{ componentField }" name="triggerLength">
			<UiFormItem>
				<UiFormLabel>{{ t('moderation.types.caps.triggerLength') }}</UiFormLabel>
				<UiFormControl>
					<UiInput type="number" v-bind="componentField" />
				</UiFormControl>
				<UiFormMessage />
			</UiFormItem>
		</UiFormField>

		<UiFormField v-slot="{ componentField }" name="maxPercentage">
			<UiFormItem>
				<UiFormLabel>{{ t('moderation.types.caps.maxPercentage') }}</UiFormLabel>
				<UiFormControl>
					<UiInput type="number" v-bind="componentField" />
				</UiFormControl>
				<UiFormMessage />
			</UiFormItem>
		</UiFormField>
	</div>
</template>
