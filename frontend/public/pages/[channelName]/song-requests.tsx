import { Anchor, Table } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { NextPage } from 'next';
import { useRouter } from 'next/router';
import { useState } from 'react';

import { useUsersByNames } from '@/services/users';

type Song = {
  title: string;
  videoId: string;
  duration: number;
  orderedByName: string;
  orderedByDisplayName: string | null;
};

const padTo2Digits = (num: number) => {
  return num.toString().padStart(2, '0');
};

const convertMillisToTime = (millis: number) => {
  let seconds = Math.floor(millis / 1000);
  let minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);

  seconds = seconds % 60;
  minutes = minutes % 60;

  return `${hours ? `${padTo2Digits(hours)}:` : ''}${padTo2Digits(minutes)}:${padTo2Digits(
    seconds,
  )}`;
};

const SongRequests: NextPage = () => {
  const router = useRouter();
  const { data: users } = useUsersByNames([router.query.channelName as string]);
  const [userId] = useState(() => users?.at(0)?.id);

  const { data: songs } = useQuery({
    queryKey: ['commands', userId],
    queryFn: async (): Promise<Song[]> => {
      const req = await fetch(`/api/v1/p/song-requests/${userId}`);

      return req.json();
    },
    initialData: [],
    enabled: !!userId,
  });

  return (
    <Table highlightOnHover>
      <thead>
        <tr>
          <th>#</th>
          <th>Title</th>
          <th>Author</th>
          <th>Duration</th>
        </tr>
      </thead>
      <tbody>
        {songs.map((song, key) => (
          <tr key={song.videoId}>
            <td>{key + 1}</td>
            <td>
              <Anchor href={'https://youtu.be/' + song.videoId} target="_blank">
                {song.title}
              </Anchor>
            </td>
            <td>{song.orderedByDisplayName ?? song.orderedByName}</td>
            <td>{convertMillisToTime(song.duration)}</td>
          </tr>
        ))}
      </tbody>
    </Table>
  );
};

export default SongRequests;
