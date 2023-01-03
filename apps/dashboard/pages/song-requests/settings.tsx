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
      enabled: true,
      acceptOnlyWhenOnline: true,
      channelPointsRewardId: '',
      maxRequests: 500,
      announcePlay: true,
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
      // for (const error of Object.values(validation.errors).flat(10) as string[]) {
      //   showNotification({
      //     title: 'Validation error',
      //     color: 'red',
      //     message: error,
      //   });
      //   console.log(error);
      // }

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
        </Grid>
      </form>
    </YouTubeSettingsFormProvider>
  );
};

export default Settings;
