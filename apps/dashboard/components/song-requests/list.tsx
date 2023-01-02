import { ActionIcon, createStyles, Flex, Table } from '@mantine/core';
import { useListState } from '@mantine/hooks';
import { IconGripVertical, IconTrash } from '@tabler/icons';
import { RequestedSong } from '@tsuwari/typeorm/entities/RequestedSong';
import { useCallback, useContext, useEffect, useState } from 'react';
import { DragDropContext, Draggable, DraggableLocation, Droppable } from 'react-beautiful-dnd';

import { millisToMinutesAndSeconds } from './helpers';

import { PlayerContext } from '@/components/song-requests/context';

const useStyles = createStyles((theme) => ({
  item: {
    backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.white,
  },

  dragHandle: {
    ...theme.fn.focusStyles(),
    width: 40,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    height: '100%',
    color: theme.colorScheme === 'dark' ? theme.colors.dark[1] : theme.colors.gray[6],
  },
}));

export const VideosList: React.FC = () => {
  const { videos, reorderVideos, skipVideo } = useContext(PlayerContext);

  const [isBrowser, setIsBrowser] = useState(false);
  const { classes } = useStyles();

  useEffect(() => {
    setIsBrowser(process.browser);
  }, []);


  const items = videos.map((video, index) => (
    <Draggable key={video.id} index={index} draggableId={video.id}>
      {(provided) => (
        <tr className={classes.item} ref={provided.innerRef} {...provided.draggableProps} hidden={index === 0}>
          <td>
            <div className={classes.dragHandle} {...provided.dragHandleProps}>
              <IconGripVertical size={18} stroke={1.5}/>
            </div>
          </td>
          <td>{index}</td>
          <td><a href={'https://youtu.be/' + video.videoId}>{video.title}</a></td>
          <td>{video.orderedByName}</td>
          <td>{millisToMinutesAndSeconds(video.duration)}</td>
          <td>
            <Flex>
              <ActionIcon variant={'filled'} color={'red'} onClick={() => skipVideo(index)}><IconTrash size={14}
              /></ActionIcon>
            </Flex>
          </td>
        </tr>
      )}
    </Draggable>
  ));

  return (
    <>{isBrowser
      ? <DragDropContext
        onDragEnd={({ destination, source }) => {
          reorderVideos(destination!, source);
        }}
      >
        <Table>
          <thead>
          <tr>
            <th style={{ width: 40 }}></th>
            <th>#</th>
            <th>Title</th>
            <th>Requested by</th>
            <th>Duration</th>
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
      : <div></div>}
    </>);
};
