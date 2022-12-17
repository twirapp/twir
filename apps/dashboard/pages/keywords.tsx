import { Badge, Button, Switch, Table } from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { Dashboard } from '@tsuwari/shared';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { useState } from 'react';
import useSWR from 'swr';

import { KeywordDrawer } from '../components/keywords/drawer';
import { swrFetcher } from '../services/swrFetcher';

export default function () {
  const [editDrawerOpened, setEditDrawerOpened] = useState(false);
  const [editableKeyword, setEditableKeyword] = useState<ChannelKeyword>({} as any);
  const [selectedDashboard] = useLocalStorage<Dashboard>({
    key: 'selectedDashboard',
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });

  const { data: keywords } = useSWR<ChannelKeyword[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/keywords` : null,
    swrFetcher,
  );

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
