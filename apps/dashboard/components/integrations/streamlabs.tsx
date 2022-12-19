import { Group, Avatar, Text, Button, Flex, Alert } from '@mantine/core';
import { IconLogout, IconLogin } from '@tabler/icons';

import { IntegrationCard } from './card';

import { useStreamLabsIntegration } from '@/services/api/integrations';

export const StreamlabsIntegration: React.FC = () => {
  const manager = useStreamLabsIntegration();
  const { data } = manager.getIntegration();

  async function login() {
    const link = await manager.getAuthLink();
    if (link) {
      window.location.replace(link);
    }
  }

  return (
    <IntegrationCard
      title="Streamlabs"
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
      {!data && <Alert>Not logged in</Alert>}
      {data && (
        <Group position="apart" mt={10}>
          <Text weight={500} size={30}>
            {data.name}
          </Text>
          <Avatar src={data.avatar} h={150} w={150} style={{ borderRadius: 900 }} />
        </Group>
      )}
    </IntegrationCard>
  );
};
