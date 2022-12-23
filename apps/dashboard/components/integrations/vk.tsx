import { Group, Avatar, Text, Button, Flex, Alert } from '@mantine/core';
import { IconBrandVk, IconLogin, IconLogout } from '@tabler/icons';
import { useTranslation } from 'next-i18next';

import { IntegrationCard } from './card';


import { useVK } from '@/services/api/integrations';

export const VKIntegration: React.FC = () => {
  const manager = useVK();
  const logout = manager.useLogout();
  const { data: profile } = manager.useData();
  const auth = manager.useGetAuthLink();
  const { t } = useTranslation('integrations');

  async function login() {
    if (auth.data) {
      window.location.replace(auth.data);
    }
  }

  return (
    <IntegrationCard
      title="VK"
      icon={IconBrandVk}
      iconColor="lightblue"
      header={
        <Flex direction="row" gap="sm">
          {profile && (
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
      {!profile && <Alert>{t('notLoggedIn')}</Alert>}
      {profile && (
        <Group position="apart" mt={10}>
          <Text weight={500} size={30}>
            {profile.first_name} {profile.last_name}
          </Text>
          <Avatar src={profile.photo_max_orig} h={150} w={150} style={{ borderRadius: 900 }} />
        </Group>
      )}
    </IntegrationCard>
  );
};
