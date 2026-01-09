<script setup lang="ts">


import { useForm } from '../composables/use-form.js'




const { formValue } = useForm()
const { t } = useI18n()
</script>

<template>
	<div class="grid items-center gap-1.5">
		<UiLabel for="ignoreTimeoutFrom">{{ t('chatAlerts.ban.ignoreTimeoutFrom') }}</UiLabel>
		<UiTagsInput
			v-if="formValue.ban && 'ignoreTimeoutFrom' in formValue.ban"
			id="ignoreTimeoutFrom"
			v-model="formValue.ban.ignoreTimeoutFrom"
		>
			<UiTagsInputItem v-for="item in formValue.ban.ignoreTimeoutFrom" :key="item" :value="item">
				<UiTagsInputItemText />
				<UiTagsInputItemDelete />
			</UiTagsInputItem>

			<UiTagsInputInput placeholder="Type usernames..." />
		</UiTagsInput>
	</div>
</template>
