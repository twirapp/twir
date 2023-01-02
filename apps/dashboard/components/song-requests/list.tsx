import { ActionIcon, Flex, Table } from '@mantine/core';
import { IconTrash } from '@tabler/icons';
import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';

import { millisToMinutesAndSeconds } from './helpers';

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
        <th>Duration</th>
        <th>Actions</th>
      </tr>
      </thead>
      <tbody>
      {props.videos?.slice(1).map((video, index) => (
        <tr key={video.id} style={{ textAlign: 'left' }}>
          <td>{index + 1}</td>
          <td><a href={'https://youtu.be/' + video.videoId}>{video.title}</a></td>
          <td>{video.orderedByName}</td>
          <td>{millisToMinutesAndSeconds(video.duration)}</td>
          <td>
            <Flex>
              <ActionIcon variant={'filled'} color={'red'}><IconTrash size={14}/></ActionIcon>
            </Flex>
          </td>
        </tr>
      ))}
      </tbody>
    </Table>
  );
};
