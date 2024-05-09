import { RpcError } from '@protobuf-ts/runtime-rpc'
import { createGlobalState } from '@vueuse/core'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useLeaveFeedback } from '@/api/feedback.js'

export const useFeedbackForm = createGlobalState(() => {
	const formApi = useLeaveFeedback()
	const { t } = useI18n()

	const error = ref<string | null>(null)
	const form = ref({
		message: '',
	})

	function clearForm() {
		form.value.message = ''
		error.value = ''
	}

	async function submitForm() {
		try {
			await formApi.mutateAsync(form.value.message)
		} catch (e) {
			if (e instanceof RpcError) {
				if (e.code === 'resource_exhausted') {
					error.value = t('feedback.rateLimited', { time: e.meta.retry_after })
				} else {
					error.value = e.message
				}

				throw e
			} else {
				error.value = 'Unexpected error happened, shared this into our discord'
				throw e
			}
		}
	}

	return {
		form,
		error,
		submitForm,
		clearForm,
	}
})
