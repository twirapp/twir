import { ActionIcon, Flex, Grid, Input, Paper, Text, TextInput } from '@mantine/core';
import { IconRotate, IconSearch } from '@tabler/icons';
import { useState } from 'react';

export const UsersList = () => {
	const [searchValue, setSearchValue] = useState('');

	const onResetParticipantsClick = () => {};
	return (
		<Paper shadow="xs" p="lg" withBorder h={'100%'}>
			<Flex direction="column" gap="md" justify="stretch" align="flex-center" w={'100%'}>
				<Flex direction="row" justify="space-between" align="space-between">
					<Text>Users</Text>
					<Text>0 Users</Text>
				</Flex>
				<Flex direction="row" justify="space-between" align="center">
					<TextInput placeholder="Search users..." icon={<IconSearch />} />
					<ActionIcon
						size={36}
						ml="xs"
						color="blue"
						variant="filled"
						onClick={onResetParticipantsClick}
					>
						<IconRotate size={36} />
					</ActionIcon>
				</Flex>
			</Flex>
		</Paper>
	);
};
