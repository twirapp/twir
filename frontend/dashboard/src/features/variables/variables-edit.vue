<script lang="ts" setup>
import { VueMonacoEditor } from '@guolao/vue-monaco-editor'
import { toTypedSchema } from '@vee-validate/zod'
import { InfoIcon, TerminalIcon } from 'lucide-vue-next'
import { Label } from 'radix-vue'
import { useForm } from 'vee-validate'
import { onMounted, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

import { formSchema, useVariablesEdit } from './composables/use-variables-edit'

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
import VariablesList from '@/components/variables-list.vue'
import { VariableScriptLanguage, VariableType } from '@/gql/graphql'
import PageLayout from '@/layout/page-layout.vue'

const route = useRoute()
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

const pythonExample = `import urllib.request

url = "https://jsonplaceholder.typicode.com/todos/1"
try:
    with urllib.request.urlopen(url, timeout=2) as response:
        return response.read().decode()
except urllib.error.URLError as e:
    return "Request failed: " + str(e)
`

const { handleSubmit, setValues, values } = useForm({
	validationSchema: toTypedSchema(formSchema),
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

watch(() => values.scriptLanguage, (newLanguage) => {
	if (newLanguage === VariableScriptLanguage.Python) {
		setValues({ evalValue: pythonExample })
	} else if (newLanguage === VariableScriptLanguage.Javascript) {
		setValues({ evalValue: jsExample })
	}
})

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
		if ('message' in error as any) {
			executionResult.value = error.message
		}
	}
}
</script>

<template>
	<form :class="{ 'blur-xs': loading }" @submit="onSubmit">
		<PageLayout stickyHeader show-back back-redirect-to="/dashboard/variables">
			<template #title>
				<span v-if="route.params.id === 'create'">Create</span>
				<span v-else>Edit "{{ title }}"</span>
			</template>

			<template #action>
				<Button type="submit" :loading="loading">
					{{ t('sharedButtons.save') }}
				</Button>
			</template>

			<template #content>
				<div class="flex flex-col gap-4 h-full">
					<FormField v-slot="{ componentField }" name="name">
						<FormItem>
							<FormLabel>{{ t('sharedTexts.name') }}</FormLabel>
							<FormControl>
								<Input
									v-bind="componentField"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="type">
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
										<SelectItem v-for="variable of VariableType" :key="variable" :value="variable">
											{{ variable }}
										</SelectItem>
									</SelectGroup>
								</SelectContent>
							</Select>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField
						v-if="values.type !== VariableType.Script" v-slot="{ componentField }"
						name="response"
					>
						<FormItem>
							<FormLabel>{{ t('sharedTexts.response') }}</FormLabel>
							<FormControl>
								<Input type="text" v-bind="componentField" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<div v-show="values.type === VariableType.Script" class="flex flex-col gap-2">
						<FormField v-slot="{ componentField }" name="scriptLanguage">
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
							<div class="bg-secondary rounded-md h-auto p-2 w-full">
								{{ executionResult || 'Run a script for test your code' }}
							</div>
							<Button type="button" class="place-self-start" @click="executeScript">
								<TerminalIcon class="size-4 mr-2" />
								Run
							</Button>
						</div>
						<div class="flex flex-col gap-2">
							<Label for="testFromUserName">Test as specific viewer</Label>
							<Input
								id="testFromUserName" v-model:model-value="testFromUserName"
								placeholder="Enter username from which perspective script will run"
							/>
						</div>

						<Alert>
							<InfoIcon class="size-4" />
							<AlertTitle>Heads up!</AlertTitle>
							<AlertDescription class="flex flex-col justify-start items-start gap-2">
								<span>
									You can use variables as you doing it in commands, like <code
										class="text-teal-200"
									>$(user.followage)</code>.
									They will be parsed and evaluated.
									But you must enclose them in quotes for proper usage!
								</span>

								<VariablesList>
									<template #trigger>
										<Button type="button" size="sm">
											Show variables list
										</Button>
									</template>
								</VariablesList>
							</AlertDescription>
						</Alert>

						<VueMonacoEditor
							:value="values.evalValue"
							:language="values.scriptLanguage!.toLowerCase()"
							class="min-h-[500px] h-full"
							theme="vs-dark"
							@change="setValues({ evalValue: $event })"
						/>
					</div>
				</div>
			</template>
		</PageLayout>
	</form>
</template>
