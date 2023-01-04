import type { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';

import { DashboardWidgets } from '@/components/dashboard';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['dashboard', 'layout'])),
  },
});

const Home: NextPage = () => {
  return <DashboardWidgets />;
};

export default Home;
