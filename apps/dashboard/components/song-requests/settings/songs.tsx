import { ActionIcon, Card, Divider, Flex, NumberInput, ScrollArea, Text, TextInput } from '@mantine/core';
import { useDebouncedState } from '@mantine/hooks';
import { IconPlus, IconX } from '@tabler/icons';
import { SearchResult } from '@tsuwari/types/api';
import React, { useEffect, useState } from 'react';

import { useYouTubeSettingsFormContext } from '@/components/song-requests/settings/form';
import { YouTubeSettingsListButtonButton } from '@/components/song-requests/settings/listButton';
import { useYoutubeModule } from '@/services/api/modules';


export const YouTubeSongsSettings: React.FC = () => {
  const form = useYouTubeSettingsFormContext();

  const youtube = useYoutubeModule();
  const search = youtube.useSearch();
  const [filterSongs, setFilterSongs] = useState('');
  const [addingNewIgnoreSong, setAddingNewIgnoreSong] = useState(false);
  const [newIgnoreSongSearch, setNewIgnoreSongSearch] = useDebouncedState('', 200);
  const [searchResults, setSearchResults] = useState<SearchResult[]>([]);

  useEffect(() => {
    if (newIgnoreSongSearch) {
      search.mutateAsync({ query: newIgnoreSongSearch, type: 'video' }).then(data => {
        setSearchResults(data);
      });
    } else {
      setSearchResults([]);
    }
  }, [newIgnoreSongSearch]);

  return (
    <Card style={{ minHeight: 500 }}>
      <Card.Section p={'xs'}><Text>Songs</Text></Card.Section>
      <Divider/>
      <Card.Section p={'md'}>
        <Flex direction={'column'} gap={'xs'}>
          <NumberInput label="Maximum number of songs in queue" {...form.getInputProps('maxRequests')} />
          <NumberInput label="Max length of song for request (minutes)" {...form.getInputProps('song.maxLength')} />
          <NumberInput label="Minimal views on song for request" {...form.getInputProps('song.minViews')} />
        </Flex>

        <Divider style={{ marginTop: 10 }}/>

        <Flex direction="row" justify="space-between" style={{ marginTop: 10 }}>
          <Text size="sm">Denied songs for request</Text>
          <ActionIcon
            onClick={() => setAddingNewIgnoreSong(!addingNewIgnoreSong)}
            color={'green'}
            size={'sm'}
          >
            <IconPlus/>
          </ActionIcon>
        </Flex>

        <Flex hidden={!addingNewIgnoreSong} direction={'column'}>
          <TextInput
            placeholder={'search...'}
            onChange={(v) => setNewIgnoreSongSearch(v.currentTarget.value)}
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
                    form.insertListItem('denyList.songs', r);
                    setAddingNewIgnoreSong(false);
                    setSearchResults([]);
                  }}
                />)
                : ''
              }

            </Flex>
          </ScrollArea>
        </Flex>

        <Flex hidden={addingNewIgnoreSong} direction={'column'}>
          {form.values.denyList.songs.length
            ? <TextInput style={{ marginTop: 10 }} placeholder="filter..."
                         onChange={(v) => setFilterSongs(v.target.value)}
            />
            : ''
          }

          <ScrollArea type={'always'} style={{ marginTop: 10 }}>
            <Flex direction={'column'} style={{ maxHeight: 300 }} gap={'sm'}>
              {form.values.denyList.songs.length
                ? form.values.denyList.songs
                  .filter(c => c.title.toLowerCase().includes(filterSongs))
                  .map((c, i) => <YouTubeSettingsListButtonButton
                    key={c.id}
                    image={c.thumbNail}
                    text={c.title}
                    onClick={() => form.removeListItem('denyList.songs', i)}
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