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
	Table,
	Text,
	TextInput,
	Tooltip,
	useMantineTheme,
} from '@mantine/core';
import { IconSearch, IconSettings, IconTrash } from '@tabler/icons';
import React, { Fragment, useState } from 'react';
import { useTranslation } from 'react-i18next';
import CategorySelector from '../commons/category-selector';
import GameAliasesCreator from '../commons/game-aliases-creator';
import { useDebouncedState, useViewportSize } from '@mantine/hooks';
import { useForm } from '@mantine/form';

type Props = {
	opened: boolean;
	setOpened: React.Dispatch<React.SetStateAction<boolean>>;
};

type ChannelCategoryAlias = {
	category: string;
	alias: string;
};

const AddCategoryAliasModal = (props: Props) => {
	const theme = useMantineTheme();
	const viewPort = useViewportSize();
	const { t } = useTranslation('dashboard');
	const { classes } = useCardStyles();
	const { useCreateOrUpdate, useGetAll, useDelete } = categoriesAliasesManager();
	const updater = useCreateOrUpdate();
	const aliases = useGetAll();
	const deleter = useDelete();

	const [category, setCategory] = useState('');
	const [alias, setAlias] = useState('');

	const form = useForm<ChannelCategoryAlias>({
		validate: {
			category: (value) => {
				if (!value.length || value.trim().length == 0) return 'Category cannot be empty';
				return null;
			},
			alias: (value) => {
				if (!value.length || value.trim().length == 0) return 'Alias cannot be empty';
				return null;
			},
		},
		initialValues: {
			category: '',
			alias: '',
		},
	});

	function onSubmit() {
		form.values.category = category;

		const validate = form.validate();
		if (validate.hasErrors) {
			console.log(validate.errors);
			return;
		}

		updater
			.mutateAsync({
				data: {
					...form.values,
				} as ChannelCategoryAlias,
			})
			.then(() => {})
			.catch((e) => console.log(e));
	}

	return (
		<Modal
			opened={props.opened}
			onClose={() => props.setOpened(false)}
			title={<Text size="xl">{t('widgets.streamManager.aliases.manage')}</Text>}
			padding="xl"
			size="xl"
			overlayColor={theme.colorScheme === 'dark' ? theme.colors.dark[9] : theme.colors.gray[2]}
			overlayOpacity={0.55}
			overlayBlur={3}
		>
			<Table mb={'md'} style={{ tableLayout: 'fixed', width: '100%' }}>
				<thead>
					<tr>
						<th style={{ width: '30%' }}>{t('widgets.streamManager.aliases.table.alias')}</th>
						<th style={{ width: '40%' }}>{t('widgets.streamManager.aliases.table.category')}</th>
						<th style={{ width: '10%' }}>{t('widgets.streamManager.aliases.table.actions')}</th>
					</tr>
				</thead>

				<tbody>
					{aliases?.data?.map((alias) => (
						<Fragment key={alias.id}>
							<tr style={{ padding: 5 }}>
								<td>{alias.alias}</td>
								<td>{alias.category}</td>
								<td>
									<Flex direction={'row'} gap="xs">
										<ActionIcon
											onClick={() => deleter.mutate(alias.id)}
											variant="filled"
											color="red"
										>
											<IconTrash size={14} />
										</ActionIcon>
									</Flex>
								</td>
							</tr>
						</Fragment>
					))}
				</tbody>
			</Table>
			<form>
				<Card withBorder radius="md">
					<Card.Section p="md" className={classes.card}>
						<CategorySelector
							label={t('widgets.streamManager.category')}
							setCategory={setCategory}
							{...form.getInputProps('category')}
						/>
						<TextInput
							label={t('name')}
							placeholder="your alias"
							withAsterisk
							{...form.getInputProps('alias')}
						/>
						<Flex mt="md">
							<Button size="md" w="30%" color="green" onClick={onSubmit}>
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
