import { Badge, Button, Table } from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { Dashboard } from '@tsuwari/shared';
import { useState } from 'react';
import useSWR from 'swr';

import { GreetingDrawer, type Greeting } from '../components/greetings/drawer';
import { swrFetcher } from '../services/swrFetcher';


export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableGreeting, setEditableGreeting] = useState<Greeting>({} as any);
  const [selectedDashboard] = useLocalStorage<Dashboard>({
    key: 'selectedDashboard',
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });

  const { data: greetings } = useSWR<Greeting[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/greetings` : null,
    swrFetcher,
  );

  return (
    <div>
      <Table>
        <thead>
          <tr>
            <th>Username</th>
            <th>Message</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {greetings &&
            greetings.map((element, idx) => (
              <tr key={element.id}>
                <td>
                  <Badge>{element.userName}</Badge>
                </td>
                <td>
                  <Badge color="cyan">{element.text}</Badge>
                </td>
                <td>
                  <Button
                    onClick={() => {
                      setEditableGreeting(greetings[idx] as any);
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

      <GreetingDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        greeting={editableGreeting}
      />
    </div>
  );
}
