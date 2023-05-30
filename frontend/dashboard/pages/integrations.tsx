import { Grid } from '@mantine/core';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import dynamic from 'next/dynamic';

import { DonationAlertsIntegration } from '../components/integrations/donationalerts';
import { LastfmIntegration } from '../components/integrations/lastfm';
import { SpotifyIntegration } from '../components/integrations/spotify';
import { StreamlabsIntegration } from '../components/integrations/streamlabs';
import { VKIntegration } from '../components/integrations/vk';

import { DonatelloIntegration } from '@/components/integrations/donatello';
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

const cols = {
	xs: 12,
	sm: 12,
	md: 5,
	lg: 5,
	xl: 5,
};

const DonateStreamIntegration = dynamic(() => import('../components/integrations/donateStream'), {
	ssr: false,
});

export default function Integrations() {
	return (
		<Grid justify="center">
			<Grid.Col {...cols}>
				<SpotifyIntegration />
			</Grid.Col>
			<Grid.Col {...cols}>
				<LastfmIntegration />
			</Grid.Col>
			<Grid.Col {...cols}>
				<VKIntegration />
			</Grid.Col>
			<Grid.Col {...cols}>
				<DonationAlertsIntegration />
			</Grid.Col>
			<Grid.Col {...cols}>
				<StreamlabsIntegration />
			</Grid.Col>
			<Grid.Col {...cols}>
				<FaceitIntegration />
			</Grid.Col>
			<Grid.Col {...cols}>
				<DonatePayIntegration />
			</Grid.Col>
			<Grid.Col {...cols}>
				<DonatelloIntegration />
			</Grid.Col>
			<Grid.Col {...cols}>
				<DonateStreamIntegration />
			</Grid.Col>
			<Grid.Col {...cols}>
				<ValorantIntegration />
			</Grid.Col>
		</Grid>
	);
}
