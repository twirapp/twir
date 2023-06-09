import { Flex, Paper, Text } from '@mantine/core';
import { TwitchChat } from 'react-twitch-embed';

import { useTheme } from '@/services/dashboard/useTheme';

type Props = {
	channel: string;
};

export const TwitchChatWrapper = ({ channel }: Props) => {
	const { colorScheme } = useTheme();
	return (
		<Paper shadow="xs" p="lg" withBorder h={'100%'}>
			<Flex direction="column" gap="md" justify="stretch" align="flex-center" w={'100%'}>
				<Flex direction="row" justify="space-between" align="space-between">
					<Text>Chat</Text>
				</Flex>
				<TwitchChat
					channel={channel}
					darkMode={colorScheme === 'dark' ? true : false}
					width={'100%'}
				/>
			</Flex>
		</Paper>
	);
};
