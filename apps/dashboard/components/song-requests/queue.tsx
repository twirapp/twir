import { ActionIcon, Badge, Card, Flex, Group, Table, Text } from '@mantine/core';
import { IconGripVertical, IconTrash } from '@tabler/icons';
import dynamic from 'next/dynamic';
import { useContext } from 'react';

import { convertMillisToTime } from './helpers';

import { PlayerContext } from '@/components/song-requests/context';
import { useCardStyles } from '@/styles/card';
import { useDraggableStyles } from '@/styles/draggable';

const DragDropContext = dynamic(
  async () => {
    const mod = await import('react-beautiful-dnd');
    return mod.DragDropContext;
  },
  { ssr: false },
);

const Droppable = dynamic(
  async () => {
    const mod = await import('react-beautiful-dnd');
    return mod.Droppable;
  },
  { ssr: false },
);

const Draggable = dynamic(
  async () => {
    const mod = await import('react-beautiful-dnd');
    return mod.Draggable;
  },
  { ssr: false },
);

export const QueueList: React.FC = () => {
  const { videos, reorderVideos, skipVideo } = useContext(PlayerContext);
  const { classes: draggableClasses } = useDraggableStyles();
  const { classes: cardClasses } = useCardStyles();

  const items = videos.map((video, index) => (
    <Draggable
      key={video.id}
      disableInteractiveElementBlocking={index === 0}
      isDragDisabled={index === 0}
      index={index}
      draggableId={video.id}
    >
      {(provided) => (
        <tr className={draggableClasses.item} ref={provided.innerRef} {...provided.draggableProps}>
          <td>
            <div className={draggableClasses.dragHandle} {...provided.dragHandleProps}>
              <IconGripVertical size={18} stroke={1.5} />
            </div>
          </td>
          <td>
            <Badge color={index === 0 ? 'green' : 'gray'} radius="md" variant="filled">
              {index === 0 ? 'Playing' : index}
            </Badge>
          </td>
          <td>
            <a href={'https://youtu.be/' + video.videoId}>{video.title}</a>
          </td>
          <td>{video.orderedByName}</td>
          <td>{convertMillisToTime(video.duration)}</td>
          <td>
            <Flex>
              <ActionIcon variant={'filled'} color={'red'} onClick={() => skipVideo(index)}>
                <IconTrash size={14} />
              </ActionIcon>
            </Flex>
          </td>
        </tr>
      )}
    </Draggable>
  ));

  return (
    <Card withBorder radius="md" p="md">
      <Card.Section withBorder inheritPadding py="xs">
        <Group position="apart">
          <Text weight={500}>Queue</Text>
        </Group>
      </Card.Section>
      <Card.Section className={cardClasses.card}>
        <DragDropContext
          onDragEnd={({ destination, source }) => {
            reorderVideos(destination!, source);
          }}
        >
          <Table>
            <thead>
              <tr>
                <th></th>
                <th>#</th>
                <th>Title</th>
                <th>Requested by</th>
                <th>
                  Duration (
                  {convertMillisToTime(
                    videos.slice(1).reduce((acc, curr) => acc + curr.duration, 0) ?? 0,
                  )}
                  )
                </th>
                <th>Actions</th>
              </tr>
            </thead>
            <Droppable droppableId="dnd-list" direction="vertical">
              {(provided) => (
                <tbody {...provided.droppableProps} ref={provided.innerRef}>
                  {items}
                  {provided.placeholder}
                </tbody>
              )}
            </Droppable>
          </Table>
        </DragDropContext>
      </Card.Section>
    </Card>
  );
};
