import { Badge, Button, Switch, Table } from '@mantine/core';
import { useViewportSize } from '@mantine/hooks';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { useState } from 'react';

import { TimerDrawer } from '../components/timers/drawer';

import { useTimers } from '@/services/api';

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableTimer, setEditableTimer] = useState<ChannelTimer>({} as any);
  const viewPort = useViewportSize();

  const { data: timers } = useTimers();

  return (
    <div>
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
            timers.map((element, idx) => (
              <tr key={element.id}>
                <td>
                  <Badge>{element.name}</Badge>
                </td>
                {viewPort.width > 550 && (
                  <td>
                    {element.responses.map((r, i) => (
                      <p key={i} style={{ margin: 0 }}>
                        {r.text}
                      </p>
                    ))}
                  </td>
                )}

                <td>{element.timeInterval} seconds</td>
                <td>{element.messageInterval}</td>
                {viewPort.width > 550 && (
                  <td>
                    <Switch
                      checked={element.enabled}
                      onChange={(event) => (element.enabled = event.currentTarget.checked)}
                    />
                  </td>
                )}
                <td>
                  <Button
                    onClick={() => {
                      setEditableTimer(timers[idx] as any);
                      setEditDrawerOpened(true);
                    }}
                  >
                    Edit
                  </Button>
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
