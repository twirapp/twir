import { Group, Avatar, Text, Flex, Button, Alert } from '@mantine/core';
import { IconBrandLastfm, IconLogin, IconLogout } from '@tabler/icons';
import { useTranslation } from 'next-i18next';

import { IntegrationCard } from './card';


import { useLastfm } from '@/services/api/integrations';

export const LastfmIntegration: React.FC = () => {
  const manager = useLastfm();
  const { data: profile } = manager.getIntegration;
  const { t } = useTranslation('integrations');

  async function login() {
    const { data } = manager.getAuthLink;
    if (data) {
      window.location.replace(data);
    }
  }

  return (
    <IntegrationCard
      title="Last.fm"
      icon={IconBrandLastfm}
      iconColor="red"
      header={
        <Flex direction="row" gap="sm">
          {profile && (
            <Button
              compact
              leftIcon={<IconLogout />}
              variant="outline"
              color="red"
              onClick={() => manager.logout.mutate()}
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
      {!profile && <Alert>{t('notLoggedIn')}</Alert>}
      {profile && (
        <Group position="apart" mt={10}>
          <Text weight={500} size={30}>
            {profile.name}
          </Text>
          <Avatar
            src={
              profile.image ||
              'https://lastfm.freetls.fastly.net/i/u/avatar170s/818148bf682d429dc215c1705eb27b98.png'
            }
            h={150}
            w={150}
            style={{ borderRadius: 900 }}
          />
        </Group>
      )}
    </IntegrationCard>
  );
};
