import { createGlobalState } from '@vueuse/core'
import { ref, unref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { nativeEnum, object, string } from 'zod'

import type { CustomVariable } from '@/api/variables.js'
import type { MaybeRef } from 'vue'
import type { TypeOf } from 'zod'

import { useVariablesApi } from '@/api/variables.js'
import { useToast } from '@/components/ui/toast'
import { VariableScriptLanguage, VariableType } from '@/gql/graphql.js'

export const formSchema = object({
	id: string().optional(),
	name: string().min(1).max(50),
	description: string().max(500).nullable().optional(),
	type: nativeEnum(VariableType),
	scriptLanguage: nativeEnum(VariableScriptLanguage).default(VariableScriptLanguage.Javascript),
}).and(object({
	response: string().max(500),
	evalValue: string().max(5000),
}).refine((data) => {
	if (!data.response && !data.evalValue) {
		return false
	}

	return true
}, {
	message: 'Script or response must be specified',
	path: ['response', 'evalValue'],
}))

export type FormSchema = TypeOf<typeof formSchema>

export const useVariablesEdit = createGlobalState(() => {
	const { toast } = useToast()
	const { t } = useI18n()
	const router = useRouter()

	const variablesApi = useVariablesApi()
	const update = variablesApi.useMutationUpdateVariable()
	const create = variablesApi.useMutationCreateVariable()
	const scriptExecutor = variablesApi.useMutationExecuteScript()

	const variable = ref<CustomVariable | null>(null)
	const testFromUserName = ref('')

	async function findVariable(id: string) {
		variable.value = null
		if (id === 'create') return

		const fetchedData = await variablesApi.variablesQuery.then((variables) => variables)
		const found = fetchedData.data?.value?.variables.find((variable) => variable.id === id)

		if (!found) {
			throw new Error('Command not found')
		}

		variable.value = found

		return found
	}

	async function submit(data: FormSchema) {
		console.log(data)
		if (data.id) {
			await update.executeMutation({
				id: data.id,
				opts: {
					...data,
					// eslint-disable-next-line ts/ban-ts-comment
					// @ts-expect-error
					id: undefined,
					description: '',
				},
			})
		} else {
			const result = await create.executeMutation({
				opts: {
					...data,
					description: '',
				},
			})

			if (result.error) {
				toast({
					title: result.error.graphQLErrors?.map(e => e.message).join(', ') ?? 'error',
					duration: 5000,
					variant: 'destructive',
				})
				return
			}

			await router.push({
				path: `/dashboard/variables/${result.data?.variablesCreate.id}`,
				state: {
					noTransition: true,
				},
			})
		}

		toast({
			title: t('sharedTexts.saved'),
			duration: 2500,
			variant: 'success',
		})
	}

	async function runScript(expression: MaybeRef<string>, language: MaybeRef<VariableScriptLanguage> = VariableScriptLanguage.Javascript) {
		const result = await scriptExecutor.executeMutation({
			expression: unref(expression),
			testFromUserName: unref(testFromUserName),
			language: unref(language),
		})

		return result.data?.executeScript ?? result.error?.message
	}

	return {
		findVariable,
		submit,
		runScript,
		testFromUserName,
	}
})
