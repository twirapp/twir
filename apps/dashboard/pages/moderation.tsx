import { Grid } from '@mantine/core';
import {
  IconLambda,
  IconLetterCaseUpper,
  IconLink,
  IconMoodSmile,
  IconPlaylistX,
  IconTextWrapDisabled,
} from '@tabler/icons';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { useState } from 'react';

import { ModerationCard } from '../components/moderation/card';
import { ModerationDrawer } from '../components/moderation/drawer';

const settings = [
  {
    id: '21aaa9ac-bcd9-4d8a-ac99-c9588ab490be',
    type: 'longMessage',
    channelId: '128644134',
    enabled: false,
    subscribers: false,
    vips: false,
    banTime: 600,
    banMessage: null,
    warningMessage: null,
    checkClips: false,
    triggerLength: 300,
    maxPercentage: 50,
    blackListSentences: [],
  },
  {
    id: 'e5e0382b-16fa-4163-a931-e84bd95b2296',
    type: 'emotes',
    channelId: '128644134',
    enabled: false,
    subscribers: false,
    vips: false,
    banTime: 600,
    banMessage: null,
    warningMessage: null,
    checkClips: false,
    triggerLength: 300,
    maxPercentage: 50,
    blackListSentences: [],
  },
  {
    id: 'd55adf93-ac9e-405a-864b-8b446f8bfe74',
    type: 'blacklists',
    channelId: '128644134',
    enabled: false,
    subscribers: false,
    vips: false,
    banTime: 600,
    banMessage: null,
    warningMessage: null,
    checkClips: false,
    triggerLength: 300,
    maxPercentage: 50,
    blackListSentences: [],
  },
  {
    id: '6b84df52-a69f-4f53-b558-7a10791b375b',
    type: 'symbols',
    channelId: '128644134',
    enabled: false,
    subscribers: false,
    vips: false,
    banTime: 600,
    banMessage: null,
    warningMessage: null,
    checkClips: false,
    triggerLength: 300,
    maxPercentage: 50,
    blackListSentences: [],
  },
  {
    id: 'fe1cdc42-a7dc-4422-83f4-c16fd5fd5351',
    type: 'links',
    channelId: '128644134',
    enabled: false,
    subscribers: false,
    vips: false,
    banTime: 2,
    banMessage: '123123',
    warningMessage: null,
    checkClips: true,
    triggerLength: 300,
    maxPercentage: 50,
    blackListSentences: [],
  },
  {
    id: '8001b95b-d5fd-4736-a3cf-48d5084ac141',
    type: 'caps',
    channelId: '128644134',
    enabled: false,
    subscribers: false,
    vips: false,
    banTime: 600,
    banMessage: null,
    warningMessage: null,
    checkClips: false,
    triggerLength: 300,
    maxPercentage: 50,
    blackListSentences: [],
  },
];

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableSettings, setEditableSettings] = useState<ChannelModerationSetting>({} as any);

  return (
    <div>
      <Grid justify="center">
        {settings.map((s, i) => (
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
