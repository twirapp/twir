import type { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import dynamic from 'next/dynamic';

const WidgetsDynamic = dynamic(
  async () => {
    const mod = await import('@/components/dashboard');
    return mod.DashboardWidgets;
  },
  { ssr: false },
);

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['dashboard', 'layout'])),
  },
});

const Home: NextPage = () => {
  return <WidgetsDynamic />;
};

export default Home;
