import { Badge, Button, Table } from '@mantine/core';
import { useState } from 'react';

import { GreetingDrawer } from '../components/greetings/drawer';
import { type Greeting, useGreetings } from '../services/api/greetings';


export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableGreeting, setEditableGreeting] = useState<Greeting>({} as any);

  const { data: greetings } = useGreetings();

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
