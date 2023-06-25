import { Avatar, Badge, Button, Flex, Text, Tooltip } from '@mantine/core';
import { IconLogin, IconLogout } from '@tabler/icons';
import { useTranslation } from 'next-i18next';
import Image from 'next/image';
import { forwardRef } from 'react';

type Props = {
	integrationKey: string
	profile?: {
		avatar?: string,
		name?: string,
	}
	login: () => void | Promise<void>
	logout: () => void | Promise<void>
}

export const OAuthComponent: React.FC<Props> = (props) => {
	const { t } = useTranslation('integrations');

	return <tr key={props.integrationKey} style={{ height: 65 }}>
		<td><Image
			src={`/dashboard/assets/icons/brands/${props.integrationKey}.svg`}
			height={30}
			width={30}
			alt={props.integrationKey}
		/>
		</td>
		<td>
			{props.profile
				? <Flex direction={'row'} align='center'>
					<Avatar
						src={props.profile.avatar}
						h={30}
						w={30}
						style={{ borderRadius: 111 }}
					/>
					<Text size={30} ml={10}>
						{props.profile.name}
					</Text>
				</Flex>
				: <Badge>Not logged in</Badge>
			}
		</td>
		<td>
			<Flex direction='row' gap='sm' justify='flex-end'>
				{props.profile && (
					<Button
						compact
						leftIcon={<IconLogout />}
						variant='light'
						color='red'
						onClick={() => props.logout()}
					>
						{t('logout')}
					</Button>
				)}
				<Button compact leftIcon={<IconLogin />} variant='light' color='green' onClick={props.login}>
					{t('login')}
				</Button>
			</Flex>
		</td>
	</tr>;
};
