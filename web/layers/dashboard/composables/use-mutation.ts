import { useMutation as _useMutation } from '@urql/vue'

import type { AnyVariables, TypedDocumentNode } from '@urql/vue'
import type { DocumentNode } from 'graphql'

export function useMutation<T = any, V extends AnyVariables = AnyVariables>(
	query: TypedDocumentNode<T, V> | DocumentNode | string,
	additionalTypenames?: string[]
) {
	const mutation = _useMutation(query)

	if (additionalTypenames) {
		const originalMutation = mutation.executeMutation
		const newMutation: typeof originalMutation = (values) => {
			return originalMutation(values, { additionalTypenames })
		}
		mutation.executeMutation = newMutation
	}

	return mutation
}
