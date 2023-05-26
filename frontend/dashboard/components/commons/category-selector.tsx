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
import React, { Component, forwardRef, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

interface Props {
	channelId: string;
}

interface ItemProps extends SelectItemProps {
	description: string;
	image: string;
}

const Category = forwardRef<HTMLDivElement, ItemProps>(
	({ description, value, image, ...others }: ItemProps, ref) => (
		<Group noWrap>
			<Avatar src={image} size="lg" />
			<div style={{ flex: 1 }}>
				<Text size="sm" weight={500}>
					{description}
				</Text>
			</div>
		</Group>
	),
);

const CategorySelector = ({ channelId }: Props) => {
	const timeoutRef = useRef<number>(-1);
	const [category, setCategory] = useState('');
	const [loading, setLoading] = useState(false);

	const { t } = useTranslation('twitch');
	const theme = useMantineTheme();
	const categories = useTwitchGameCategories(category, channelId);

	const handleChange = (val: string) => {
		window.clearTimeout(timeoutRef.current);
		setCategory(val);

		if (val.trim().length === 0) {
			setLoading(false);
		} else {
			setLoading(true);
			timeoutRef.current = window.setTimeout(() => {
				setLoading(false);
				categories.refetch();
			}, 1000);
		}
	};

	const data = categories?.data?.map((item) => ({ image: item.box_art_url, value: item.name }));

	return (
		<Autocomplete
			label="Select a category"
			placeholder="Pick one"
			itemComponent={Category}
			data={data ?? []}
			onChange={onSearchInputChange}
		/>
	);
};

export default CategorySelector;
