import { Anchor, Badge, Table, Text } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { NextPage } from 'next';
import { useRouter } from 'next/router';

import { useUsersByNames } from '@/services/users';

type Song = {
  title: string
  videoId: string
  duration: number
  orderedByName: string
  orderedByDisplayName: string | null
}

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

  const {
    data: songs,
  } = useQuery({
    queryKey: ['commands', users?.at(0)?.id],
    queryFn: async (): Promise<Song[]> => {
      const req = await fetch(`/api/v1/p/song-requests/${users?.at(0)?.id}`);

      return req.json();
    },
    enabled: !!users?.at(0)?.id,
  });

  return (<Table highlightOnHover>
    <thead>
    <tr>
      <th>#</th>
      <th>Title</th>
      <th>Author</th>
      <th>Duration</th>
    </tr>
    </thead>
    <tbody>
    {songs?.map((c, i) => <tr>
      <td>{i + 1}</td>
      <td><Anchor href={'https://youtu.be/' + c.videoId} target={'_blank'}>{c.title}</Anchor></td>
      <td>{c.orderedByDisplayName ?? c.orderedByName}</td>
      <td>{convertMillisToTime(c.duration)}</td>
    </tr>)}
    </tbody>
  </Table>);
};

export default SongRequests;