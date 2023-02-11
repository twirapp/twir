import { ActionIcon, Badge, Table } from '@mantine/core';
import { IconPencil } from '@tabler/icons';
import type { Event, EventType } from '@tsuwari/typeorm/entities/events/Event';
import { OperationType } from '@tsuwari/typeorm/entities/events/EventOperation';
import { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { EventsDrawer } from '@/components/events/drawer';
import { Greeting } from '@/services/api';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['layout'])),
  },
});

const Events: NextPage<{ operations: typeof OperationType }> = (props) => {
  const events: Event[] = [
    {
      id: '',
      type: 'FOLLOW' as EventType,
      rewardId: '',
      commandId: '',
      operations: [
        {
          type: 'SEND_MESSAGE' as OperationType,
          id: '',
          delay: 0,
          input: 'Im message',
          eventId: '',
          repeat: 1,
          order: 0,
        },
      ],
      description: 'Send message when user follows channel',
      channelId: '',
    },
    {
      id: '',
      type: 'COMMAND_USED' as EventType,
      rewardId: '',
      commandId: '',
      operations: [
        {
          type: 'SEND_MESSAGE' as OperationType,
          id: '',
          delay: 0,
          input: 'Im message',
          eventId: '',
          repeat: 1,
          order: 0,
        },
      ],
      description: 'Send message when user follows channel',
      channelId: '',
    },
    {
      id: '',
      type: 'REDEMPTION_CREATED' as EventType,
      rewardId: '',
      commandId: '',
      operations: [
        {
          type: 'SEND_MESSAGE' as OperationType,
          id: '',
          delay: 0,
          input: 'Im message',
          eventId: '',
          repeat: 1,
          order: 0,
        },
      ],
      description: 'Send message when user follows channel',
      channelId: '',
    },
  ];
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableEvent, setEditableEvent] = useState<Event | undefined>();

  return (
    <>
      {props.operations}
      <Table>
        <thead>
        <tr>
          <th>Event</th>
          <th>Description</th>
          <th>Actions</th>
        </tr>
        </thead>
        <tbody>
        {events.map((e, idx) => <tr key={e.id}>
          <td><Badge>{e.type}</Badge></td>
          <td>{e.description}</td>
          <td>
            <ActionIcon
              onClick={() => {
                setEditableEvent(events[idx] as any);
                setEditDrawerOpened(true);
              }}
              variant="filled"
              color="blue"
            >
              <IconPencil size={14} />
            </ActionIcon>
          </td>
        </tr>)}
        </tbody>
      </Table>

      <EventsDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        event={editableEvent}
      />
    </>
  );
};

export default Events;