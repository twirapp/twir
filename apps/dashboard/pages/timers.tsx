import { Badge, Button, Switch, Table } from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { Dashboard } from '@tsuwari/shared';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { useState } from 'react';
import useSWR from 'swr';

import { TimerDrawer } from '../components/timers/drawer';
import { swrFetcher } from '../services/swrFetcher';

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableTimer, setEditableTimer] = useState<ChannelTimer>({} as any);
  const [selectedDashboard] = useLocalStorage<Dashboard>({
    key: 'selectedDashboard',
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });

  const { data: timers } = useSWR<ChannelTimer[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/timers` : null,
    swrFetcher,
  );

  return (
    <div>
      <Table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Responses</th>
            <th>Time Interval</th>
            <th>Messages Interval</th>
            <th>Status</th>
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
                <td>
                  {element.responses.map((r, i) => (
                    <p key={i} style={{ margin: 0 }}>
                      {r.text}
                    </p>
                  ))}
                </td>
                <td>{element.timeInterval} seconds</td>
                <td>{element.messageInterval}</td>
                <td>
                  <Switch
                    checked={element.enabled}
                    onChange={(event) => (element.enabled = event.currentTarget.checked)}
                  />
                </td>
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
