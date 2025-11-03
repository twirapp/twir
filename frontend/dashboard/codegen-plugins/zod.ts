// *super puper vibe-coded graphql validation -> zod schemas plugin *

import type { CodegenPlugin, Types } from '@graphql-codegen/plugin-helpers'
import type { DirectiveNode, GraphQLSchema, InputObjectTypeDefinitionNode } from 'graphql'
import { print } from 'graphql'

/**
 * Plugin configuration
 */
export interface ZodValidatorPluginConfig {
	/** Enable debug logging */
	debug?: boolean
	/** Custom import path for zod, e.g. '@/lib/zod' */
	importFrom?: string
}

export const plugin: (
	schema: GraphQLSchema,
	_documents: Types.DocumentFile[],
	config: ZodValidatorPluginConfig
) => {
	prepend: any[]
	content: string
} = (schema: GraphQLSchema, _documents: Types.DocumentFile[], config: ZodValidatorPluginConfig) => {
	const { debug = false, importFrom } = config
	const zodImport = importFrom ? `import { z } from '${importFrom}';` : `import { z } from 'zod';`

	// --------------------------------------------------------------------- //
	// 1. Find all InputObjectTypeDefinition nodes
	// --------------------------------------------------------------------- //
	const inputTypes = Object.values(schema.getTypeMap())
		.map((type) => type.astNode)
		.filter(
			(node): node is InputObjectTypeDefinitionNode => node?.kind === 'InputObjectTypeDefinition'
		)

	if (inputTypes.length === 0) {
		return { prepend: [], content: '// No input types found in schema.\n' }
	}

	// --------------------------------------------------------------------- //
	// 2. Generate Zod schema for each input
	// --------------------------------------------------------------------- //
	const generated = inputTypes.map((input) => generateZodSchema(input, schema))

	const content = `${zodImport}\n\n${generated.join('\n\n')}\n`

	if (debug) {
		console.log(
			'[zod-validator-plugin] Generated schemas for:',
			inputTypes.map((t) => t.name.value).join(', ')
		)
	}

	return { prepend: [], content }
}

/* --------------------------------------------------------------------- */
/* Helper: generate Zod object schema for a single input type            */
/* --------------------------------------------------------------------- */
function generateZodSchema(input: InputObjectTypeDefinitionNode, schema: GraphQLSchema): string {
	const name = input.name.value
	const fields = input.fields ?? []

	const zodFields = fields.map((field) => {
		const fieldName = field.name.value

		// ---- optionality ------------------------------------------------- //
		const isNonNull = field.type.kind === 'NonNullType'
		const hasQuestionMark = field.type.kind === 'NamedType' && field.type.name.value.endsWith('?')
		const rawConstraint = getValidateConstraint(field)
		const hasOmitempty = rawConstraint?.includes('omitempty') ?? false
		const isOptional = !isNonNull || hasQuestionMark || hasOmitempty

		// ---- base Zod type ----------------------------------------------- //
		let validator = getBaseZodType(field.type, schema)

		// ---- apply constraints (except omitempty) ------------------------ //
		const constraints = rawConstraint
			? rawConstraint
					.split(',')
					.map((s) => s.trim())
					.filter((s) => s !== 'omitempty')
			: []

		// apply simple constraints
		for (const part of constraints) {
			if (part.startsWith('max=')) {
				const val = part.split('=')[1]
				validator = `${validator}.max(${val})`
			} else if (part.startsWith('min=')) {
				const val = part.split('=')[1]
				validator = `${validator}.min(${val})`
			} else if (part.startsWith('lte=')) {
				const val = part.split('=')[1]
				validator = `${validator}.max(${val})`
			}
		}

		const isList =
			field.type.kind === 'ListType' ||
			(field.type.kind === 'NonNullType' && field.type.type.kind === 'ListType')

		if (isList && constraints.includes('dive')) {
			const innerNode = field.type.kind === 'ListType' ? field.type.type : field.type.type

			let inner = getBaseZodType(innerNode, schema)
			const diveConstraints = constraints.filter((c) => c !== 'dive')

			for (const part of diveConstraints) {
				if (part.startsWith('max=')) inner = `${inner}.max(${part.split('=')[1]})`
				if (part.startsWith('min=')) inner = `${inner}.min(${part.split('=')[1]})`
				if (part.startsWith('lte=')) inner = `${inner}.max(${part.split('=')[1]})`
			}

			validator = `z.array(${inner})`
		}

		// ---- make optional ------------------------------------------------ //
		if (isOptional) {
			validator = `z.optional(${validator})`
		}

		return `  ${fieldName}: ${validator},`
	})

	return `export const ${name}Schema = z.object({\n${zodFields.join(
		'\n'
	)}\n});\n\nexport type ${name}Input = z.infer<typeof ${name}Schema>;`
}

/* --------------------------------------------------------------------- */
/* Helper: extract @validate(constraint: "...") value                    */
/* --------------------------------------------------------------------- */
function getValidateConstraint(field: any): string | null {
	const dir: DirectiveNode | undefined = field.directives?.find(
		(d: any) => d.name.value === 'validate'
	)
	if (!dir) return null

	const arg = dir.arguments?.find((a: any) => a.name.value === 'constraint')
	return arg?.value?.kind === 'StringValue' ? arg.value.value : null
}

/* --------------------------------------------------------------------- */
/* Helper: map GraphQL type → Zod base type (string)                     */
/* --------------------------------------------------------------------- */
function getBaseZodType(typeNode: any, schema: GraphQLSchema): string {
	// NonNull → recurse
	if (typeNode.kind === 'NonNullType') {
		return getBaseZodType(typeNode.type, schema)
	}

	// List
	if (typeNode.kind === 'ListType') {
		const inner = getBaseZodType(typeNode.type, schema)
		return `z.array(${inner})`
	}

	// Named type
	if (typeNode.kind === 'NamedType') {
		const rawName = typeNode.name.value.replace(/\?$/, '') // strip ?
		const scalarMap: Record<string, string> = {
			String: 'z.string()',
			Int: 'z.number().int()',
			Float: 'z.number()',
			Boolean: 'z.boolean()',
			ID: 'z.string()',
			UUID: 'z.string().uuid()',
			Time: 'z.number()',
			Upload: 'z.instanceof(File)',
		}

		if (scalarMap[rawName]) return scalarMap[rawName]

		// Enum
		const gqlType = schema.getType(rawName)
		if (gqlType?.astNode?.kind === 'EnumTypeDefinition') {
			const values = gqlType.astNode.values!.map((v: any) => `'${v.name.value}'`).join(', ')
			return `z.enum([${values}])`
		}

		// Nested input → lazy
		return `z.lazy(() => ${rawName}Schema)`
	}

	return 'z.any()'
}
