import { Group, Avatar, Text } from '@mantine/core';

import { IntegrationCard } from './card';

export const StreamlabsIntegration: React.FC = () => {
  return (
    <IntegrationCard title="Streamlabs">
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
