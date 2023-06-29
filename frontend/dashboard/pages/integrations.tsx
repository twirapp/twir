import { Center, Grid, Table, useMantineTheme } from '@mantine/core';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import dynamic from 'next/dynamic';

import { DonationAlertsIntegration } from '../components/integrations/donationalerts';
import { LastfmIntegration } from '../components/integrations/lastfm';
import { SpotifyIntegration } from '../components/integrations/spotify';
import { StreamlabsIntegration } from '../components/integrations/streamlabs';
import { VKIntegration } from '../components/integrations/vk';

import { DonatePayIntegration } from '@/components/integrations/donatepay';
import { FaceitIntegration } from '@/components/integrations/faceit';
import { ValorantIntegration } from '@/components/integrations/valorant';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
	props: {
		...(await serverSideTranslations(locale, ['integrations', 'layout', 'common'])),
	},
});

const DonateStreamIntegration = dynamic(() => import('../components/integrations/donateStream'), {
	ssr: false,
});

const DonatelloIntegration = dynamic(() => import('../components/integrations/donatello'), {
	ssr: false,
});

export default function Integrations() {
	const theme = useMantineTheme();

	return (
		<Center>
			<Table
				style={{
					backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[2],
				}}
				highlightOnHover
				w={'50%'}
				horizontalSpacing='xl'
			>
				<thead style={{ display: 'none' }}>
				<tr>
					<th style={{ width: 50 }}></th>
					<th></th>
					<th></th>
				</tr>
				</thead>
				<tbody>
					<SpotifyIntegration />
					<LastfmIntegration />
					<VKIntegration />
					<StreamlabsIntegration />
					<FaceitIntegration />
					<DonationAlertsIntegration />
					<DonatelloIntegration />
					<DonatePayIntegration />
					<DonateStreamIntegration />
					<ValorantIntegration />
				</tbody>
			</Table>
		</Center>
	);
}
