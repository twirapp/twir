<script setup lang="ts">





const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col gap-4">
		<div class="flex flex-row flex-wrap w-full gap-2">
			<UiFormField v-slot="{ componentField }" name="triggerLength">
				<UiFormItem class="w-full xl:w-[45%]">
					<UiFormLabel>{{ t('moderation.types.one_man_spam.triggerLength') }}</UiFormLabel>
					<UiFormControl>
						<UiInput type="number" v-bind="componentField" />
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ componentField }" name="oneManSpamMinimumStoredMessages">
				<UiFormItem class="w-full xl:w-[45%]">
					<UiFormLabel>{{ t('moderation.types.one_man_spam.minimumStoredMessages') }}</UiFormLabel>
					<UiFormControl>
						<UiInput type="number" v-bind="componentField" min="1" max="20" />
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>
		</div>

		<div class="flex flex-row w-full flex-wrap gap-2">
			<UiFormField v-slot="{ componentField }" name="oneManSpamMessageMemorySeconds">
				<UiFormItem class="w-full xl:w-[45%]">
					<UiFormLabel>{{ t('moderation.types.one_man_spam.messageMemorySeconds') }}</UiFormLabel>
					<UiFormControl>
						<UiInput type="number" v-bind="componentField" min="1" max="600" />
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ componentField }" name="maxPercentage">
				<UiFormItem class="w-full xl:w-[45%]">
					<UiFormLabel>{{ t('moderation.types.one_man_spam.maxPercentage') }}</UiFormLabel>
					<UiFormControl>
						<UiInput type="number" v-bind="componentField" min="1" max="100" />
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>
		</div>
	</div>
</template>
