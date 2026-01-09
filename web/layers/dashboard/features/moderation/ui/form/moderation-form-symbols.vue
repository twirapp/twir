<script setup lang="ts">





const { t } = useI18n()
</script>

<template>
	<div>
		<UiFormField v-slot="{ componentField }" name="triggerLength">
			<UiFormItem>
				<UiFormLabel>{{ t('moderation.types.symbols.triggerLength') }}</UiFormLabel>
				<UiFormControl>
					<UiInput type="number" v-bind="componentField" />
				</UiFormControl>
				<UiFormMessage />
			</UiFormItem>
		</UiFormField>

		<UiFormField v-slot="{ componentField }" name="maxPercentage">
			<UiFormItem>
				<UiFormLabel>{{ t('moderation.types.symbols.maxPercentage') }}</UiFormLabel>
				<UiFormControl>
					<UiInput type="number" v-bind="componentField" />
				</UiFormControl>
				<UiFormMessage />
			</UiFormItem>
		</UiFormField>
	</div>
</template>
