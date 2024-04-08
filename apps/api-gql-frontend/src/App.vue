<script setup lang="ts">
import { useQuery, useMutation } from '@urql/vue';
import { ref } from 'vue';

import { graphql } from './gql';

const query = ref({ userId: '2' });

const createCommandMutation = graphql(`
	mutation newCommand($name: String!, $description: String) {
		createCommand(name: $name, description: $description) {
			id
			name
		}
	}
`);

const createNotificationMutation = graphql(`
	mutation newNotification($text: String!, $userId: String) {
		createNotification(text: $text, userId: $userId) {
			id
			userId
			text
		}
	}
`);

const { executeMutation: createCommand } = useMutation(createCommandMutation);
const { executeMutation: createNotification } = useMutation(createNotificationMutation);

const { data } = useQuery({
		query: graphql(`
			query getCommandsAndNotifications($userId: String!) {
				commands {
					id
					name
					aliases
				}
				notifications(userId: $userId) {
					id
					userId
					text
				}
			}
  `),
	variables: query,
});

</script>

<template>
	<div style="display: flex; flex-direction: column; gap: 12px;">
		<button @click="createCommand({ name: 'test'})">Create command</button>
		<button @click="createNotification({  text: 'im new', userId: '2'})">Create notification</button>
		<pre style="text-align: left;">
			{{ JSON.stringify(data, null, 2) }}
		</pre>
	</div>
</template>

<style scoped>
.logo {
	height: 6em;
	padding: 1.5em;
	will-change: filter;
	transition: filter 300ms;
}

.logo:hover {
	filter: drop-shadow(0 0 2em #646cffaa);
}

.logo.vue:hover {
	filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
