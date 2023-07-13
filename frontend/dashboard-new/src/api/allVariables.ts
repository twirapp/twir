import { useQuery } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp.js';

type Variable = {
	name: string;
	example: string;
	description?: string;
	visible: boolean;
	isBuiltIn: boolean;
}

export const useAllVariables = () => {
	return useQuery({
		queryKey: ['allVariables'],
		queryFn: async () => {
			const [builtIn, custom] = await Promise.all([
				protectedApiClient.builtInVariablesGetAll({}),
				protectedApiClient.variablesGetAll({}),
			]);

			const variables: Variable[] = [];

			for (const variable of custom.response.variables) {
				variables.push({
					name: variable.name,
					description: variable.description,
					visible: true,
					example: `$(customvar|${variable.name})`,
					isBuiltIn: false,
				});
			}

			for (const variable of builtIn.response.variables) {
				variables.push({
					name: variable.name,
					description: variable.description,
					visible: variable.visible,
					example: variable.example,
					isBuiltIn: true,
				});
			}

			return variables;
		},
	});
};
