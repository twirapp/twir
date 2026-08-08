import { describe, expect, it } from 'vitest'

import { isEditableHotkeyTarget } from './composables/isEditableHotkeyTarget'

// Regression: the canvas window-level Backspace/Delete preventDefault must not
// fire for editable targets, otherwise typing in fields like the font search
// (teleported popover inputs without @keydown.stop) silently breaks.
describe('isEditableHotkeyTarget', () => {
	it('returns true for input, textarea and contentEditable targets', () => {
		const input = document.createElement('input')
		const textarea = document.createElement('textarea')
		const editable = document.createElement('div')
		editable.contentEditable = 'true'

		expect(isEditableHotkeyTarget(input)).toBe(true)
		expect(isEditableHotkeyTarget(textarea)).toBe(true)
		expect(isEditableHotkeyTarget(editable)).toBe(true)
	})

	it('returns false for non-editable targets and null', () => {
		const div = document.createElement('div')

		expect(isEditableHotkeyTarget(div)).toBe(false)
		expect(isEditableHotkeyTarget(document.body)).toBe(false)
		expect(isEditableHotkeyTarget(null)).toBe(false)
	})
})
