import { Badge, Button, Switch, Table } from '@mantine/core';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { useState } from 'react';

import { KeywordDrawer } from '@/components/keywords/drawer';
import { useKeywords } from '@/services/api';

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableKeyword, setEditableKeyword] = useState<ChannelKeyword>({} as any);

  const { data: keywords } = useKeywords();

  return (
    <div>
      <Table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {keywords &&
            keywords.map((element, idx) => (
              <tr key={element.id}>
                <td>
                  <Badge>{element.text}</Badge>
                </td>
                <td>
                  <Switch
                    checked={element.enabled}
                    onChange={(event) => (element.enabled = event.currentTarget.checked)}
                  />
                </td>
                <td>
                  <Button
                    onClick={() => {
                      setEditableKeyword(keywords[idx] as any);
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

      <KeywordDrawer
        opened={editDrawerOpened}
        setOpened={setEditDrawerOpened}
        keyword={editableKeyword}
      />
    </div>
  );
}
