import {
  ActionIcon,
  Anchor,
  Card,
  Center,
  Group,
  Loader,
  ScrollArea,
  Table,
  Text,
  Tooltip,
} from '@mantine/core';
import { IconGripVertical, IconTrash } from '@tabler/icons';
import dynamic from 'next/dynamic';
import { useContext } from 'react';

import { resolveUserName } from '../../../util/resolveUserName';
import { convertMillisToTime } from '../helpers';

import { confirmDelete } from '@/components/confirmDelete';
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
  const { videos, reorderVideos, skipVideo, clearQueue } = useContext(PlayerContext);
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
            <Center>{index + 1}</Center>
          </td>
          <td>
            <Anchor target="_blank" href={'https://youtu.be/' + video.videoId}>
              {video.title}
            </Anchor>
          </td>
          <td>
            <Center>{resolveUserName(video.orderedByName, video.orderedByDisplayName)}</Center>
          </td>
          <td>
            <Center>{convertMillisToTime(video.duration)}</Center>
          </td>
          <td>
            <ActionIcon mx="sm" variant="transparent" color="red" onClick={() => skipVideo(index)}>
              <IconTrash size={14} />
            </ActionIcon>
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
          <Tooltip withinPortal position="top" label="Clear queue">
            <ActionIcon
              onClick={() =>
                confirmDelete({
                  onConfirm: () => clearQueue(),
                  text: 'Are you sure you wanna clear song list?',
                  title: 'Song requests',
                })
              }
            >
              <IconTrash size={14} />
            </ActionIcon>
          </Tooltip>
        </Group>
      </Card.Section>
      <Card.Section className={cardClasses.card} style={{ overflow: 'auto' }}>
        {items.length ? (
          <DragDropContext
            onDragEnd={({ destination, source }) => {
              reorderVideos(destination!, source);
            }}
          >
            <ScrollArea.Autosize maxHeight="80vh">
              <Table highlightOnHover>
                <thead className={draggableClasses.thead}>
                  <tr>
                    <th />
                    <th>
                      <Center>#</Center>
                    </th>
                    <th>Title</th>
                    <th>
                      <Center>Author</Center>
                    </th>
                    <th>
                      <Center>Duration</Center>
                    </th>
                    <th />
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
                <tfoot className={draggableClasses.tfoot}>
                  <tr>
                    <th />
                    <th>
                      <Center>{videos.length}</Center>
                    </th>
                    <th />
                    <th />
                    <th>
                      <Center>
                        {convertMillisToTime(
                          videos.reduce((acc, curr) => acc + curr.duration, 0) ?? 0,
                        )}
                      </Center>
                    </th>
                    <th />
                  </tr>
                </tfoot>
              </Table>
            </ScrollArea.Autosize>
          </DragDropContext>
        ) : (
          <Center py="lg">
            <Loader variant="dots" />
          </Center>
        )}
      </Card.Section>
    </Card>
  );
};
