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
	const header = `
/**
	Code generated. DO NOT EDIT.
*/`
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

	const content = `${header}\n\n${zodImport}\n\n${generated.join('\n\n')}\n`

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

		// ---- constraints ------------------------------------------------- //
		const constraints = rawConstraint
			? rawConstraint
					.split(',')
					.map((s) => s.trim())
					.filter((s) => s !== 'omitempty')
			: []

		const hasDive = constraints.includes('dive')
		const isList =
			field.type.kind === 'ListType' ||
			(field.type.kind === 'NonNullType' && field.type.type.kind === 'ListType')

		// Length-related tags
		const lengthTags = new Set(['min', 'max', 'len', 'lte', 'gte'])
		const lengthParts = constraints.filter((part) => {
			const [tag] = part.split('=')
			return lengthTags.has(tag!)
		})

		// Parts to apply to outer validator
		const partsToApply = isList && hasDive ? lengthParts : constraints.filter((c) => c !== 'dive')

		// ---- base validator ---------------------------------------------- //
		let validator: string
		let inner: string | null = null
		if (isList) {
			const innerNode = field.type.kind === 'ListType' ? field.type.type : field.type.type
			inner = getBaseZodType(innerNode, schema)
			if (hasDive) {
				inner = applyConstraints(
					inner,
					constraints.filter((c) => c !== 'dive')
				)
			}
			validator = `z.array(${inner})`
		} else {
			validator = getBaseZodType(field.type, schema)
		}

		// Apply constraints to outer
		validator = applyConstraints(validator, partsToApply)

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
/* Helper: apply constraints to a base Zod validator string             */
/* --------------------------------------------------------------------- */
function applyConstraints(base: string, parts: string[]): string {
	let v = base
	for (const part of parts) {
		const [tag, paramStr] = part.split('=')
		const param = paramStr ?? null
		if (tag === 'min' && param !== null) {
			v = `${v}.min(${param})`
		} else if (tag === 'max' && param !== null) {
			v = `${v}.max(${param})`
		} else if (tag === 'len' && param !== null) {
			v = `${v}.length(${param})`
		} else if (tag === 'lte' && param !== null) {
			v = `${v}.max(${param})`
		} else if (tag === 'gte' && param !== null) {
			v = `${v}.min(${param})`
		} else if (tag === 'email') {
			v = `${v}.email()`
		} else if (tag === 'url') {
			v = `${v}.url()`
		} else if (tag === 'startswith' && param !== null) {
			const p = JSON.stringify(param)
			v = `${v}.startsWith(${p})`
		} else if (tag === 'endswith' && param !== null) {
			const p = JSON.stringify(param)
			v = `${v}.endsWith(${p})`
		} else if (tag === 'startsnotwith' && param !== null) {
			const p = JSON.stringify(param)
			v = `${v}.refine(val => !val.startsWith(${p}))`
		} else if (tag === 'endsnotwith' && param !== null) {
			const p = JSON.stringify(param)
			v = `${v}.refine(val => !val.endsWith(${p}))`
		} else if (tag === 'contains' && param !== null) {
			const p = JSON.stringify(param)
			v = `${v}.includes(${p})`
		} else if (tag === 'excludes' && param !== null) {
			const p = JSON.stringify(param)
			v = `${v}.refine(val => !val.includes(${p}))`
		} else if (tag === 'alpha') {
			v = `${v}.regex(/^[a-zA-Z]+$/i)`
		} else if (tag === 'required') {
			// Basic required handling: ensure non-empty for strings/arrays
			if (v.includes('z.string') || v.includes('z.array')) {
				v = `${v}.min(1)`
			}
		}
		// Add more constraints here as needed from go-playground/validator tags
	}
	return v
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
