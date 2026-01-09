import type { SetOptional } from 'type-fest'

import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import type { GetCustomAndBuiltInVariablesQuery } from '@/gql/graphql.js'

import { graphql } from '@/gql/gql.js'
import { VariableType } from '@/gql/graphql.js'
import { useMutation } from '~/composables/use-mutation.js'

const invalidationKey = 'VariablesInvalidateKey'

export type CustomVariable = GetCustomAndBuiltInVariablesQuery['variables'][number]
export type EditableCustomVariable = Omit<SetOptional<CustomVariable, 'id'>, '__typename'>

export const useVariablesApi = createGlobalState(() => {
	const variablesQuery = useQuery({
		variables: {},
		context: { additionalTypenames: [invalidationKey] },
		query: graphql(`
			query GetCustomAndBuiltInVariables {
				variables {
					id
					description
					type
					name
					evalValue
					response
					scriptLanguage
				}
				variablesBuiltIn {
					name
					example
					description
					visible
					canBeUsedInRegistry
					links {
						href
						name
					}
				}
			}
		`),
	})

	const customVariables = computed(() => {
		const mapped =
			variablesQuery.data.value?.variables.map((variable) => ({
				id: variable.id,
				name: variable.name,
				description: variable.description,
				visible: true,
				example: `customvar|${variable.name}`,
				isBuiltIn: false,
				canBeUsedInRegistry: variable.type !== VariableType.Script,
				type: variable.type,
				response: variable.response,
				evalValue: variable.evalValue,
				scriptLanguage: variable.scriptLanguage,
				links: [],
			})) ?? []

		return mapped
	})

	const builtInVariables = computed(() => {
		const mapped =
			variablesQuery.data.value?.variablesBuiltIn.map((variable) => ({
				name: variable.name,
				description: variable.description,
				visible: variable.visible,
				example: variable.example || `${variable.name}`,
				isBuiltIn: true,
				canBeUsedInRegistry: variable.canBeUsedInRegistry,
				links: variable.links,
			})) ?? []

		return mapped
	})

	const allVariables = computed(() => {
		return [...customVariables.value, ...builtInVariables.value]
	})

	const isLoading = computed(() => {
		return variablesQuery.fetching.value
	})

	const useMutationCreateVariable = () =>
		useMutation(
			graphql(`
		mutation CreateVariable($opts: VariableCreateInput!) {
			variablesCreate(opts: $opts) {
				id
			}
		}
	`),
			[invalidationKey]
		)

	const useMutationUpdateVariable = () =>
		useMutation(
			graphql(`
		mutation UpdateVariable($id: UUID!, $opts: VariableUpdateInput!) {
			variablesUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`),
			[invalidationKey]
		)

	const useMutationRemoveVariable = () =>
		useMutation(
			graphql(`
		mutation RemoveVariable($id: UUID!) {
			variablesDelete(id: $id)
		}
	`),
			[invalidationKey]
		)

	const useMutationExecuteScript = () =>
		useMutation(
			graphql(`
		mutation ExecuteScript($expression: String!, $language: VariableScriptLanguage!, $testFromUserName: String) {
			executeScript(script: $expression, language: $language, testAsUserName: $testFromUserName)
		}
	`)
		)

	return {
		variablesQuery,
		customVariables,
		builtInVariables,
		allVariables,
		isLoading,
		useMutationCreateVariable,
		useMutationUpdateVariable,
		useMutationRemoveVariable,
		useMutationExecuteScript,
	}
})
