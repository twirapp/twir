import { Badge, Button, Switch, Table } from '@mantine/core';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { useState } from 'react';

import { TimerDrawer } from '../components/timers/drawer';

const timers = [
  {
    id: 'e27d9ddc-8e35-4c02-a964-53e8d7e58333',
    channelId: '128644134',
    name: 'Телега',
    enabled: true,
    timeInterval: 11,
    messageInterval: 5,
    lastTriggerMessageNumber: 1400,
    responses: [
      {
        id: '1fe70a65-d8c9-4952-b37c-6634b1e4d0bf',
        text: 'Пишу около айти, около моей жизни: https://t.me/satontdev',
        isAnnounce: false,
        timerId: 'e27d9ddc-8e35-4c02-a964-53e8d7e58333',
      },
    ],
  },
  {
    id: '0ad53238-6b63-4073-8b1e-628c2a02807d',
    channelId: '128644134',
    name: 'Donate',
    enabled: true,
    timeInterval: 10,
    messageInterval: 5,
    lastTriggerMessageNumber: 1400,
    responses: [
      {
        id: '5fe6b730-852b-4e3e-b1ad-c6711af15a72',
        text: 'https://www.donationalerts.com/r/s4tont Для украинцев: https://new.donatepay.ru/@88670',
        isAnnounce: false,
        timerId: '0ad53238-6b63-4073-8b1e-628c2a02807d',
      },
    ],
  },
];
export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableTimer, setEditableTimer] = useState<ChannelTimer>({} as any);

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
          {timers.map((element, idx) => (
            <tr key={element.id}>
              <td>
                <Badge>{element.name}</Badge>
              </td>
              <td>{element.responses.map((r) => r.text).join('\n')}</td>
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
