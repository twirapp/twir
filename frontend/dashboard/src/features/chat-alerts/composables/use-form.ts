import { createGlobalState, useDebounceFn, watchIgnorable } from '@vueuse/core'
import { useNotification } from 'naive-ui'
import { ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import type { ChatAlerts } from '@/gql/graphql'
import type { KeysOfUnion, RequiredDeep, SetNonNullable } from 'type-fest'

import { useChatAlertsApi } from '@/api/chat-alerts.js'

export type FormKey = Exclude<KeysOfUnion<RequiredDeep<SetNonNullable<ChatAlerts>>>, '__typename'>

type OmitDeep<T, K extends string> = T extends object
	? T extends Array<infer U>
		? OmitDeep<U, K>[]
		: {
			[P in keyof T as P extends K ? never : P]: OmitDeep<T[P], K>
		}
	: T

type NonNullableFields<T> = {
	[P in keyof T]-?: NonNullable<T[P]>;
}

type Form = OmitDeep<NonNullableFields<ChatAlerts>, '__typename'>

export const useForm = createGlobalState(() => {
	const message = useNotification()
	const { t } = useI18n()
	const formRef = ref<HTMLFormElement>()

	const { useChatAlertsQuery, useMutationUpdateChatAlerts } = useChatAlertsApi()
	const { data } = useChatAlertsQuery()
	const updateChatAlerts = useMutationUpdateChatAlerts()

	const formValue = ref<Form>({
		chatCleared: {
			enabled: false,
			messages: [],
			cooldown: 2,
		},
		cheers: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		donations: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		firstUserMessage: {
			enabled: false,
			messages: [],
			cooldown: 2,
		},
		followers: {
			enabled: false,
			messages: [],
			cooldown: 3,
		},
		raids: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		redemptions: {
			enabled: false,
			messages: [],
			cooldown: 0,
			ignoredRewardsIds: [],
		},
		streamOffline: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		streamOnline: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		subscribers: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		ban: {
			enabled: false,
			messages: [],
			cooldown: 2,
			ignoreTimeoutFrom: [],
		},
		unbanRequestCreate: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		unbanRequestResolve: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
		messageDelete: {
			enabled: false,
			messages: [],
			cooldown: 0,
		},
	})

	async function save() {
		const input = toRaw(formValue.value)
		if (!input) return

		try {
			await updateChatAlerts.executeMutation({ input })
		} catch (error) {
			message.error({
				title: t('sharedTexts.errorOnSave'),
				duration: 2500,
			})
		}
	}

	const debouncedSave = useDebounceFn(save, 1000)

	const { ignoreUpdates } = watchIgnorable(
		formValue,
		() => {
			if (!formRef.value) return
			if (!formRef.value?.reportValidity()) return

			debouncedSave()
		},
		{ deep: true, immediate: true },
	)

	watch(data, (v) => {
		ignoreUpdates(() => {
			if (!v?.chatAlerts) return
			for (const key of Object.keys(formValue.value)) {
				// eslint-disable-next-line ts/ban-ts-comment
				// @ts-expect-error
				if (!v.chatAlerts[key]) continue
				// eslint-disable-next-line ts/ban-ts-comment
				// @ts-expect-error
				formValue.value[key] = v.chatAlerts[key]
			}
		})
	}, { immediate: true })

	return {
		formValue,
		save: debouncedSave,
		formRef,
	}
})
