import { Group, Avatar, Text, Button, Flex } from '@mantine/core';
import { IconBrandVk, IconLogin, IconLogout } from '@tabler/icons';

import { IntegrationCard } from './card';

import { useVkIntegration } from '@/services/api/integrations';

export const VKIntegration: React.FC = () => {
  const manager = useVkIntegration();
  const { data } = manager.getProfile();

  async function login() {
    const link = await manager.getAuthLink();
    if (link) {
      window.location.replace(link);
    }
  }

  return (
    <IntegrationCard
      title="VK"
      icon={IconBrandVk}
      iconColor="lightblue"
      header={
        <Flex direction="row" gap="sm">
          {data && (
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
      <Group position="apart" mt={10}>
        <Text weight={500} size={30}>
          Satont WorldWide
        </Text>
        <Avatar
          src={
            'https://images.unsplash.com/photo-1527004013197-933c4bb611b3?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=720&q=80'
          }
          h={150}
          w={150}
          style={{ borderRadius: 900 }}
        />
      </Group>
    </IntegrationCard>
  );
};
