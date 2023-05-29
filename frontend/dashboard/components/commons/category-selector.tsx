import { useTwitchGameCategories } from '@/services/api';
import { printError } from '@/services/api/error';
import {
	Autocomplete,
	Avatar,
	Group,
	Loader,
	useMantineTheme,
	Text,
	SelectItemProps,
	Paper,
} from '@mantine/core';
import { useDebouncedState } from '@mantine/hooks';
import React, { ReactNode, forwardRef, useRef, useState, Children } from 'react';
interface ItemProps extends SelectItemProps {
	image: string;
	id: string;
}

const Category = forwardRef<HTMLDivElement, ItemProps>(
	({ value, image, ...others }: ItemProps, ref) => (
		<div ref={ref} {...others}>
			<Group noWrap>
				<Avatar src={image} size="lg" />
				<div style={{ flex: 1 }}>
					<Text size="sm" weight={500}>
						{value}
					</Text>
				</div>
			</Group>
		</div>
	),
);

interface DropdownComponentProps {
	children: ReactNode[];
}

const dropdownComponent = ({ children }: DropdownComponentProps) => {
	const theme = useMantineTheme();

	return (
		<div
			style={{
				position: 'fixed',
				overflow: 'auto',
				backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.colors.gray[0],
				borderRadius: 10,
				maxHeight: 300,
				minWidth: 400,
				scrollbarWidth: 'thin',
				scrollbarColor:
					theme.colorScheme === 'dark'
						? `${theme.colors.dark[4]} ${theme.colors.dark[6]}`
						: `${theme.colors.gray[3]} ${theme.colors.gray[6]}`,
			}}
		>
			{children}
		</div>
	);
};

interface Props {
	label: string;
	setCategory: (value: CategoryType) => void;
	outerCategory: CategoryType;
	withAsterisk: boolean;
}

export interface CategoryType {
	name: string;
	id: string;
}

const CategorySelector = ({ label, setCategory, withAsterisk, outerCategory }: Props) => {
	const [category, setInnerCategory] = useState(outerCategory.name);

	const categories = useTwitchGameCategories(category);

	const handleChange = (val: string) => {
		setInnerCategory(val);
		const findedCategory = categories.data?.find((category) => category.name == val);
		setCategory({
			id: findedCategory?.id ?? '',
			name: findedCategory?.name ?? '',
		});
	};

	const data = categories?.data?.map((item) => ({
		image: item.box_art_url,
		value: item.name,
		id: item.id,
	}));
	return (
		<Autocomplete
			value={category}
			rightSection={categories.isLoading ? <Loader w={20} /> : <></>}
			label={label}
			itemComponent={Category}
			dropdownComponent={dropdownComponent}
			data={data ?? []}
			withAsterisk={withAsterisk}
			onChange={handleChange}
		/>
	);
};

export default CategorySelector;
