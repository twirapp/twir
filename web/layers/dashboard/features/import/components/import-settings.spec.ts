import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import { ImportFailureReason } from '~/gql/graphql.js'
import { createImportsTestI18n } from './test-i18n.js'
import ImportSettings from './import-settings.vue'

const i18n = createImportsTestI18n()

const stubs = {
	Alert: { template: '<div><slot /></div>' },
	AlertDescription: { template: '<p><slot /></p>' },
	AlertTitle: { template: '<strong><slot /></strong>' },
	Badge: { template: '<span><slot /></span>' },
	Button: { props: ['disabled'], template: '<button :disabled="disabled"><slot /></button>' },
	Card: { template: '<section><slot /></section>' },
	CardContent: { template: '<div><slot /></div>' },
	CardDescription: { template: '<p><slot /></p>' },
	CardFooter: { template: '<footer><slot /></footer>' },
	CardHeader: { template: '<header><slot /></header>' },
	CardTitle: { template: '<h3><slot /></h3>' },
	NuxtIcon: true,
}

describe('ImportSettings', () => {
	it('keeps command and timer actions independent', async () => {
		const wrapper = mount(ImportSettings, {
			props: {
				connected: true,
				canManageCommands: true,
				canManageTimers: false,
			},
			global: { plugins: [i18n], stubs },
		})

		await wrapper.get('[data-test="import-commands"]').trigger('click')
		expect(wrapper.emitted('importCommands')).toHaveLength(1)
		expect(wrapper.get('[data-test="import-timers"]').attributes('disabled')).toBeDefined()
	})

	it('renders partial results and translated stable failure reasons', () => {
		const wrapper = mount(ImportSettings, {
			props: {
				connected: true,
				canManageCommands: true,
				canManageTimers: true,
				commandsReport: {
					importedCount: 2,
					failedCount: 1,
					failures: [{ name: '!hello', reason: ImportFailureReason.Duplicate }],
				},
			},
			global: { plugins: [i18n], stubs },
		})

		expect(wrapper.text()).toContain('2 imported')
		expect(wrapper.text()).toContain('1 failed')
		expect(wrapper.text()).toContain('!hello')
		expect(wrapper.text()).toContain('Already exists')
	})
})
