import type { McpRequestedScope, McpScopeGroup } from '~/utils/mcp-scopes.js'

interface GroupSelection {
	readonly read: boolean
	readonly edit: boolean
}

type ScopeSelection = Partial<Record<McpScopeGroup, GroupSelection>>
type SelectionMode = 'all' | 'read' | 'none'

function createSelection(
	scopes: readonly McpRequestedScope[],
	mode: SelectionMode,
): ScopeSelection {
	return Object.fromEntries(
		scopes.map((scope) => [
			scope.group,
			{
				read: mode !== 'none',
				edit: mode === 'all' && scope.actions.includes('edit'),
			},
		]),
	)
}

export function useMcpConsentScopeSelection() {
	const selection = ref<ScopeSelection>({})
	const showSelectionError = ref(false)

	const hasEditSelected = computed(() =>
		Object.values(selection.value).some((groupSelection) => groupSelection?.edit === true),
	)

	function initializeSelection(scopes: readonly McpRequestedScope[]): void {
		selection.value = createSelection(scopes, 'read')
		showSelectionError.value = false
	}

	function setRead(group: McpScopeGroup, value: boolean): void {
		selection.value = {
			...selection.value,
			[group]: value
				? { read: true, edit: selection.value[group]?.edit === true }
				: { read: false, edit: false },
		}
		showSelectionError.value = false
	}

	function setEdit(group: McpScopeGroup, value: boolean): void {
		selection.value = {
			...selection.value,
			[group]: value
				? { read: true, edit: true }
				: { read: selection.value[group]?.read === true, edit: false },
		}
		showSelectionError.value = false
	}

	function hasEveryRequestedActionSelected(scopes: readonly McpRequestedScope[]): boolean {
		return scopes.every((scope) => {
			const groupSelection = selection.value[scope.group]
			return (
				groupSelection?.read === true &&
				(!scope.actions.includes('edit') || groupSelection.edit === true)
			)
		})
	}

	function toggleAll(scopes: readonly McpRequestedScope[]): void {
		selection.value = createSelection(
			scopes,
			hasEveryRequestedActionSelected(scopes) ? 'none' : 'all',
		)
		showSelectionError.value = false
	}

	return {
		selection,
		showSelectionError,
		hasEditSelected,
		initializeSelection,
		setRead,
		setEdit,
		hasEveryRequestedActionSelected,
		toggleAll,
	}
}
