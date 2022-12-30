import { GetServerSideProps, NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useEffect, useState } from 'react';
import YouTube, { YouTubeProps } from 'react-youtube';
import type { YouTubePlayer } from 'youtube-player/dist/types';

export const getServerSideProps: GetServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale!, ['layout'])),
  },
});

const Player: NextPage = () => {
  const [player, setPlayer] = useState<YouTubePlayer>();

  const onPlayerReady: YouTubeProps['onReady'] = (event) => {
    console.log('player ready');
    setPlayer(event.target);
  };

  const onStateChange: YouTubeProps['onStateChange'] = (event) => {
    console.log(event);
  };

  useEffect(() => {
    if (!player) return;

    for (const song of ['FhwsDelZPV8', 'Bu9fAfD5YKk']) {
      player.cueVideoById(song, 0);
    }

    setTimeout(() => {
      player!.playVideoAt(0);
    }, 5000);
  }, [player]);

  return <YouTube onReady={onPlayerReady} onStateChange={onStateChange}/>;

};

export default Player;