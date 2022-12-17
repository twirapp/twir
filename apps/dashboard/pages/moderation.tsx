import { Grid } from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { Dashboard } from '@tsuwari/shared';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { useState } from 'react';
import useSWR from 'swr';

import { ModerationCard } from '../components/moderation/card';
import { ModerationDrawer } from '../components/moderation/drawer';
import { swrFetcher } from '../services/swrFetcher';

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableSettings, setEditableSettings] = useState<ChannelModerationSetting>({} as any);
  const [selectedDashboard] = useLocalStorage<Dashboard>({
    key: 'selectedDashboard',
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });

  const { data: settings } = useSWR<ChannelModerationSetting[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/moderation` : null,
    swrFetcher,
  );

  return (
    <div>
      <Grid justify="center">
        {settings &&
          settings.map((s, i) => (
            <Grid.Col key={i} xs={12} sm={12} md={5} lg={5} xl={5}>
              <ModerationCard
                settings={s as any}
                setEditableSettings={setEditableSettings}
                setEditDrawerOpened={setEditDrawerOpened}
              />
            </Grid.Col>
          ))}
      </Grid>

      <ModerationDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        settings={editableSettings}
      />
    </div>
  );
}
