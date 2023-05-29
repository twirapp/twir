import { useCardStyles } from '@/styles/card';
import {
	ActionIcon,
	Button,
	Card,
	Flex,
	Group,
	Skeleton,
	Text,
	TextInput,
	Tooltip,
} from '@mantine/core';
import { useTranslation } from 'next-i18next';
import { categoriesAliasesManager, useGetStream, useModerationManager } from '@/services/api';
import { IconRefresh, IconSettings } from '@tabler/icons';
import CategorySelector, { CategoryType } from '../commons/category-selector';
import { useState } from 'react';
import AddCategoryAliasModal from './add-category-alias-modal';
import React from 'react';

const ManageCategoriesAliases = () => {
	const { t } = useTranslation('dashboard');
	const { classes } = useCardStyles();
	const [editDrawerOpened, setEditDrawerOpened] = useState(false);
	const [category, setCategory] = useState<CategoryType>({
		id: '',
		name: '',
	});
	const [title, setTitle] = useState('');

	const manager = useModerationManager();
	const titleUpdater = manager.useUpdateTitle();
	const categoryUpdater = manager.useUpdateCategory();

	const onSaveClick = () => {
		titleUpdater.mutateAsync(title);
		categoryUpdater.mutateAsync(category.id);
	};

	const streamManager = useGetStream();
	const { data: stream, refetch: refetchStream } = streamManager.useGet();

	React.useEffect(() => {
		const interval = setInterval(() => {
			setCategory({
				id: stream?.gameId ?? '',
				name: stream?.gameName ?? '',
			});
			setTitle(stream?.title ?? '');
			return () => clearInterval(interval);
		}, 60 * 1000);
	}, []);

	React.useEffect(() => {
		setCategory({
			id: stream?.gameId ?? '',
			name: stream?.gameName ?? '',
		});
		setTitle(stream?.title ?? '');
	}, [stream]);

	return (
		<Skeleton radius="md" visible={false}>
			<Card withBorder radius="md">
				<Card.Section withBorder inheritPadding py="sm">
					<Group position="apart">
						<Text weight={500}>{t('widgets.streamManager.title')}</Text>
						<Flex>
							<Tooltip label={'Refresh'} withArrow>
								<ActionIcon
									size={'lg'}
									variant={'default'}
									component="a"
									mr="xs"
									target={'_blank'}
									onClick={() => {
										refetchStream();
									}}
								>
									<IconRefresh width={20} />
								</ActionIcon>
							</Tooltip>
							<Tooltip label={'Settings'} withArrow>
								<ActionIcon
									size={'lg'}
									variant={'default'}
									component="a"
									target={'_blank'}
									onClick={() => {
										setEditDrawerOpened(true);
									}}
								>
									<IconSettings width={20} />
								</ActionIcon>
							</Tooltip>
						</Flex>
					</Group>
				</Card.Section>
				<Card.Section p="md" className={classes.card}>
					<CategorySelector
						label={t('widgets.streamManager.setCategory')}
						setCategory={setCategory}
						outerCategory={category}
						withAsterisk={false}
					/>
					<TextInput
						mt="md"
						label={t('widgets.streamManager.setTitle')}
						value={title}
						onChange={(event) => setTitle(event.currentTarget.value)}
					/>
					<Flex mt="md">
						<Button size="md" w="30%" color="green" onClick={onSaveClick}>
							{t('widgets.streamManager.save')}
						</Button>
					</Flex>
				</Card.Section>
				<AddCategoryAliasModal opened={editDrawerOpened} setOpened={setEditDrawerOpened} />
			</Card>
		</Skeleton>
	);
};

export default ManageCategoriesAliases;
