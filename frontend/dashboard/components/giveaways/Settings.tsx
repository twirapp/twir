import { Button, Text, Paper, Flex, Radio, Group, Grid, TextInput, Switch } from '@mantine/core';
import { IconCheck, IconCross } from '@tabler/icons';
import { useState } from 'react';

export const Settings = () => {
	const [running, setRunning] = useState(false);
	const [giveawayType, setGiveawayType] = useState('keyword');
	const [keyword, setKeyword] = useState('');
	const [isNeedAnnounce, setIsNeedAnnounce] = useState(false);
	const onRunningClick = () => {
		setRunning(!running);
	};

	return (
		<Paper shadow="xs" p="lg" withBorder h={'100%'}>
			<Flex direction="column" gap="md" justify="stretch" align="flex-center" w={'100%'}>
				<Flex direction="row" justify="space-between" align="space-between">
					<Text>Settings</Text>
					{running ? (
						<Button onClick={onRunningClick} color="green" leftIcon={<IconCheck />}>
							Running
						</Button>
					) : (
						<Button onClick={onRunningClick} color="red" leftIcon={<IconCheck />}>
							Stopped
						</Button>
					)}
				</Flex>
				<Radio.Group
					name="keywordType"
					label="Giveaway type"
					value={giveawayType}
					onChange={setGiveawayType}
					withAsterisk
				>
					<Group>
						<Radio label="Random number" value="number" />
						<Radio label="Keyword" value="keyword" />
					</Group>
				</Radio.Group>
				<TextInput
					value={keyword}
					disabled={giveawayType !== 'keyword'}
					onChange={(event) => setKeyword(event.currentTarget.value)}
				/>
				<Switch
					checked={isNeedAnnounce}
					onChange={(event) => setIsNeedAnnounce(event.currentTarget.checked)}
					label="Announce winner in chat"
				/>
				<Button>Roll it!</Button>
			</Flex>
		</Paper>
	);
};
