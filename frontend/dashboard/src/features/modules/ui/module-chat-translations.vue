<script setup lang="ts">
import { LanguagesIcon, SettingsIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { toTypedSchema } from 'vee-validate/zod'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import Card from '@/components/card/card.vue'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'

const { t } = useI18n()

const showSettings = ref(false)

const formSchema = z.object()

const translationsForm = useForm({
	validationSchema: toTypedSchema(formSchema),
})

const handleSubmit = translationsForm.handleSubmit(async (values) => {
	console.log(values)
})
</script>

<template>
	<Card
		title="Chat translations"
		:is-loading="false"
		:icon="LanguagesIcon"
		icon-height="30px"
		icon-width="30px"
		description="Translate chat in real time"
	>
		<template #footer>
			<Button class="flex gap-2 items-center" variant="secondary" @click="showSettings = !showSettings">
				{{ t('sharedTexts.settings') }}
				<SettingsIcon class="size-4" />
			</Button>
		</template>
	</Card>

	<Dialog v-model:open="showSettings">
		<DialogOrSheet>
			<form @submit.prevent="handleSubmit"></form>
		</DialogOrSheet>
	</Dialog>
</template>
