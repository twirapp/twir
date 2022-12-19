import { ActionIcon, Badge, Button, Flex, Switch, Table, Text } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { IconPencil, IconTrash } from '@tabler/icons';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { useState } from 'react';

import { TimerDrawer } from '../components/timers/drawer';

import { confirmDelete } from '@/components/confirmDelete';
import { useTimersManager } from '@/services/api';

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableTimer, setEditableTimer] = useState<ChannelTimer | undefined>();
  const viewPort = useViewportSize();

  const manager = useTimersManager();
  const { data: timers } = manager.getAll();

  return (
    <div>
      <Flex direction="row" justify="space-between">
        <Text size="lg">Timers</Text>
        <Button
          color="green"
          onClick={() => {
            setEditableTimer(undefined);
            setEditDrawerOpened(true);
          }}
        >
          Create
        </Button>
      </Flex>
      <Table>
        <thead>
          <tr>
            <th>Name</th>
            {viewPort.width > 550 && <th>Responses</th>}
            <th>Time Interval</th>
            <th>Messages Interval</th>
            {viewPort.width > 550 && <th>Status</th>}
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {timers &&
            timers.map((timer, idx) => (
              <tr key={timer.id}>
                <td>
                  <Badge>{timer.name}</Badge>
                </td>
                {viewPort.width > 550 && (
                  <td>
                    {timer.responses.map((r, i) => (
                      <p key={i} style={{ margin: 0 }}>
                        {r.text}
                      </p>
                    ))}
                  </td>
                )}

                <td>{timer.timeInterval} seconds</td>
                <td>{timer.messageInterval}</td>
                {viewPort.width > 550 && (
                  <td>
                    <Switch
                      checked={timer.enabled}
                      onChange={(event) => (timer.enabled = event.currentTarget.checked)}
                    />
                  </td>
                )}
                <td>
                  <Flex direction="row" gap="xs">
                    <ActionIcon
                      onClick={() => {
                        setEditableTimer(timers[idx] as any);
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
                          onConfirm: () => manager.delete(timer.id),
                        })
                      }
                      variant="filled"
                      color="red"
                    >
                      <IconTrash size={14} />
                    </ActionIcon>
                  </Flex>
                </td>
              </tr>
            ))}
        </tbody>
      </Table>

      <TimerDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        timer={editableTimer}
      />
    </div>
  );
}
