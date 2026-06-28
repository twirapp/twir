import { useQuery } from "@urql/vue";
import { createGlobalState } from "@vueuse/core";
import { computed } from "vue";

import { useMutation } from "@/composables/use-mutation.js";
import { graphql } from "@/gql/gql.js";

const invalidationKey = "SecretsInvalidateKey";

export type Secret = {
	id: string;
	name: string;
	description?: string | null;
	value: string;
};

export const useSecretsApi = createGlobalState(() => {
	const secretsQuery = useQuery({
		variables: {},
		context: { additionalTypenames: [invalidationKey] },
		query: graphql(`
			query GetSecrets {
				secrets {
					id
					name
					description
					value
				}
			}
		`),
	});

	const secrets = computed(() => {
		return secretsQuery.data.value?.secrets ?? [];
	});

	const isLoading = computed(() => {
		return secretsQuery.fetching.value;
	});

	const useMutationCreateSecret = () =>
		useMutation(
			graphql(`
				mutation CreateSecret($opts: SecretCreateInput!) {
					secretCreate(opts: $opts) {
						id
					}
				}
			`),
			[invalidationKey],
		);

	const useMutationUpdateSecret = () =>
		useMutation(
			graphql(`
				mutation UpdateSecret($id: UUID!, $opts: SecretUpdateInput!) {
					secretUpdate(id: $id, opts: $opts) {
						id
					}
				}
			`),
			[invalidationKey],
		);

	const useMutationRemoveSecret = () =>
		useMutation(
			graphql(`
				mutation RemoveSecret($id: UUID!) {
					secretDelete(id: $id)
				}
			`),
			[invalidationKey],
		);

	return {
		secretsQuery,
		secrets,
		isLoading,
		useMutationCreateSecret,
		useMutationUpdateSecret,
		useMutationRemoveSecret,
	};
});
