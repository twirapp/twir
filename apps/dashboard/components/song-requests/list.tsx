import { Table } from '@mantine/core';
import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';

import { formatDuration } from './helpers';

type Props = {
  videos: RequestedSong[]
}

export const VideosList: React.FC<Props> = (props) => {
  return (
    <Table>
      <thead>
      <tr style={{ textAlign: 'left' }}>
        <th>#</th>
        <th>Title</th>
        <th>Requested by</th>
        <th>Duration</th>
      </tr>
      </thead>
      <tbody>
      {props.videos?.map((video, index) => (
        <tr key={video.id} style={{ textAlign: 'left' }}>
          <th>{index + 1}</th>
          <th><a href={'https://youtu.be/' + video.videoId}>{video.title}</a></th>
          <th>{video.orderedByName}</th>
          <th>{formatDuration(video.duration)}</th>
        </tr>
      ))}
      </tbody>
    </Table>
  );
};
