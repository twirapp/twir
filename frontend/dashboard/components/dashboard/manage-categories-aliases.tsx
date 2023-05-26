import { useCardStyles } from '@/styles/card';
import {
	ActionIcon,
	Button,
	Card,
	CardSection,
	Flex,
	Group,
	Skeleton,
	Text,
	Tooltip,
} from '@mantine/core';
import { useTranslation } from 'next-i18next';
import GameAliasesCreator from '../commons/game-aliases-creator';
import { categoriesAliasesManager } from '@/services/api';
import { IconSettings } from '@tabler/icons';
import CategorySelector from '../commons/category-selector';
import { useState } from 'react';
import AddCategoryAliasModal from './add-category-alias-modal';

const ManageCategoriesAliases = () => {
	const { t } = useTranslation('dashboard');
	const { classes } = useCardStyles();
	const { useCreateOrUpdate } = categoriesAliasesManager();
	const updater = useCreateOrUpdate();
	const [editDrawerOpened, setEditDrawerOpened] = useState(false);

	return (
		<Skeleton radius="md" visible={false}>
			<Card withBorder radius="md">
				<Card.Section withBorder inheritPadding py="sm">
					<Group position="apart">
						<Text weight={500}>{t('categoriesAliases.title')}</Text>
						<Tooltip label={'Settings'} withArrow>
							<ActionIcon size={'lg'} variant={'default'} component="a" target={'_blank'}>
								<IconSettings width={20} />
							</ActionIcon>
						</Tooltip>
					</Group>
				</Card.Section>
				<Card.Section p="md" className={classes.card}>
					<CategorySelector />
					<GameAliasesCreator />
					<Flex mt="md">
						<Button
							size="md"
							w="30%"
							color="green"
							onClick={() => {
								setEditDrawerOpened(true);
							}}
						>
							Добавить
						</Button>
					</Flex>
				</Card.Section>
				<AddCategoryAliasModal opened={editDrawerOpened} setOpened={setEditDrawerOpened} />
			</Card>
		</Skeleton>
	);
};

export default ManageCategoriesAliases;
