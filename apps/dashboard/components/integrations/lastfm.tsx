import { Group, Avatar, Text, Flex, Button, Alert } from '@mantine/core';
import { IconBrandLastfm, IconLogin, IconLogout } from '@tabler/icons';

import { IntegrationCard } from './card';

import { useLastfmIntegration } from '@/services/api/integrations/lastfm';

export const LastfmIntegration: React.FC = () => {
  const manager = useLastfmIntegration();
  const { data: profile } = manager.getProfile();

  async function login() {
    const link = await manager.getAuthLink();
    if (link) {
      window.location.replace(link);
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
              onClick={manager.logout}
            >
              Logout
            </Button>
          )}
          <Button compact leftIcon={<IconLogin />} variant="outline" color="green" onClick={login}>
            Login
          </Button>
        </Flex>
      }
    >
      {!profile && <Alert>Not logged in</Alert>}
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
