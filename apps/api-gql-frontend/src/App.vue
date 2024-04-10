<script setup lang="ts">
import { useQuery, useMutation, useSubscription } from '@urql/vue';

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

const { data: commandSubscription } = useSubscription({
	query: graphql(`
		subscription newC {
			newCommand {
				id
			}
		}
	`),
});

const { data: notificationsSubscripction } = useSubscription({
	query: graphql(`
		subscription newN {
			newNotification {
				id
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

		<pre style="text-align: left;">
			{{ JSON.stringify(commands, null, 2) }}
		</pre>


		<h1>Example graphql subscription of new command created</h1>
		<button @click="createCommand({ name: 'test' })">
			Create command (then you need to refresh page)
		</button>
		<pre>
			{{ JSON.stringify(commandSubscription, null, 2) }}
		</pre>

		<h1>Example graphql subscription of new notifications created</h1>
		<pre>
			{{ JSON.stringify(notificationsSubscripction, null, 2) }}
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
