<script setup lang="ts">
import { useQuery, useMutation } from '@urql/vue';

import { graphql } from './gql';

const createCommandMutation = graphql(`
	mutation newCommand($name: String!, $aliases: [String!], $description: String, $responses: [CreateCommandResponseInput!]) {
		createCommand(
    opts: {name: $name, description: $description, aliases: $aliases, responses: $responses}
  ) {
    id
  }
	}
`);

const { executeMutation: createCommand } = useMutation(createCommandMutation);

const { data: commands } = useQuery({
	query: graphql(`
			query getCommands {
				commands {
					name
					id
					aliases
					responses {
						text
					}
					description
					createdAt
					updatedAt
				}
			}
  `),
});

const { data: user } = useQuery({
	query: graphql(`
		query getUser {
			authedUser {
			id
			apiKey
			channel {
				botId
				isBotModerator
				isEnabled
			}
			hideOnLandingPage
			isBanned
			isBotAdmin
			}
		}
	`),
});
</script>

<template>
	<div style="display: flex; flex-direction: column; gap: 12px;">
		<h1>Authed as</h1>
		<pre style="text-align: left;">
			{{ JSON.stringify(user, null, 2) }}
		</pre>

		<button @click="createCommand({ name: 'test' })">
			Create command (then you need to refresh page)
		</button>
		<pre style="text-align: left;">
			{{ JSON.stringify(commands, null, 2) }}
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
