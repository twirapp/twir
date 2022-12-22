import { Grid } from '@mantine/core';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { ModerationCard } from '../components/moderation/card';
import { ModerationDrawer } from '../components/moderation/drawer';

import { useModerationSettings } from '@/services/api';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
    props: {
        ...(await serverSideTranslations(locale, ['moderation', 'layout'])),
    },
});

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableSettings, setEditableSettings] = useState<ChannelModerationSetting>({} as any);

  const { data: settings } = useModerationSettings().getAll;

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
