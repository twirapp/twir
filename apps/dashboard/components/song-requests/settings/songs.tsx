import { Card, Divider, NumberInput, Text } from '@mantine/core';
import React from 'react';

import { useYouTubeSettingsFormContext } from '@/components/song-requests/settings/form';

export const YouTubeSongsSettings: React.FC = () => {
  const form = useYouTubeSettingsFormContext();

  return (
    <Card style={{ minHeight: 500 }}>
      <Card.Section p={'xs'}><Text>Songs</Text></Card.Section>
      <Divider/>
      <Card.Section p={'md'}>
        <NumberInput label="Maximum number of songs in queue" {...form.getInputProps('maxRequests')} />
      </Card.Section>
    </Card>
  );
};