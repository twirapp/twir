import { Group, Avatar, Text, Flex, Button, Alert } from '@mantine/core';
import { IconLogout, IconLogin, IconInfoCircle } from '@tabler/icons';
import { useTranslation } from 'next-i18next';

import { IntegrationCard } from './card';

import { useDonationAlerts } from '@/services/api/integrations';

export const DonationAlertsIntegration: React.FC = () => {
  const manager = useDonationAlerts();
  const logout = manager.useLogout();
  const { data } = manager.useData();
  const auth = manager.useGetAuthLink();
  const { t } = useTranslation('integrations');

  async function login() {
    if (auth.data) {
      window.location.replace(auth.data);
    }
  }

  return (
    <IntegrationCard
      title="DonationAlerts"
      header={
        <Flex direction="row" gap="sm">
          {data && (
            <Button
              compact
              leftIcon={<IconLogout />}
              variant="outline"
              color="red"
              onClick={() => logout.mutate()}
            >
              {t('logout')}
            </Button>
          )}
          <Button compact leftIcon={<IconLogin />} variant="outline" color="green" onClick={login}>
            {t('login')}
          </Button>
        </Flex>
      }
    >
      {!data && <Alert>{t('notLoggedIn')}</Alert>}
      {data && (
        <Group position="apart" mt={10}>
          <Text weight={500} size={30}>
            {data.name}
          </Text>
          <Avatar src={data.avatar} h={150} w={150} style={{ borderRadius: 900 }} />
        </Group>
      )}

      <Alert color={'lime'} icon={<IconInfoCircle />} mt={5}>
        <Text dangerouslySetInnerHTML={{ __html: t('info.donations') }} />
      </Alert>
    </IntegrationCard>
  );
};
