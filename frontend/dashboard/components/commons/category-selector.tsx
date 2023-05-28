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
} from '@mantine/core';
import { useDebouncedState } from '@mantine/hooks';
import React, { forwardRef, useRef, useState } from 'react';

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

interface Props {
	label: string;
	setCategory: (value: CategoryType) => void;
	withAsterisk: boolean;
}

export interface CategoryType {
	name: string;
	id: string;
}

const CategorySelector = ({ label, setCategory, withAsterisk }: Props) => {
	const [category, setInnerCategory] = useDebouncedState('', 200);

	const theme = useMantineTheme();
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
			rightSection={categories.isLoading ? <Loader w={20} /> : <></>}
			label={label}
			itemComponent={Category}
			data={data ?? []}
			withAsterisk={withAsterisk}
			onChange={handleChange}
		/>
	);
};

export default CategorySelector;
