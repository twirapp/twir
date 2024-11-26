import crypto from 'node:crypto'

import { Kind, print, visit } from 'graphql'

import type { DocumentNode } from 'graphql'

export function generatePersistHash(documentNode: DocumentNode) {
	const sanitizedDocument = visit(documentNode, {
		[Kind.FIELD](field) {
			if (field.directives?.some(directive => directive.name.value === 'client')) {
				return null
			}
		},
		[Kind.DIRECTIVE](directive) {
			if (directive.name.value === 'connection') {
				return null
			}
		},
	})
	const documentStr = print(sanitizedDocument)

	return {
		hash: computeQueryHash(documentStr),
	}
}

function computeQueryHash(query: string) {
	const hash = crypto.createHash('sha256')
	hash.update(query)
	return hash.digest('hex')
}
