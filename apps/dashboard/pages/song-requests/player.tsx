'use client';

import type { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { GetServerSideProps, NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import dynamic from 'next/dynamic';
import { useState } from 'react';

import { PlayerContext } from '@/components/song-requests/context';

export const getServerSideProps: GetServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale!, ['layout'])),
  },
});

const PlayerComponent = dynamic(
  () => import('../../components/song-requests/player'),
  { ssr: false },
);

const Player: NextPage = () => {
  const [videos, setVideos] = useState<RequestedSong[]>([]);

  return <PlayerContext.Provider value={{ videos, setVideos }}><PlayerComponent/></PlayerContext.Provider>;
};

export default Player;