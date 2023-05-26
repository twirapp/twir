import { SimpleGrid } from '@mantine/core';
import React from 'react';

import { BotManage } from './bot-manage';
import ManageCategoriesAliases from './manage-categories-aliases';
// import { ServiceManage } from './service-manage';

export const DashboardWidgets = () => {
	return (
		<SimpleGrid
			spacing="lg"
			breakpoints={[
				{
					minWidth: 'md',
					cols: 2,
				},
				{
					minWidth: 'sm',
					cols: 1,
				},
			]}
		>
			<BotManage />
			<ManageCategoriesAliases />
			{/* <ServiceManage /> */}
		</SimpleGrid>
	);
};
