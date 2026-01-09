<script setup lang="ts">






const { t } = useI18n()
</script>

<template>
	<UiFormField v-slot="{ componentField }" name="triggerLength">
		<UiFormItem>
			<UiFormLabel>{{ t('moderation.types.long_message.triggerLength') }}</UiFormLabel>
			<UiFormControl>
				<UiInput type="number" v-bind="componentField" />
			</UiFormControl>
			<UiFormMessage />
		</UiFormItem>
	</UiFormField>

	<UiSeparator />
</template>
