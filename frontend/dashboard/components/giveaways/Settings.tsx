import {
	Button,
	Text,
	Paper,
	Flex,
	Radio,
	Group,
	TextInput,
	Switch,
	NumberInput,
	MultiSelect,
	RangeSlider,
	Slider,
} from '@mantine/core';
import { IconCheck, IconShieldHalfFilled } from '@tabler/icons';
import { useState } from 'react';

import { DurationPicker } from './ui/DurationPicker';

import { useRolesApi } from '@/services/api';

const SLIDER_MARKS = [
	{ value: 0, label: 0 },
	{ value: 1, label: 1 },
	{ value: 2, label: 2 },
	{ value: 3, label: 3 },
	{ value: 4, label: 4 },
	{ value: 5, label: 5 },
	{ value: 6, label: 6 },
	{ value: 7, label: 7 },
	{ value: 8, label: 8 },
	{ value: 9, label: 9 },
	{ value: 10, label: 10 },
];

export const Settings = () => {
	const [running, setRunning] = useState(false);
	const [giveawayType, setGiveawayType] = useState('keyword');
	const [keyword, setKeyword] = useState('');
	const [isNeedAnnounce, setIsNeedAnnounce] = useState(false);
	const [minimumWatchTime, setMinimumWatchTime] = useState<number | null>(0);
	const [minimumFollowTime, setMinimumFollowTime] = useState<number | null>(0);
	const [minimumMessages, setMinimumMessages] = useState<number | undefined>(0);
	const [minimumSubTier, setMinimumSubTier] = useState<number | undefined>(1);
	const [minimumSubTime, setMinimumSubTime] = useState<number | null>(0);
	const [rolesIds, setRolesIds] = useState<string[]>([]);
	const [giveawayMinimumNumber, setGiveawayMinimumNumber] = useState<number | undefined>(0);
	const [giveawayMaximumNumber, setGiveawayMaximumNumber] = useState<number | undefined>(0);
	const [winnersCount, setWinnersCount] = useState<number | undefined>(0);
	const [subLuck, setSubLuck] = useState(0);
	const [subTier1Luck, setSubTier1Luck] = useState(0);
	const [subTier2Luck, setSubTier2Luck] = useState(0);
	const [subTier3Luck, setSubTier3Luck] = useState(0);

	const rolesManager = useRolesApi();
	const { data: roles } = rolesManager.useGetAll();

	const onRunningClick = () => {
		setRunning(!running);
	};

	const onRollItClick = () => {};

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
				{giveawayType === 'keyword' ? (
					<TextInput
						label="Keyword"
						value={keyword}
						disabled={giveawayType !== 'keyword'}
						onChange={(event) => setKeyword(event.currentTarget.value)}
					/>
				) : (
					<>
						<NumberInput
							label="Minimum number"
							value={giveawayMinimumNumber}
							onChange={setGiveawayMinimumNumber}
						/>
						<NumberInput
							min={1}
							max={3}
							label="Maximum number"
							value={giveawayMaximumNumber}
							onChange={setGiveawayMaximumNumber}
						/>
					</>
				)}
				<Switch
					checked={isNeedAnnounce}
					onChange={(event) => setIsNeedAnnounce(event.currentTarget.checked)}
					label="Announce winner in chat"
				/>
				<DurationPicker
					placeholder="24h"
					label="Minimum watch time"
					onChange={setMinimumWatchTime}
				/>
				<DurationPicker
					placeholder="24h"
					label="Minimum follow time"
					onChange={setMinimumFollowTime}
				/>
				<NumberInput
					label="Minimum messages"
					value={minimumMessages}
					onChange={setMinimumMessages}
				/>
				<NumberInput
					min={1}
					max={3}
					label="Minimum sub tier"
					value={minimumSubTier}
					onChange={setMinimumSubTier}
				/>
				<DurationPicker placeholder="24h" label="Minimum sub time" onChange={setMinimumSubTime} />
				<MultiSelect
					data={
						roles?.map((r) => ({
							value: r.id,
							label: r.name,
							group: r.type !== 'CUSTOM' ? 'System' : 'Custom',
						})) ?? []
					}
					icon={<IconShieldHalfFilled size={18} />}
					label={'Roles'}
					placeholder="That roles will access to this giveaway."
					description={'Leave blank for everyone.'}
					clearButtonLabel="Clear selection"
					clearable
					value={rolesIds}
					onChange={setRolesIds}
				/>
				<NumberInput
					min={1}
					label="Winners count"
					value={winnersCount}
					onChange={setWinnersCount}
				/>
				<Text>Sub luck</Text>
				<Slider
					mb="xs"
					value={subLuck}
					onChange={setSubLuck}
					defaultValue={0}
					min={0}
					max={10}
					marks={SLIDER_MARKS}
				/>
				<Text>Tier 1 sub luck</Text>
				<Slider
					mb="xs"
					value={subTier1Luck}
					onChange={setSubTier1Luck}
					defaultValue={0}
					min={0}
					max={10}
					marks={SLIDER_MARKS}
				/>
				<Text>Tier 2 sub luck</Text>
				<Slider
					mb="xs"
					value={subTier2Luck}
					onChange={setSubTier2Luck}
					defaultValue={0}
					min={0}
					max={10}
					marks={SLIDER_MARKS}
				/>
				<Text>Tier 3 sub luck</Text>
				<Slider
					mb="xs"
					value={subTier3Luck}
					onChange={setSubTier3Luck}
					defaultValue={0}
					min={0}
					max={10}
					marks={SLIDER_MARKS}
				/>
				<Button onClick={onRollItClick}>Roll it!</Button>
			</Flex>
		</Paper>
	);
};
