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
  MultiSelect, ScrollArea, Menu,
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
    <Tabs defaultValue="users" onTabChange={(t) => {
      if (t === 'songs') setSearchType('video');
      if (t === 'channels') setSearchType('channel');
    }}>
      <Tabs.List position='center'>
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

        <Menu shadow="md" width={400} position={'top'}>
          <Menu.Target>
            <Flex direction='row' justify='space-between' style={{ marginTop: 10 }}>
              <Text size='sm'>Denied songs for request</Text>
              <ActionIcon color='green' variant={'filled'} size={'sm'}><IconPlus /></ActionIcon>
            </Flex>
          </Menu.Target>

          <Menu.Dropdown>
            {searchResults.length ? searchResults.map(r =>
              <Menu.Item 
                icon={<Avatar size={40} src={r.thumbNail} />}
                onClick={() => {
                  if (form.values.blackList.songs.some(s => s.id === r.id)) {
                    showNotification({
                      title: 'Song exists',
                      color: 'red',
                      message: `Song ${r.title} already added to the ignore list!`,
                    });
                  } else {
                    form.insertListItem('blackList.songs', r);
                  }
                }}
              >
                {r.title}
              </Menu.Item>)
              : ''}

            <TextInput
              placeholder={'filter...'}
              onChange={(v) => setSearchValue(v.target.value)}
            />

          </Menu.Dropdown>
        </Menu>

        {form.values.blackList.songs.length
          ? <TextInput style={{ marginTop: 10 }} placeholder='search...' onChange={(v) => setFilterSongs(v.target.value)} />
          : ''
        }

        <ScrollArea type={'auto'}>
          {form.values.blackList.songs.length ?
            form.values.blackList.songs
              .filter(s => s.title.includes(filterSongs))
              .map((s, i) => <Group style={{ maxHeight: 280 }}>
              <Flex
                direction='row'
                justify='space-between'
                style={{ width: '95%', marginTop: 10 }}
                gap='sm'
              >
                <Avatar size={40} color="blue" src={s.thumbNail} />
                <Text size={'sm'} lineClamp={4}>{s.title}</Text>
                <ActionIcon onClick={() => form.removeListItem('blackList.songs', i)}>
                  <IconX />
                </ActionIcon>
              </Flex>
            </Group>)
            : ''}
        </ScrollArea>

      </Tabs.Panel>
      <Tabs.Panel value="channels" pt="xs">
        Channels
      </Tabs.Panel>
    </Tabs>

    <Divider style={{ marginTop: 10, marginBottom: 10 }} />
    <Button color='green' onClick={submit}>Save</Button>
  </form>;
};
