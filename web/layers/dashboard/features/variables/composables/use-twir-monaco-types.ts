import { computed } from 'vue'

import { useSecretsApi } from '~~/layers/dashboard/features/channels-secret/composables/use-secrets-api'

export function useTwirMonacoTypes() {
	const secretsApi = useSecretsApi()
	const secrets = secretsApi.secrets

	const twirTypeDefinitions = computed(() => {
		const secretNames = secrets.value.map((s) => s.name)

		const secretsOverload = secretNames.length > 0
			? secretNames.map((name) => `        get(name: '${name}'): string | null;`).join('\n')
			: '        get(name: string): string | null;'

		return `
interface TwirChannel {
    /** Current channel ID */
    id: string;
}

interface TwirSecrets {
${secretsOverload}
}

interface Twir {
    secrets: TwirSecrets;
    channel: TwirChannel;
}

declare const twir: Twir;
`
	})

	return { twirTypeDefinitions }
}
