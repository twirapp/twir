<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { Label } from '@/components/ui/label'
import { TagsInput, TagsInputInput, TagsInputItem, TagsInputItemDelete, TagsInputItemText } from '@/components/ui/tags-input'
import { useForm } from '@/features/chat-alerts/composables/use-form.js'

const { formValue } = useForm()
const { t } = useI18n()
</script>

<template>
	<div class="grid items-center gap-1.5">
		<Label for="ignoreTimeoutFrom">{{ t('chatAlerts.ban.ignoreTimeoutFrom') }}</Label>
		<TagsInput
			v-if="formValue.ban && 'ignoreTimeoutFrom' in formValue.ban"
			id="ignoreTimeoutFrom"
			v-model="formValue.ban.ignoreTimeoutFrom"
		>
			<TagsInputItem v-for="item in formValue.ban.ignoreTimeoutFrom" :key="item" :value="item">
				<TagsInputItemText />
				<TagsInputItemDelete />
			</TagsInputItem>

			<TagsInputInput placeholder="Type usernames..." />
		</TagsInput>
	</div>
</template>
