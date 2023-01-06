import { Card, Divider, Grid, Text, Textarea } from '@mantine/core';
import { useTranslation } from 'next-i18next';
import React from 'react';

import { useYouTubeSettingsFormContext } from '@/components/song-requests/settings/form';

export const YouTubeTranslationSettings: React.FC = () => {
  const form = useYouTubeSettingsFormContext();
  const [t] = useTranslation('song-requests-settings');

  const createKeys = (obj: Record<string, string | Record<string, string>>) => {
    return Object.keys(form.values.translations).reduce((acc, curr) => {
      if (typeof obj[curr] === 'string') return [...acc, curr];
      return [
        ...acc,
        ...Object.keys(obj[curr]).map(k => `${curr}.${k}`),
      ];
    }, [] as string[]);
  };

  return (
    <Card style={{ minHeight: 500 }}>
      <Card.Section p={'xs'}>
        <Text>{t('translations.title')}</Text>
      </Card.Section>
      <Divider />
      <Card.Section p={'md'}>
        <Grid grow>
          {createKeys(form.values.translations as any)
              .flat()
              .map(k =>
                <Grid.Col span={4}>
                  <Textarea
                    autosize
                    minRows={2}
                    {...form.getInputProps(`translations.${k}`)}
                  />
                </Grid.Col>)
          }
        </Grid>
      </Card.Section>
    </Card>
  );
};
