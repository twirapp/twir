import {
  ActionIcon, Avatar,
  Button,
  Divider,
  Flex,
  Group,
  NumberInput,
  Select,
  SelectItem,
  Switch,
  Tabs,
  Text,
  TextInput,
  Autocomplete,
  MultiSelect, ScrollArea, Menu, Modal, UnstyledButton,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useDebouncedState, useDebouncedValue } from '@mantine/hooks';
import { openModal, closeAllModals } from '@mantine/modals';
import { showNotification } from '@mantine/notifications';
import { IconAdjustmentsHorizontal, IconPlus, IconTrash, IconUsers, IconX } from '@tabler/icons';
import { YoutubeSettings, SearchResult } from '@tsuwari/types/api';
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

  const search = youtube.useSearch();
  const [searchType, setSearchType] = useState<'video' | 'channel'>('video');
  const [searchValue, setSearchValue] = useDebouncedState('', 200);
  const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
  const [filterSongs, setFilterSongs] = useState('');

  useEffect(() => {
    if (searchValue) {
      search.mutateAsync({ query: searchValue, type: searchType }).then(data => {
        setSearchResults(data);
      });
    } else {
      setSearchResults([]);
    }
  }, [searchValue]);

  const [searchModalOpened, setSearchModalOpened] = useState(false);
  function insertBySearch(r: SearchResult) {
    const notificationOpts = {
      title: 'Already exists',
      color: 'red',
      message: `${r.title} already in list`,
    };

    if (searchType === 'video') {
      if (form.values.blackList.songs.some(s => s.id === r.id)) {
        showNotification(notificationOpts);
      } else {
        form.insertListItem('blackList.songs', r);
        setSearchModalOpened(false);
      }
    }

    if (searchType === 'channel') {
      if (form.values.blackList.channels.some(s => s.id === r.id)) {
        showNotification(notificationOpts);
      } else {
        form.insertListItem('blackList.channels', r);
        setSearchModalOpened(false);
      }
    }

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

    <Text size='lg'>Restrictions</Text>
    <Tabs defaultValue="users">
      <Tabs.List position='center' grow>
        <Tabs.Tab color='pink' value="users" icon={<IconUsers size={14} />}>Users</Tabs.Tab>
        <Tabs.Tab color='grape' value="songs" icon={<IconUsers size={14} />}>Songs</Tabs.Tab>
        <Tabs.Tab color='violet' value="channels" icon={<IconUsers size={14} />}>Channels</Tabs.Tab>
      </Tabs.List>

      <Tabs.Panel value="users" pt="xs">
        <NumberInput label="Maximum songs by user in queue" {...form.getInputProps('user.maxRequests')} />
        <NumberInput label="Minimal watch time of user for request song" {...form.getInputProps('user.minWatchTime')} />
        <NumberInput label="Minimal messages by user for request song" {...form.getInputProps('user.minMessages')} />
        <NumberInput label="Minimal follow time for request song" {...form.getInputProps('user.minFollowTime')} />

      </Tabs.Panel>
      <Tabs.Panel value="songs" pt="xs">
        <NumberInput label="Max length of song for request" {...form.getInputProps('song.maxLength')} />
        <NumberInput label="Minimal views on song for request" {...form.getInputProps('song.minViews')} />

        <Divider style={{ marginTop: 10 }} />

        <Flex direction='row' justify='space-between' style={{ marginTop: 10 }}>
          <Text size='sm'>Denied songs for request</Text>
          <ActionIcon
            color='green'
            variant={'filled'}
            size={'sm'}
            onClick={() => {
              setSearchType('video');
              setSearchModalOpened(true);
            }}
          ><IconPlus /></ActionIcon>
        </Flex>

        {form.values.blackList.songs.length
          ? <TextInput style={{ marginTop: 10 }} placeholder='filter...' onChange={(v) => setFilterSongs(v.target.value)} />
          : ''
        }

        <ScrollArea type={'always'}>
          <div style={{ maxHeight: 300 }}>
            {form.values.blackList.songs.length ?
              form.values.blackList.songs
                .filter(s => s.title.toLowerCase().includes(filterSongs.toLowerCase()))
                .map((s, i) => <Group style={{ maxHeight: 280 }}>
                <Flex
                  key={s.id}
                  direction='row'
                  justify='space-between'
                  style={{ width: '95%', marginTop: 10 }}
                  gap='sm'
                >
                  <Avatar size={40} src={s.thumbNail} />
                  <Text size={'sm'} lineClamp={4}>{s.title}</Text>
                  <ActionIcon onClick={() => form.removeListItem('blackList.songs', i)}>
                    <IconX />
                  </ActionIcon>
                </Flex>
              </Group>)
              : ''}
          </div>
        </ScrollArea>

      </Tabs.Panel>
      <Tabs.Panel value="channels" pt="xs">
        <Flex direction='row' justify='space-between' style={{ marginTop: 10 }}>
          <Text size='sm'>Denied channels for request</Text>
          <ActionIcon
            color='green'
            variant={'filled'}
            size={'sm'}
            onClick={() => {
              setSearchType('channel');
              setSearchModalOpened(true);
            }}
          ><IconPlus /></ActionIcon>
        </Flex>

        <ScrollArea type={'always'}>
          <div style={{ maxHeight: 300 }}>
            {form.values.blackList.channels.length ?
              form.values.blackList.channels
                .filter(s => s.title.toLowerCase().includes(filterSongs.toLowerCase()))
                .map((s, i) => <Group style={{ maxHeight: 280 }}>
                  <Flex
                    key={s.id}
                    direction='row'
                    justify='space-between'
                    style={{ width: '95%', marginTop: 10 }}
                    gap='sm'
                  >
                    <Avatar size={40} src={s.thumbNail} />
                    <Text size={'sm'} lineClamp={4}>{s.title}</Text>
                    <ActionIcon onClick={() => form.removeListItem('blackList.channels', i)}>
                      <IconX />
                    </ActionIcon>
                  </Flex>
                </Group>)
              : ''}
          </div>
        </ScrollArea>
      </Tabs.Panel>
    </Tabs>

    <Divider style={{ marginTop: 10, marginBottom: 10 }} />
    <Button color='green' onClick={submit}>Save</Button>

    <Modal
      opened={searchModalOpened}
      onClose={() => setSearchModalOpened(false)}
      title="Search"
      zIndex={300}
    >
      <TextInput
        placeholder={'search...'}
        onChange={(v) => setSearchValue(v.target.value)}
        style={{ marginBottom: 10 }}
      />

      {searchResults.length
        ? searchResults.map(r => <UnstyledButton
          onClick={() => {
            insertBySearch(r);
          }}
        >
          <Flex key={r.id} direction={'row'} gap={'md'}>
            <Avatar size={40} src={r.thumbNail} />
            <Text size={'sm'}>{r.title}</Text>
          </Flex>
        </UnstyledButton>)
        : ''
      }
    </Modal>
  </form>;
};
