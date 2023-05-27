import { useTwitchGameCategories } from '@/services/api';
import {
	Autocomplete,
	Avatar,
	Group,
	Loader,
	useMantineTheme,
	Text,
	SelectItemProps,
} from '@mantine/core';
import React, { forwardRef, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

interface ItemProps extends SelectItemProps {
	image: string;
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
	setCategory: (value: string) => void;
}

const CategorySelector = ({ label, setCategory }: Props) => {
	const timeoutRef = useRef<number>(-1);
	const [category, setInnerCategory] = useState('');

	const theme = useMantineTheme();
	const categories = useTwitchGameCategories(category);

	const handleChange = (val: string) => {
		// window.clearTimeout(timeoutRef.current);
		setInnerCategory(val);
		setCategory(val);
		// if (val.trim().length === 0) {
		// 	setLoading(false);
		// } else {
		// 	setLoading(true);
		// 	timeoutRef.current = window.setTimeout(() => {
		// 		setLoading(false);
		// 	}, 1000);
		// }
	};
	const data = categories?.data?.map((item) => ({ image: item.box_art_url, value: item.name }));
	return (
		<Autocomplete
			rightSection={categories.isLoading ? <Loader w={20} /> : <></>}
			label={label}
			itemComponent={Category}
			className="scr"
			data={data ?? []}
			onChange={handleChange}
		/>
	);
};

export default CategorySelector;
