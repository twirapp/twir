import { Table } from '@mantine/core';
import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';

type Props = {
  videos: RequestedSong[]
}

export const VideosList: React.FC<Props> = (props) => {
  return (
    <Table>
      <thead>
      <tr>
        <th>#</th>
        <th>Title</th>
        <th>Requested by</th>
      </tr>
      </thead>
      <tbody>
      {props.videos?.map((video, index) => (
        <tr key={video.id}>
          <th>{index + 1}</th>
          <th><a href={'https://youtu.be/' + video.videoId}>{video.title}</a></th>
          <th>{video.orderedByName}</th>
        </tr>
      ))}
      </tbody>
    </Table>
  );
};
