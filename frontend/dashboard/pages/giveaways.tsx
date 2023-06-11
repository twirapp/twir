import { Grid } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';

import { Settings } from '@/components/giveaways/Settings';
import { TwitchChatWrapper } from '@/components/giveaways/TwitchChatWrapper';
import { UsersList } from '@/components/giveaways/UsersList';
import { useTheme } from '@/services/dashboard/useTheme';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
	props: {
		...(await serverSideTranslations(locale, ['giveaways', 'layout'])),
	},
});

export default function () {
	const { theme, colorScheme } = useTheme();
	const { width, height } = useViewportSize();
	// const manager = useGiveawaysSettings();
	// const { data: giveaways } = manager.useGet();

	return (
		<div>
			<Grid>
				<Grid.Col span="auto">
					<UsersList />
				</Grid.Col>
				<Grid.Col span={6}>
					<Settings />
				</Grid.Col>
				<Grid.Col span="auto">
					<TwitchChatWrapper channel="danluki" />
				</Grid.Col>
			</Grid>
		</div>
	);
}
