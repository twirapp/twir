import { useCardStyles } from '@/styles/card';
import { Card, CardSection, Group, Skeleton, Text } from '@mantine/core';
import { useTranslation } from 'next-i18next';
import GameAliasesCreator from '../commons/game-aliases-creator';
import { categoriesAliasesManager } from '@/services/api';

const ManageCategoriesAliases = () => {
	const { t } = useTranslation('dashboard');
	const { classes } = useCardStyles();
	const { useUpdate } = categoriesAliasesManager();
	const updater;

	return (
		<Skeleton radius="md" visible={false}>
			<Card withBorder radius="md">
				<Card.Section withBorder inheritPadding py="sm">
					<Group position="apart">
						<Text weight={500}>{t('categoriesAliases.title')}</Text>
					</Group>
				</Card.Section>
				<Card.Section p="md" className={classes.card}>
					<GameAliasesCreator />
				</Card.Section>
			</Card>
		</Skeleton>
	);
};

export default ManageCategoriesAliases;
