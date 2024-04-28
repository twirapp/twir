import { useQuery } from '@urql/vue';
import { defineStore } from 'pinia';
import { SetOptional } from 'type-fest';
import { computed } from 'vue';

import { useMutation } from '@/composables/use-mutation';
import { graphql } from '@/gql';
import { GetCustomAndBuiltInVariablesQuery, VariableType } from '@/gql/graphql';

const invalidationKey = 'VariablesInvalidateKey';

export type CustomVariable = GetCustomAndBuiltInVariablesQuery['variables'][number]
export type EditableCustomVariable = SetOptional<CustomVariable, 'id'>

export const useVariablesApi = defineStore('api/variables', () => {
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
				}
				variablesBuiltIn {
					name
					example
					description
					visible
					canBeUsedInRegistry
				}
			}
		`),
	});

	const customVariables = computed(() => {
		const mapped = variablesQuery.data.value?.variables.map((variable) => ({
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
		})) ?? [];

		return mapped;
	});

	const builtInVariables = computed(() => {
		const mapped = variablesQuery.data.value?.variablesBuiltIn.map((variable) => ({
			name: variable.name,
			description: variable.description,
			visible: variable.visible,
			example: variable.example || `${variable.name}`,
			isBuiltIn: true,
			canBeUsedInRegistry: variable.canBeUsedInRegistry,
		})) ?? [];

		return mapped;
	});

	const allVariables = computed(() => {
		return [
			...customVariables.value,
			...builtInVariables.value,
		];
	});

	const isLoading = computed(() => {
		return variablesQuery.fetching.value;
	});

	const useMutationCreateVariable = () => useMutation(graphql(`
		mutation CreateVariable($opts: VariableCreateInput!) {
			variablesCreate(opts: $opts) {
				id
			}
		}
	`), [invalidationKey]);

	const useMutationUpdateVariable = () => useMutation(graphql(`
		mutation UpdateVariable($id: ID!, $opts: VariableUpdateInput!) {
			variablesUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`), [invalidationKey]);

	const useMutationRemoveVariable = () => useMutation(graphql(`
		mutation RemoveVariable($id: ID!) {
			variablesDelete(id: $id)
		}
	`), [invalidationKey]);

	return {
		customVariables,
		builtInVariables,
		allVariables,
		isLoading,
		useMutationCreateVariable,
		useMutationUpdateVariable,
		useMutationRemoveVariable,
	};
});
