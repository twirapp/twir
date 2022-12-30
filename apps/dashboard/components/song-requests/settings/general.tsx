import {
  ActionIcon,
  Card,
  Divider,
  Flex,
  ScrollArea,
  Select,
  Switch,
  Text,
  TextInput,
  useMantineTheme,
} from '@mantine/core';
import { useDebouncedState } from '@mantine/hooks';
import { IconPlus, IconX } from '@tabler/icons';
import { SearchResult } from '@tsuwari/types/api';
import React, { useEffect, useState } from 'react';

import { useYouTubeSettingsFormContext } from '@/components/song-requests/settings/form';
import { YouTubeSettingsListButtonButton } from '@/components/song-requests/settings/listButton';
import { RewardItem, RewardItemProps } from '@/components/song-requests/settings/reward';
import { useRewards } from '@/services/api';
import { useYoutubeModule } from '@/services/api/modules';

export const YouTubeGeneralSettings: React.FC = () => {
  const theme = useMantineTheme();
  const form = useYouTubeSettingsFormContext();
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
  const search = youtube.useSearch();

  const [filterChannels, setFilterChannels] = useState('');
  const [addingNewIgnoreChannel, setAddingNewIgnoreChannel] = useState(false);
  const [newIgnoreChannelSearch, setNewIgnoreChannelSearch] = useDebouncedState('', 200);
  const [searchResults, setSearchResults] = useState<SearchResult[]>([]);

  useEffect(() => {
    if (newIgnoreChannelSearch) {
      search.mutateAsync({ query: newIgnoreChannelSearch, type: 'channel' }).then(data => {
        setSearchResults(data);
      });
    } else {
      setSearchResults([]);
    }
  }, [newIgnoreChannelSearch]);

  return (
    <Card style={{ minHeight: 500 }}>
      <Card.Section p={'xs'}><Text>General</Text></Card.Section>
      <Divider/>
      <Card.Section p={'md'}>
        <Flex direction={'column'} gap={'xs'}>
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
            label="Channel reward for requesting songs"
            placeholder="..."
            searchable
            itemComponent={RewardItem}
            dropdownPosition={'bottom'}
            allowDeselect
            data={rewards}
            {...form.getInputProps('channelPointsRewardId')}
          />

        </Flex>

        <Divider style={{ marginTop: 10 }}/>

        <Flex direction="row" justify="space-between" style={{ marginTop: 10 }}>
          <Text size="sm">Denied channels for request</Text>
          <ActionIcon
            onClick={() => setAddingNewIgnoreChannel(!addingNewIgnoreChannel)}
            color={'green'}
            size={'sm'}
          >
            <IconPlus/>
          </ActionIcon>
        </Flex>

        <Flex hidden={!addingNewIgnoreChannel} direction={'column'}>
          <TextInput
            placeholder={'search...'}
            onChange={(v) => setNewIgnoreChannelSearch(v.currentTarget.value)}
            style={{ marginBottom: 10 }}
          />

          <ScrollArea type={'always'} style={{ marginTop: 10 }}>
            <Flex direction={'column'} style={{ maxHeight: 300 }} gap={'sm'}>
              {searchResults.length
                ? searchResults.map((r) => <YouTubeSettingsListButtonButton
                  key={r.id}
                  text={r.title}
                  image={r.thumbNail}
                  onClick={() => {
                    form.insertListItem('denyList.channels', r);
                    setAddingNewIgnoreChannel(false);
                  }}
                />)
                : ''
              }

            </Flex>
          </ScrollArea>
        </Flex>

        <Flex hidden={addingNewIgnoreChannel} direction={'column'}>
          {form.values.denyList.channels.length
            ? <TextInput style={{ marginTop: 10 }} placeholder="filter..."
                         onChange={(v) => setFilterChannels(v.target.value)}
            />
            : ''
          }

          <ScrollArea type={'always'} style={{ marginTop: 10 }}>
            <Flex direction={'column'} style={{ maxHeight: 300 }}>
              {form.values.denyList.channels.length
                ? form.values.denyList.channels
                  .filter(c => c.title.includes(filterChannels))
                  .map((c, i) => <YouTubeSettingsListButtonButton
                    key={c.id}
                    image={c.thumbNail}
                    text={c.title}
                    onClick={() => form.removeListItem('denyList.channels', i)}
                    icon={IconX}
                  />)
                : ''
              }
            </Flex>
          </ScrollArea>
        </Flex>

      </Card.Section>
    </Card>
  );
};