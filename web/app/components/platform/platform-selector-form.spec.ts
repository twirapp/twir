import { mount } from '@vue/test-utils'
import { Field as FormField, useForm } from 'vee-validate'
import { describe, expect, it } from 'vitest'
import { defineComponent, h } from 'vue'

import { formSchema } from '~~/layers/dashboard/features/commands/composables/use-command-edit-v2'

import PlatformSelector from './platform-selector.vue'

const Host = defineComponent({
	setup() {
		useForm({
			validationSchema: formSchema,
			initialValues: {
				enabled: true,
				aliases: [],
				responses: [],
				description: '',
				rolesIds: [],
				deniedUsersIds: [],
				allowedUsersIds: [],
				requiredMessages: 0,
				requiredUsedChannelPoints: 0,
				requiredWatchTime: 0,
				cooldown: 0,
				cooldownType: 'GLOBAL',
				isReply: true,
				visible: true,
				keepResponsesOrder: true,
				onlineOnly: false,
				offlineOnly: false,
				enabledCategories: [],
				expiresType: null,
				expiresAt: null,
				roleCooldowns: [],
			},
			keepValuesOnUnmount: true,
		})

		return () =>
			h(
				FormField,
				{ name: 'platforms' },
				{
					default: ({ field }: any) =>
						h(PlatformSelector, {
							'modelValue': field.value,
							'onUpdate:modelValue': field['onUpdate:modelValue'],
						}),
				},
			)
	},
})

const iconStub = { NuxtIcon: { props: ['name'], template: '<i :data-icon="name" />' } }

describe('PlatformSelector inside commands create form', () => {
	it('applies selection on click', async () => {
		const wrapper = mount(Host, { global: { stubs: iconStub } })

		const twitch = wrapper.findAll('button').find((b) => b.text().includes('Twitch'))!
		await twitch.trigger('click')

		expect(twitch.attributes('data-active')).toBe('true')
	})
})
