import { useQuery } from '@tanstack/vue-query';
import { VariableType } from '@twir/api/messages/variables/variables';

import { protectedApiClient } from '@/api/twirp.js';

type Variable = {
	name: string;
	example: string;
	description?: string;
	visible: boolean;
	isBuiltIn: boolean;
	canBeUsedInRegistry: boolean,
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
					example: `customvar|${variable.name}`,
					isBuiltIn: false,
					canBeUsedInRegistry: variable.type !== VariableType.SCRIPT,
				});
			}

			for (const variable of builtIn.response.variables) {
				variables.push({
					name: variable.name,
					description: variable.description,
					visible: variable.visible,
					example: variable.example,
					isBuiltIn: true,
					canBeUsedInRegistry: variable.canBeUsedInRegistry,
				});
			}

			return variables;
		},
	});
};
