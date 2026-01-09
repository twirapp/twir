<script setup lang="ts">





const { t } = useI18n()
</script>

<template>
	<div>
		<UiFormField v-slot="{ componentField }" name="triggerLength">
			<UiFormItem>
				<UiFormLabel>{{ t('moderation.types.emotes.triggerLength') }}</UiFormLabel>
				<UiFormControl>
					<UiInput type="number" v-bind="componentField" />
				</UiFormControl>
				<UiFormMessage />
			</UiFormItem>
		</UiFormField>
	</div>
</template>
