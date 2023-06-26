import {
  ActionIcon,
  Card,
  Divider,
  Flex,
  NumberInput,
  ScrollArea,
  Text,
  TextInput,
} from '@mantine/core';
import { useDebouncedState } from '@mantine/hooks';
import { IconPlus, IconX } from '@tabler/icons';
import { SearchResult } from '@twir/types/api';
import { useTranslation } from 'next-i18next';
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
  const [t] = useTranslation('song-requests-settings');

  useEffect(() => {
    if (newIgnoreSongSearch) {
      search.mutateAsync({ query: newIgnoreSongSearch, type: 'video' }).then((data) => {
        setSearchResults(data);
      });
    } else {
      setSearchResults([]);
    }
  }, [newIgnoreSongSearch]);

  return (
    <Card style={{ minHeight: 500 }}>
      <Card.Section p={'xs'}>
        <Text>{t('songs.title')}</Text>
      </Card.Section>
      <Divider />
      <Card.Section p={'md'}>
        <Flex direction={'column'} gap={'xs'}>
          <NumberInput label={t('songs.maxRequests')} {...form.getInputProps('maxRequests')} />
          <NumberInput label={t('songs.minLength')} {...form.getInputProps('song.minLength')} />
          <NumberInput label={t('songs.maxLength')} {...form.getInputProps('song.maxLength')} />
          <NumberInput label={t('songs.minViews')} {...form.getInputProps('song.minViews')} />
        </Flex>

        <Divider style={{ marginTop: 10 }} />

        <Flex direction="row" justify="space-between" style={{ marginTop: 10 }}>
          <Text size="sm">{t('songs.denied')}</Text>
          <ActionIcon
            onClick={() => setAddingNewIgnoreSong(!addingNewIgnoreSong)}
            color={'green'}
            size={'sm'}
          >
            <IconPlus />
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
                ? searchResults.map((r) => (
                    <YouTubeSettingsListButtonButton
                      key={r.id}
                      text={r.title}
                      image={r.thumbNail}
                      onClick={() => {
                        form.insertListItem('denyList.songs', r);
                        setAddingNewIgnoreSong(false);
                        setSearchResults([]);
                      }}
                    />
                  ))
                : ''}
            </Flex>
          </ScrollArea>
        </Flex>

        <Flex hidden={addingNewIgnoreSong} direction={'column'}>
          {form.values.denyList.songs.length ? (
            <TextInput
              style={{ marginTop: 10 }}
              placeholder="filter..."
              onChange={(v) => setFilterSongs(v.target.value)}
            />
          ) : (
            ''
          )}

          <ScrollArea type={'always'} style={{ marginTop: 10 }}>
            <Flex direction={'column'} style={{ maxHeight: 300 }} gap={'sm'}>
              {form.values.denyList.songs.length
                ? form.values.denyList.songs
                    .filter((c) => c.title.toLowerCase().includes(filterSongs))
                    .map((c, i) => (
                      <YouTubeSettingsListButtonButton
                        key={c.id}
                        image={c.thumbNail}
                        text={c.title}
                        onClick={() => form.removeListItem('denyList.songs', i)}
                        icon={IconX}
                      />
                    ))
                : ''}
            </Flex>
          </ScrollArea>
        </Flex>
      </Card.Section>
    </Card>
  );
};
