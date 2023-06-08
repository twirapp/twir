import { Button, Flex, Grid, TextInput, Alert, Text, NativeSelect } from '@mantine/core';
import { showNotification } from '@mantine/notifications';
import { IconDeviceFloppy, IconInfoCircle } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import { useEffect, useState } from 'react';

import { IntegrationCard } from './card';

import { printError } from '@/services/api/error';
import { useRocketLeagueIntegration } from '@/services/api/integrations/rocketLeague';

type Platform = {
  value: string,
  code: string,
}

const platforms: Platform[] = [
	{
		value: 'Steam',
		code: 'steam',
	},
	{
		value: 'Epic Games',
		code: 'epic',
	},
	{
		value: 'Xbox',
		code: 'xbox',
	},
];
export const RocketLeagueIntegration: React.FC = () => {
	const manager = useRocketLeagueIntegration();
	const { data } = manager.useData();

  const { t: integrationsTranslate } = useTranslation('integrations');
  const { t } = useTranslation('common');
	const update = manager.usePost();

  const [username, setUsername] = useState<string>('');
  const [platform, setPlatform] = useState<Platform>(platforms[0]);
  const [platformSelectValue, setPlatformSelectValue] = useState<string>('Steam');

	useEffect(() => {
		if (typeof data?.username !== 'undefined' && typeof data?.code !== 'undefined') {
			setUsername(data.username);
			const plat = platforms.find(p => p.code === data.code);
			if (!plat) return;
			setPlatformSelectValue(plat.value);
			setPlatform(plat);
		}
	}, [data]);

  async function save() {
		if (!platform || !username) {
			printError('Platform and UserID must be filled');
			return; 
		}
		try{
			await update.mutateAsync({ username, code: platform.code });
			showNotification({
				title: 'Successful',
				message: (
					<div>
						Added new Rocket League integration into your account.
					</div>
				),
				color: 'green',
			});
		} catch(e) {
			console.log(e);
		}
  }

  function onChangePlatform(val: string) {
    const plat = platforms.find(plat => plat.value === val);
		if (!plat) {
			printError('Unable to find platform');
			return; 
		}

    setPlatform(plat);
    setPlatformSelectValue(val);
  }

  return (
		<IntegrationCard
			title="Rocket League"
			header={
				<Flex direction="row" gap="sm">
					<Button
						compact
						leftIcon={<IconDeviceFloppy />}
						variant="outline"
						color="green"
						onClick={save}
					>
						{t('save')}
					</Button>
				</Flex>
			}
		>
			<Grid align="flex-end">
				<Grid.Col span={12}>
					<TextInput
						label="User ID"
						value={username}
						onChange={(v) => setUsername(v.currentTarget.value)}
						placeholder={'Name'}
						mb="xs"
					/>
					<NativeSelect
						label="Choose the platform"
						placeholder="Pick one"
            value={platformSelectValue}
            onChange={(event) => onChangePlatform(event.currentTarget.value)}
						data={['Steam', 'Epic Games', 'Xbox']}
					/>
				</Grid.Col>
			</Grid>
			<Alert color={'lime'} icon={<IconInfoCircle />} mt={5}>
				<Text dangerouslySetInnerHTML={{ __html: integrationsTranslate('info.rocketleague') }} />
			</Alert>
		</IntegrationCard>
	);
};