import { categoriesAliasesManager } from '@/services/api';
import { useCardStyles } from '@/styles/card';
import {
	ActionIcon,
	Button,
	Card,
	Flex,
	Group,
	Modal,
	Skeleton,
	Text,
	Tooltip,
	useMantineTheme,
} from '@mantine/core';
import { IconSettings } from '@tabler/icons';
import React from 'react';
import { useTranslation } from 'react-i18next';
import CategorySelector from '../commons/category-selector';
import GameAliasesCreator from '../commons/game-aliases-creator';

type Props = {
	opened: boolean;
	setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

const AddCategoryAliasModal = (props: Props) => {
	const theme = useMantineTheme();
	const { t } = useTranslation('dashboard');
	const { classes } = useCardStyles();
	const { useCreateOrUpdate } = categoriesAliasesManager();
	const updater = useCreateOrUpdate();

	function onSubmit() {}

	return (
		<Modal
			opened={props.opened}
			onClose={() => props.setOpened(false)}
			title={
				<Button size="xs" color="green" onClick={onSubmit}>
					{'save'}
				</Button>
			}
			padding="xl"
			size="xl"
			overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
			overlayOpacity={0.55}
			overlayBlur={3}
		>
			<form>
				<Card withBorder radius="md">
					<Card.Section p="md" className={classes.card}>
						<CategorySelector />
						<GameAliasesCreator />
						<Flex mt="md">
							<Button size="md" w="30%" color="green">
								Добавить
							</Button>
						</Flex>
					</Card.Section>
				</Card>
			</form>
		</Modal>
	);
};

export default AddCategoryAliasModal;
