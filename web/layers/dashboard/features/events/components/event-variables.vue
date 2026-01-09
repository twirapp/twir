<script setup lang="ts">
import { useField } from 'vee-validate'
import { computed } from 'vue'


import { flatEvents } from '../constants/helpers.ts'

import type { EventType } from '~/gql/graphql.ts'




const { value } = useField<EventType>('type')
const { t } = useI18n()

const event = computed(() => {
	return flatEvents[value.value]
})
</script>

<template>
	<UiCard>
		<UiCardHeader>
			<UiCardTitle> Variables </UiCardTitle>
		</UiCardHeader>
		<UiCardContent>
			<UiAlert v-if="!event.variables?.length">
				No variables available for this event type.
			</UiAlert>
			<div v-else>
				<p class="text-sm mb-2">
					Variables for <strong>{{ event.name }}</strong> event:
				</p>
				<ul class="list-disc pl-5 space-y-2">
					<li v-for="variable in event.variables" :key="variable">
						<span class="font-bold bg-zinc-800 p-1 rounded">
							{{ '{' + `${variable}` + '}' }}
						</span>
						-
						<span>
							{{ t(`events.variables.${variable}`) ?? 'No description' }}
						</span>
					</li>
				</ul>
			</div>
		</UiCardContent>
	</UiCard>
</template>
