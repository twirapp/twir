import { useEffect, useState } from 'react';

import { useBuiltInVariables } from '@/services/api/builtInVariables.js';
import { useVariablesManager } from '@/services/api/crud.js';

type Variable = {
	name: string;
	example: string;
	description?: string;
	visible: boolean;
}

export const useAllVariables = () => {
	const [variables, setVariables] = useState<Variable[]>([]);

	const variablesManager = useVariablesManager();

	const customVariables = variablesManager.getAll({});
	const builtInVariables = useBuiltInVariables();

	useEffect(() => {
		setVariables((v) => [
			...v,
			...customVariables.data?.variables.map(variable => ({
				name: variable.name,
				description: variable.description,
				visible: true,
				example: `$(customvar|${variable.name})`,
			})) ?? [],
			...builtInVariables.data?.variables.map(variable => ({
				name: variable.name,
				description: variable.description,
				visible: variable.visible,
				example: variable.example,
			})) ?? [],
		]);
	}, [customVariables.data, builtInVariables.data]);

	return variables;
};
