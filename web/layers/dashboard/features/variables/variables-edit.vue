<script lang="ts" setup>
import { VueMonacoEditor, useMonaco } from '@guolao/vue-monaco-editor'
import { Label } from 'reka-ui'
import { useForm } from 'vee-validate'
import { onMounted, ref, toRaw, watch } from 'vue'
import { useRoute } from 'vue-router'
import VariablesList from '~~/layers/dashboard/components/variables-list.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import Button from '@/components/ui/button/Button.vue'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { VariableScriptLanguage, VariableType } from '~/gql/graphql.js'

import { formSchema, useVariablesEdit } from './composables/use-variables-edit'

const { monacoRef } = useMonaco()

const TWIR_TYPE_DEFS = `
interface TwirSecrets {
	/** Get a secret value by name. Returns null if not found. */
	get(name: string): string | null;
}

interface TwirChannel {
	/** Current channel ID */
	id: string;
}

interface Twir {
	/** Access channel secrets stored in the dashboard */
	secrets: TwirSecrets;
	/** Information about the current channel */
	channel: TwirChannel;
}

/** Twir API available in custom variable scripts */
declare const twir: Twir;

/** Fetch API available in custom variable scripts */
declare function fetch(url: string, options?: RequestInit): Promise<Response>;
`

watch(monacoRef, (monaco) => {
	if (!monaco) return

	monaco.languages.typescript.javascriptDefaults.setCompilerOptions({
		target: monaco.languages.typescript.ScriptTarget.ESNext,
		module: monaco.languages.typescript.ModuleKind.ESNext,
		allowNonTsExtensions: true,
	})

	monaco.languages.typescript.javascriptDefaults.addExtraLib(
		TWIR_TYPE_DEFS,
		'twir://globals.d.ts',
	)
}, { immediate: true })

const route = useRoute<'dashboard-variables-id'>()
const { t } = useI18n()
const { findVariable, submit, runScript, testFromUserName } = useVariablesEdit()

const loading = ref(true)
const title = ref('')

const jsExample = `// semicolons (;) matters, do not forget put them on end of statements.
// you can use commands variables:
// const userFollowAge = '$(user.followage)'

const request = await fetch('https://jsonplaceholder.typicode.com/todos/1');
const response = await request.json();
// you should return value from your script
return response.title;
`

const { handleSubmit, setValues, values } = useForm({
	validationSchema: formSchema,
	initialValues: {
		description: null,
		type: VariableType.Text,
		response: '',
		evalValue: jsExample,
		scriptLanguage: VariableScriptLanguage.Javascript,
	},
	keepValuesOnUnmount: true,
})

onMounted(async () => {
	if (typeof route.params.id !== 'string') {
		return
	}

	const variable = await findVariable(route.params.id)
	if (variable) {
		setValues(toRaw(variable))
		title.value = variable.name
	}

	loading.value = false
})

watch(
	() => values.scriptLanguage,
	(newLanguage) => {
		if (newLanguage === VariableScriptLanguage.Javascript) {
			setValues({ evalValue: jsExample })
		}
	}
)

const onSubmit = handleSubmit(submit)

const executionResult = ref('')

async function executeScript() {
	if (!values.evalValue) {
		return
	}

	executionResult.value = 'Executing...'

	try {
		const result = await runScript(values.evalValue, values.scriptLanguage)
		if (result) {
			executionResult.value = result
		}
	} catch (error: any) {
		if (('message' in error) as any) {
			executionResult.value = error.message
		}
	}
}
</script>

<template>
	<form
		:class="{ 'blur-xs': loading }"
		@submit="onSubmit"
	>
		<PageLayout
			stickyHeader
			show-back
			back-redirect-to="/dashboard/variables"
		>
			<template #title>
				<span v-if="route.params.id === 'create'">Create</span>
				<span v-else>Edit "{{ title }}"</span>
			</template>

			<template #action>
				<Button
					type="submit"
					:loading="loading"
				>
					{{ t('sharedButtons.save') }}
				</Button>
			</template>

			<template #content>
				<div class="flex h-full flex-col gap-4">
					<FormField
						v-slot="{ componentField }"
						name="name"
					>
						<FormItem>
							<FormLabel>{{ t('sharedTexts.name') }}</FormLabel>
							<FormControl>
								<Input v-bind="componentField" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField
						v-slot="{ componentField }"
						name="type"
					>
						<FormItem>
							<FormLabel>{{ t('variables.type') }}</FormLabel>

							<Select v-bind="componentField">
								<FormControl>
									<SelectTrigger>
										<SelectValue placeholder="Select a verified email to display" />
									</SelectTrigger>
								</FormControl>
								<SelectContent>
									<SelectGroup>
										<SelectItem
											v-for="variable of VariableType"
											:key="variable"
											:value="variable"
										>
											{{ variable }}
										</SelectItem>
									</SelectGroup>
								</SelectContent>
							</Select>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField
						v-if="values.type !== VariableType.Script"
						v-slot="{ componentField }"
						name="response"
					>
						<FormItem>
							<FormLabel>{{ t('sharedTexts.response') }}</FormLabel>
							<FormControl>
								<Input
									type="text"
									v-bind="componentField"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<div
						v-show="values.type === VariableType.Script"
						class="flex flex-col gap-2"
					>
						<FormField
							v-slot="{ componentField }"
							name="scriptLanguage"
						>
							<FormItem>
								<FormLabel>Script Language</FormLabel>
								<Select v-bind="componentField">
									<FormControl>
										<SelectTrigger>
											<SelectValue placeholder="Select script language" />
										</SelectTrigger>
									</FormControl>
									<SelectContent>
										<SelectGroup>
											<SelectItem
												v-for="language in Object.values(VariableScriptLanguage)"
												:key="language"
												:value="language"
											>
												{{ language }}
											</SelectItem>
										</SelectGroup>
									</SelectContent>
								</Select>
								<FormMessage />
							</FormItem>
						</FormField>

						<span>Execution result</span>
						<div class="flex flex-row gap-2">
							<div class="bg-secondary h-auto w-full rounded-md p-2">
								{{ executionResult || 'Run a script for test your code' }}
							</div>
							<Button
								type="button"
								class="place-self-start"
								@click="executeScript"
							>
								<Icon
									name="lucide:terminal"
									class="mr-2 size-4"
								/>
								Run
							</Button>
						</div>
						<div class="flex flex-col gap-2">
							<Label for="testFromUserName">Test as specific viewer</Label>
							<Input
								id="testFromUserName"
								v-model:model-value="testFromUserName"
								placeholder="Enter username from which perspective script will run"
							/>
						</div>

						<Alert>
							<Icon
								name="lucide:info"
								class="size-4"
							/>
							<AlertTitle>Heads up!</AlertTitle>
							<AlertDescription class="flex flex-col items-start justify-start gap-2">
								<span>
									You can use variables as you doing it in commands, like
									<code class="text-teal-200">$(user.followage)</code>. They will be parsed and
									evaluated. But you must enclose them in quotes for proper usage!
								</span>

								<VariablesList>
									<template #trigger>
										<Button
											type="button"
											size="sm"
										>
											Show variables list
										</Button>
									</template>
								</VariablesList>
							</AlertDescription>
						</Alert>

						<VueMonacoEditor
							:value="values.evalValue"
							:language="values.scriptLanguage!.toLowerCase()"
							class="h-full min-h-[500px]"
							theme="vs-dark"
							@change="setValues({ evalValue: $event })"
						/>
					</div>
				</div>
			</template>
		</PageLayout>
	</form>
</template>
