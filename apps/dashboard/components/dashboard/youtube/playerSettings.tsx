import { ActionIcon, Button, Divider, Flex, Group, NumberInput, Select, SelectItem, Switch, Text, TextInput } from '@mantine/core';
import { useForm } from '@mantine/form';
import { openModal, closeAllModals } from '@mantine/modals';
import { showNotification } from '@mantine/notifications';
import { IconAdjustmentsHorizontal } from '@tabler/icons';
import { YoutubeSettings } from '@tsuwari/types/api';
import React, { useContext, useEffect, useState } from 'react';

import { noop } from '../../../util/chore';
import { RewardItem, RewardItemProps } from './reward';

import { useRewards } from '@/services/api';
import { useYoutubeModule } from '@/services/api/modules';

export const PlayerSettings: React.FC  = () => {
  return <ActionIcon
    onClick={() => {
      openModal({
          title: 'YouTube',
          children: <SettingsModal />,
      });
    }}>
    <IconAdjustmentsHorizontal />
  </ActionIcon>;
};


const SettingsModal: React.FC = () => {
  const form = useForm<YoutubeSettings>({
    initialValues: {
      enabled: true,
      acceptOnlyWhenOnline: true,
      channelPointsRewardId: '',
      maxRequests: 500,
      blackList: {
        artistsNames: [],
        songs: [],
        users: [],
        channels: [],
      },
      song: {
        maxLength: 10,
        minViews: 50000,
        acceptedCategories: [],
      },
      user: {
        maxRequests: 20,
        minWatchTime: 0,
        minFollowTime: 0,
        minMessages: 0,
      },
    },
    validate: {
      maxRequests: (v) => v > 500 && 'Max number of songs in queue is 500',
    },
  });

  const rewardsManager = useRewards();
  const { data: rewardsData } = rewardsManager();
  const [rewards, setRewards] = useState<RewardItemProps[]>([]);

  useEffect(() => {
    if (rewardsData) {
      const data = rewardsData
        .sort((a, b) => a.is_user_input_required === b.is_user_input_required ? 1 : -1)
        .map(r => ({
          value: r.id,
          label: r.title,
          description: r.is_user_input_required ? '' : 'Cannot be picked because have no no require input',
          image: r.image?.url_4x || r.default_image?.url_4x,
          disabled: !r.is_user_input_required,
        } as RewardItemProps));

      setRewards(data);
    }
  }, [rewardsData]);

  const youtube = useYoutubeModule();
  const { mutateAsync: updateSettings } = youtube.useUpdate();

  const { data: youtubeSettings } = youtube.useSettings();

  useEffect(() => {
    if (youtubeSettings) {
      form.setValues(youtubeSettings);
    }
  }, [youtubeSettings]);

  async function submit() {
    const validation = form.validate();
    console.log(validation.errors);
    if (validation.hasErrors) {
      for (const error of Object.values(validation.errors).flat(10) as string[]) {
        showNotification({
          title: 'Validation error',
          color: 'red',
          message: error,
        });
        console.log(error);
      }

      return;
    }

    updateSettings(form.values)
      .then(() => closeAllModals())
      .catch(noop);
  }

  return <form>
    <Group>
      <Switch
        label="Enabled"
        labelPosition='left'
        {...form.getInputProps('enabled', { type: 'checkbox' })}
      />
      <Switch
        label="Accept requests only when stream online"
        labelPosition='left'
        {...form.getInputProps('acceptOnlyWhenOnline', { type: 'checkbox' })}
      />
      <Select
        label="Channel points reward for requesting songs"
        placeholder="..."
        searchable
        itemComponent={RewardItem}
        allowDeselect
        data={rewards}
        {...form.getInputProps('channelPointsRewardId')}
      />

      <NumberInput label="Maximum number of songs in queue" {...form.getInputProps('maxRequests')} />
    </Group>

    <Divider style={{ marginTop: 10, marginBottom: 10 }} />

    <Group>
      <Text size={'lg'}>Restrictions</Text>
      <Flex direction={'row'} wrap={'wrap'} gap={'xs'}>
        <Button onClick={openSecondModal} size={'xs'} uppercase color="pink">Users</Button>
        <Button onClick={openSecondModal} size={'xs'} uppercase color="grape">Songs</Button>
        <Button onClick={openSecondModal} size={'xs'} uppercase color="violet">Channels</Button>
      </Flex>
    </Group>

    <Divider style={{ marginTop: 10, marginBottom: 10 }} />
    <Button color='green' onClick={submit}>Save</Button>
  </form>;
};

function openSecondModal() {
  return openModal({
    title: 'Please confirm your action',
    children: (
      <Text size="sm">
        This is second
      </Text>
    ),
  });
}