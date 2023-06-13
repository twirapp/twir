import { Button, Flex, Grid, Text } from '@mantine/core';
import { closeAllModals } from '@mantine/modals';
import { showNotification } from '@mantine/notifications';
import { GetServerSideProps, NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import React, { useEffect } from 'react';

import { noop } from '../../util/chore';

import {
  useYouTubeSettingsForm,
  YouTubeSettingsFormProvider,
} from '@/components/song-requests/settings/form';
import { YouTubeGeneralSettings } from '@/components/song-requests/settings/general';
import { YouTubeSongsSettings } from '@/components/song-requests/settings/songs';
import { YouTubeTranslationSettings } from '@/components/song-requests/settings/translations';
import { YouTubeUsersSettings } from '@/components/song-requests/settings/users';
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
    ...(await serverSideTranslations(locale!, ['song-requests-settings', 'layout'])),
  },
});

const Settings: NextPage = () => {
  const form = useYouTubeSettingsForm({
    initialValues: {
      enabled: false,
      acceptOnlyWhenOnline: true,
      channelPointsRewardId: '',
      maxRequests: 500,
      announcePlay: true,
      neededVotesVorSkip: 30,
      denyList: {
        artistsNames: [],
        songs: [],
        users: [],
        channels: [],
      },
      song: {
        maxLength: 10,
        minLength: 0,
        minViews: 50000,
        acceptedCategories: [],
      },
      user: {
        maxRequests: 20,
        minWatchTime: 0,
        minFollowTime: 0,
        minMessages: 0,
      },
      translations: {
        notEnabled: 'Song requests not enabled.',
        nowPlaying: 'Now playing "{{songTitle}}" {{songLink} requested by @{{orderedByDisplayName}}',
        noText: 'You should provide text for song request.',
        acceptOnlyWhenOnline: 'Requests accepted only on online streams.',
        song: {
          notFound: 'Song not found.',
          alreadyInQueue: 'Song already in queue.',
          ageRestrictions: 'Age restriction on that song.',
          cannotGetInformation: 'Cannot get information about song.',
          live: 'Seems like that song is live, which is disallowed.',
          denied: 'That song is denied for requests.',
          requestedMessage: 'Song "{{songTitle}}" requested, queue position {{position}}. Estimated wait time before your track will be played is {{waitTime}}.',
          maximumOrdered: 'Maximum number of songs is queued ({{maximum}}).',
          minViews: 'Song {{songTitle}} ({{songViews}} views) haven\'t {{neededViews}} views for being ordered',
          maxLength: 'Maximum length of song is {{maxLength}}',
          minLength: 'Minimum length of song is {{minLength}}',
        },
        user: {
          denied: 'You are denied to request any song.',
          maxRequests: 'Maximum number of songs ordered by you ({{count}})',
          minMessages: 'You have only {{userMessages}} messages, but needed {{neededMessages}} for requesting song',
          minWatched: 'You\'ve only watched {{userWatched}} but needed {{neededWatched}} to request a song.',
          minFollow: 'You are followed for {{userFollow}} minutes, but needed {{neededFollow}} for requesting song',
        },
        channel: {
          denied: 'That channel is denied for requests.',
        },
      },
    },
    validate: {
      maxRequests: (v) => {
        if (v > 500) return 'Max number of songs in queue is 500';
        if (v <= 0) return 'Minimum number cannot be lower then 0';
      },
      song: {
        maxLength: (v) => {
          if (v < 0) return 'Minimal number cannot be lower then 0';
        },
        minViews: (v) => {
          if (v < 0) return 'Minimal number cannot be lower then 0';
        },
      },
      user: {
        maxRequests: (v) => {
          if (v < 0) return 'Minimal number cannot be lower then 0';
        },
        minWatchTime: (v) => {
          if (v < 0) return 'Minimal number cannot be lower then 0';
        },
        minFollowTime: (v) => {
          if (v < 0) return 'Minimal number cannot be lower then 0';
        },
        minMessages: (v) => {
          if (v < 0) return 'Minimal number cannot be lower then 0';
        },
      },
    },
  });

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

  return (
    <YouTubeSettingsFormProvider form={form}>
      <form>
        <Flex justify={'space-between'}>
          <Text size={'lg'}>Song Requests settings</Text>
          <Button color={'green'} onClick={submit}>
            Save
          </Button>
        </Flex>
        <Grid justify={'center'} style={{ marginTop: 10 }}>
          <Grid.Col {...cols}>
            <YouTubeGeneralSettings />
          </Grid.Col>
          <Grid.Col {...cols}>
            <YouTubeUsersSettings />
          </Grid.Col>
          <Grid.Col {...cols}>
            <YouTubeSongsSettings />
          </Grid.Col>
          <Grid.Col span={12}>
            <YouTubeTranslationSettings />
          </Grid.Col>
        </Grid>
      </form>
    </YouTubeSettingsFormProvider>
  );
};

export default Settings;
