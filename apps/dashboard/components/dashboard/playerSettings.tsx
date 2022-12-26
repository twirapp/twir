import { useForm } from '@mantine/form';
import { YoutubeSettings } from '@tsuwari/types/api';

export const PlayerSettings: React.FC  = () => {
  const form = useForm<YoutubeSettings>({
    initialValues: {
      acceptOnlyWhenOnline: true,
      channelPointsRewardName: '',
      maxRequests: 20,
      enabled: true,
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
  });

  return <div></div>;
};