import { ActionIcon, Badge, Button, Flex, Group, Table, Text, TextInput } from '@mantine/core';
import { IconPencil, IconSearch, IconTrash } from '@tabler/icons';
import type { Event, EventType } from '@tsuwari/typeorm/entities/events/Event';
import { OperationType } from '@tsuwari/typeorm/entities/events/EventOperation';
import { NextPage } from 'next';
import { serverSideTranslations } from 'next-i18next/serverSideTranslations';
import { useState } from 'react';

import { confirmDelete } from '@/components/confirmDelete';
import { EventsDrawer } from '@/components/events/drawer';
import { eventsManager } from '@/services/api';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
export const getServerSideProps = async ({ locale }) => ({
  props: {
    ...(await serverSideTranslations(locale, ['events', 'layout'])),
  },
});

const Events: NextPage<{ operations: typeof OperationType }> = (props) => {
  const manager = eventsManager();
  const { data: events } = manager.useGetAll();
  const deleter = manager.useDelete();

  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableEvent, setEditableEvent] = useState<Event | undefined>();

  return (
    <>
      <Flex direction="row" justify="space-between">
        <Group>
          <Text size="lg">Events</Text>
        </Group>
        <Button
          color="green"
          onClick={() => {
            setEditableEvent(undefined);
            setEditDrawerOpened(true);
          }}
        >
          Create
        </Button>
      </Flex>
      <Table>
        <thead>
        <tr>
          <th>Event</th>
          <th>Description</th>
          <th>Actions</th>
        </tr>
        </thead>
        <tbody>
        {events?.map((e, idx) => <tr key={e.id}>
          <td><Badge>{e.type.split('_').join(' ')}</Badge></td>
          <td>{e.description}</td>
          <td>
            <Flex direction="row" gap="xs">
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
              <ActionIcon
                onClick={() =>
                  confirmDelete({
                    onConfirm: () => deleter.mutate(e.id),
                  })
                }
                variant="filled"
                color="red"
              >
                <IconTrash size={14} />
              </ActionIcon>
            </Flex>
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