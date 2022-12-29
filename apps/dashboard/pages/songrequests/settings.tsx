import {
  ActionIcon,
  Button,
  Divider,
  Flex,
  Grid,
  Group,
  NumberInput,
  Popover,
  ScrollArea,
  Select,
  Switch,
  Text,
  TextInput,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { closeAllModals } from '@mantine/modals';
import { showNotification } from '@mantine/notifications';
import { IconCheck, IconPlus, IconX } from '@tabler/icons';
import { YouTubeSettings } from '@tsuwari/types/api';
import { GetServerSideProps, NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import React, { useEffect, useState } from 'react';

import { noop } from '../../util/chore';

import { RewardItem, RewardItemProps } from '@/components/dashboard/youtube/reward';
import { useRewards } from '@/services/api';
import { useYoutubeModule } from '@/services/api/modules';

const cols = {
  xs: 12,
  sm: 12,
  md: 6,
  lg: 4,
  xl: 4,
};

export const getServerSideProps: GetServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale!, ['layout'])),
  },
});

const Settings: NextPage = () => {
  const form = useForm<YouTubeSettings>({
    initialValues: {
      enabled: true,
      acceptOnlyWhenOnline: true,
      channelPointsRewardId: '',
      maxRequests: 500,
      denyList: {
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

  const youtube = useYoutubeModule();
  const { mutateAsync: updateSettings } = youtube.useUpdate();

  const { data: youtubeSettings } = youtube.useSettings();

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
          description: r.is_user_input_required ? '' : 'Cannot be picked because have no require input',
          image: r.image?.url_4x || r.default_image?.url_4x,
          disabled: !r.is_user_input_required,
        } as RewardItemProps));

      setRewards(data);
    }
  }, [rewardsData]);

  async function submit() {
    const validation = form.validate();
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
      .then(() => {
        closeAllModals();
        showNotification({
          message: 'Settings updated',
          color: 'green',
        });
      })
      .catch(noop);
  }

  const [newDenyUser, setNewDenyUser] = useState('');
  const [filterUsers, setFilterUsers] = useState('');
  const [newDenyUserPopover, setNewDenyUserPopover] = useState(false);

  function insertDenyUser() {
    if (!newDenyUser) return;
    form.insertListItem('denyList.users', { userName: newDenyUser, id: '' });
    setNewDenyUser('');
    setNewDenyUserPopover(false);
  }

  return (
    <form>
      <Flex justify={'space-between'}>
        <Text size={'lg'}>Songrequests settings</Text>
        <Button color={'green'} onClick={submit}>Save</Button>
      </Flex>
      <Grid
        align="flex-start"
        justify="center"
        style={{}}
        grow
        gutter="xl"
      >
        <Grid.Col {...cols}>
          <Text size={'xl'}>General</Text>
          <Switch
            label="Enabled"
            labelPosition="left"
            {...form.getInputProps('enabled', { type: 'checkbox' })}
          />
          <Switch
            label="Accept requests only when stream online"
            labelPosition="left"
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
        </Grid.Col>
        <Grid.Col {...cols}>
          <Text size={'xl'}>Users</Text>

          <NumberInput label="Maximum songs by user in queue" {...form.getInputProps('user.maxRequests')} />
          <NumberInput
            label="Minimal watch time of user for request song (minutes)" {...form.getInputProps('user.minWatchTime')} />
          <NumberInput label="Minimal messages by user for request song" {...form.getInputProps('user.minMessages')} />
          <NumberInput
            label="Minimal follow time for request song (minutes)" {...form.getInputProps('user.minFollowTime')} />

          <Divider style={{ marginTop: 10 }}/>

          <Flex direction="row" justify="space-between" style={{ marginTop: 10 }}>
            <Text size="sm">Denied users for request</Text>
            <Popover width={200} position="bottom" withArrow shadow="md" opened={newDenyUserPopover}>
              <Popover.Target>
                <ActionIcon
                  color="green"
                  variant={'filled'}
                  size={'sm'}
                  onClick={() => setNewDenyUserPopover(!newDenyUserPopover)}
                ><IconPlus/></ActionIcon>
              </Popover.Target>
              <Popover.Dropdown>
                <Flex direction={'row'} gap={'sm'}>
                  <TextInput placeholder="enter username" onChange={(v) => setNewDenyUser(v.currentTarget.value)}/>
                  <ActionIcon onClick={() => insertDenyUser()}><IconCheck/></ActionIcon>
                </Flex>
              </Popover.Dropdown>
            </Popover>
          </Flex>

          {form.values.denyList.users.length
            ? <TextInput style={{ marginTop: 10 }} placeholder="filter..."
                         onChange={(v) => setFilterUsers(v.target.value)}
            />
            : ''
          }

          <ScrollArea type={'always'}>
            <div style={{ maxHeight: 300 }}>
              {form.values.denyList.users.length ?
                form.values.denyList.users
                  .filter(s => s.userName.toLowerCase().includes(filterUsers.toLowerCase()))
                  .map((s, i) => <Group style={{ maxHeight: 280 }}>
                    <Flex
                      key={s.userId}
                      direction="row"
                      justify="space-between"
                      style={{ width: '95%', marginTop: 10 }}
                      gap="sm"
                    >
                      <Text size={'sm'} lineClamp={4}>{s.userName}</Text>
                      <ActionIcon onClick={() => form.removeListItem('denyList.users', i)}>
                        <IconX/>
                      </ActionIcon>
                    </Flex>
                  </Group>)
                : ''}
            </div>
          </ScrollArea>
        </Grid.Col>
        <Grid.Col {...cols}>
          <Text size={'xl'}>Songs</Text>

          <NumberInput label="Max length of song for request (minutes)" {...form.getInputProps('song.maxLength')} />
          <NumberInput label="Minimal views on song for request" {...form.getInputProps('song.minViews')} />

          <Divider style={{ marginTop: 10 }}/>
        </Grid.Col>
      </Grid>
    </form>
  );
};

export default Settings;