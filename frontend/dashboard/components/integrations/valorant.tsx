import {
	TextInput,
} from '@mantine/core';
import { useEffect, useState } from 'react';

import { ManualComponent } from './manualComponent';

import { useValorantIntegration } from '@/services/api/integrations/valorant';

export const ValorantIntegration: React.FC = () => {
	const manager = useValorantIntegration();
	const { data } = manager.useData();
	const update = manager.usePost();

	const [username, setUsername] = useState<string>();

	useEffect(() => {
		if (typeof data?.username !== 'undefined') {
			setUsername(data.username);
		}
	}, [data]);

	async function save() {
		if (typeof username == 'undefined') return;
		await update.mutateAsync({ username });
	}

	return (
		<ManualComponent
			integrationKey={'valorant'}
			save={save}
			imageSize={50}
			body={
				<TextInput
					label='Valorant username'
					value={username}
					onChange={(v) => setUsername(v.currentTarget.value)}
					placeholder={'Name#tag'}
				/>
			}
		/>
	);
};
