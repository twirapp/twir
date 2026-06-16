<script setup lang="ts">
import { useField } from 'vee-validate'
import { computed } from 'vue'

import type { EventType } from '~/gql/graphql.js'

import { Alert } from '@/components/ui/alert'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

import { flatEvents } from '../constants/helpers.js'

const { value } = useField<EventType>('type')
const { t } = useI18n()

const event = computed(() => {
	return flatEvents[value.value]
})
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle> Variables </CardTitle>
		</CardHeader>
		<CardContent>
			<Alert v-if="!event.variables?.length"> No variables available for this event type. </Alert>
			<div v-else>
				<p class="mb-2 text-sm">
					Variables for <strong>{{ event.name }}</strong> event:
				</p>
				<ul class="list-disc space-y-2 pl-5">
					<li
						v-for="variable in event.variables"
						:key="variable"
					>
						<span class="rounded bg-zinc-800 p-1 font-bold">
							{{ '{' + `${variable}` + '}' }}
						</span>
						-
						<span>
							{{ t(`events.variables.${variable}`) ?? 'No description' }}
						</span>
					</li>
				</ul>
			</div>
		</CardContent>
	</Card>
</template>
